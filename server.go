package main

import (
	"context"
	"fmt"
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
	cMainCfg := &cfg.Cfg{}
	err := cfg.LoadConfig(cMainCfg)
	CheckError(err)

	if cMainCfg.Port == "" {
		cMainCfg.Port = defaultPort
	}

	dBase := graph.Connect(cMainCfg.PostgresURL)

	redisAddr := fmt.Sprintf("%s:%s", cMainCfg.RedisHost, cMainCfg.RedisPort)
	rConn, err := redis.Dial("tcp", redisAddr)
	if err != nil {
		log.Fatal(err)
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
			msg, err := kReader.ReadMessage(context.Background())
			if err != nil {
				log.Println("couldn't read msg: ", err.Error())
			}

			var decMsg types.MsgBase = types.MsgBase{}

			err = CheckMsg(&msg, &decMsg)
			if err != nil {
				var msgFailed types.MsgError
				msgFailed.Errname = err.Error()
				msgFailed.IncorrectMsg = string(msg.Value)

				msgData, err := json.Marshal(msgFailed)
				CheckError(err)

				errWrite := kWriter.WriteMessages(context.Background(), kafka.Message{
					Value: msgData,
				})
				CheckError(errWrite)

				continue
			}

			var enrMsg types.MsgEnriched = types.MsgEnriched{}

			enrich.Enrichment(&decMsg, &enrMsg)

			err = storage.WriteData(dBase, rConn, &enrMsg)
			CheckError(err)
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

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", cMainCfg.Port)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	finishCh := make(chan struct{})

	go func() {
		err := http.ListenAndServe(":"+cMainCfg.Port, nil)
		if err != nil {
			log.Fatalf("couldn't start server GraphQL: %v\n", err)
		}
	}()

	go func() {
		s := <-signalCh
		log.Printf("\ngot signal %v, graceful shutdown...", s)
		dBase.Close()
		finishCh <- struct{}{}
	}()

	<-finishCh
	fmt.Println("Finished shutdown")
}

func CheckMsg(kMsg *kafka.Message, dMsg *types.MsgBase) error {
	err := json.Unmarshal(kMsg.Value, dMsg)
	if err != nil {
		fmt.Println("Error unmarshal: ", err.Error())
		return err
	}

	v := validator.New()
	err = v.Struct(dMsg)
	if err != nil {
		fmt.Println("Error validate: ", err.Error())
		return err
	}

	return nil
}

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}
