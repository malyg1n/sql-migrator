#slqx-migrator

[![Version](https://img.shields.io/badge/version-v0.0.1-green.svg)](https://github.com/malyg1n/sqlx-migrator/releases)

Golang пакет для создания миграций с использованием [`sqlx`](https://github.com/jmoiron/sqlx).

## Использование

### Установка
```bigquery
go get -u github.com/malyg1n/sqlx-migrator
```

### Создание файлов миграций
В корне проекта необходимо выполнить команду со следующей сигнатурой:
```bigquery
sqlx-migrator create [migrations-directory] migration-name
```
Пример:
```bigquery
sqlx-migrator create migrations create-users-table
```
После этого в папке migrations в корне проекта появятся два файла
```bigquery
[date]-create-users-table-up.sql
[date]-create-users-table-down.sql
```
В них необходимо прописать SQL код для накатывая и откатывания миграции соотвественно
### Миграции
Чтобы накатить миграции, воспользуйтесь командой:
```bigquery
sqlx-migrator up
```
B чтобы откатить:
```bigquery
sqlx-migrator down
```