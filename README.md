# nameenrich

Сервер GraphQL и Atreugo, слушающий очередь Kafka, обогащающий
данные, пишущий их в PostgreSQL и осуществляющий кеширование
в Redis.

Для запуска сервера необходимы установленные и запущенные Kafka
с темами FIO и FIO_FAILED, Redis, PostgreSQL.

Для создания базы данных Postgres (название: enrichment) необходим
инструмент migrate.

Создать базу данных:

```
migrate -database ${POSTGRESQL_URL} -path db/migrations up
```

Запустить сервер:

```
go run server.go
```

Получить данные о человеке по ID через GraphQL:

```
mutation getperson($id: Int!) {
    getPerson(id: $id) {
        name
        surname
    }
}


variables
{
    "id": 1
}
```