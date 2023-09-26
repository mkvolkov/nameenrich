package main

import (
	"context"
	"log"
	"nameenrich/aserver"
	"nameenrich/cfg"
	"nameenrich/enrich"
	"nameenrich/graph"
	"nameenrich/storage"
	"nameenrich/types"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-json"
	"github.com/gomodule/redigo/redis"
	"github.com/segmentio/kafka-go"
)

const (
	defaultPort = "8080"

	topic       = "FIO"
	failedTopic = "FIO_FAILED"
)

func main() {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	cMainCfg := &cfg.Cfg{}
	err := cfg.LoadConfig(cMainCfg)
	LogError(errorLog, "Error in LoadConfig: ", err)

	if cMainCfg.Port == "" {
		cMainCfg.Port = defaultPort
	}

	dBase := graph.Connect(cMainCfg.PostgresURL)

	// конкатенация строк может быть быстрее, но здесь эта операция одноразовая
	redisAddr := cMainCfg.RedisHost + ":" + cMainCfg.RedisPort
	rConn, err := redis.Dial("tcp", redisAddr)
	if err != nil {
		errorLog.Fatalln("Couldn't connect to Redis: ", err)
	}
	defer rConn.Close()

	kLog := log.New(os.Stdout, "kafka reader: ", 0)

	kReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   topic,
		Logger:  kLog,
	})

	kWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   failedTopic,
		Logger:  kLog,
	})

	go func() {
		for {
			msg, _ := kReader.ReadMessage(context.Background())

			var decMsg types.MsgBase = types.MsgBase{}

			err = CheckMsg(&msg, &decMsg, errorLog)
			if err != nil {
				var msgFailed types.MsgError
				msgFailed.Errname = err.Error()
				msgFailed.IncorrectMsg = string(msg.Value)

				msgData, err := json.Marshal(msgFailed)
				LogError(errorLog, "Marshalling error", err)

				errWrite := kWriter.WriteMessages(context.Background(), kafka.Message{
					Value: msgData,
				})
				LogError(errorLog, "Kafka writing error", errWrite)

				continue
			}

			var enrMsg types.MsgEnriched = types.MsgEnriched{}

			err = enrich.Enrichment(&decMsg, &enrMsg)
			LogError(errorLog, "Enrichment error", err)

			err = storage.WriteData(dBase, rConn, &enrMsg)
			LogError(errorLog, "Writing data error", err)
		}
	}()

	aCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		aServer := aserver.NewServer(cMainCfg.AtrHost, cMainCfg.AtrPort, dBase, rConn)

		err = aServer.Run(aCtx)
		if err != nil {
			log.Fatalln("Couldn't run Atreugo server, exiting...")
		}
	}()

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{DB: dBase, RConn: rConn}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	infoLog.Printf("connect to http://localhost:%s/ for GraphQL playground", cMainCfg.Port)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	finishCh := make(chan struct{})

	go func() {
		err := http.ListenAndServe(":"+cMainCfg.Port, nil)
		if err != nil {
			log.Fatalf("Couldn't start GraphQL server: %v\n", err)
		}
	}()

	go func() {
		s := <-signalCh
		infoLog.Printf("\ngot signal %v, graceful shutdown...", s)
		dBase.Close()
		finishCh <- struct{}{}
	}()

	<-finishCh
	infoLog.Println("Finished shutdown")
}

func CheckMsg(kMsg *kafka.Message, dMsg *types.MsgBase, errLog *log.Logger) error {
	err := json.Unmarshal(kMsg.Value, dMsg)
	if err != nil {
		errLog.Println("Error unmarshal: ", err.Error())
		return err
	}

	v := validator.New()
	err = v.Struct(dMsg)
	if err != nil {
		errLog.Println("Error validate: ", err.Error())
		return err
	}

	return nil
}

func LogError(lg *log.Logger, msg string, err error) {
	if err != nil {
		lg.Println(err.Error())
	}
}
