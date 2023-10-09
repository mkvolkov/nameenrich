# nameenrich

Сервер GraphQL и Atreugo, слушающий очередь Kafka, обогащающий
данные, пишущий их в PostgreSQL и осуществляющий кеширование
в Redis.

Для запуска сервера необходимы установленные и запущенные Kafka
с темами FIO и FIO_FAILED, Redis, PostgreSQL.

Для создания базы данных Postgres (название: enrichment) необходим
инструмент migrate.

Создать базу данных:

Шаг 1 - создать базу данных enrichment (выполняется в psql)

```
create database enrichment;
```

Шаг 2 - создать таблицы с помощью миграции. Переменная среды POSTGRESQL_URL
должна совпадать с переменной POSTGRES_URL, указанной в конфигурации в файле
app.env

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