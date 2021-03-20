# sql-migrator

[![Version](https://img.shields.io/badge/version-v0.0.6-green.svg)](https://github.com/malyg1n/sql-migrator/releases)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/malyg1n/sql-migrator/blob/master/LICENSE.md)
[![Coverage Status](https://coveralls.io/repos/github/malyg1n/sql-migrator/badge.svg?branch=master)](https://coveralls.io/github/malyg1n/sql-migrator?branch=master)

Golang utility for managing migrations using [`database/sql`](https://golang.org/pkg/database/sql).

The package includes the following drivers: [`postgres`](https://github.com/lib/pq), [`mysql`](https://github.com/go-sql-driver/mysql), [`sqlite3`](https://github.com/mattn/go-sqlite3).
## Usage

### Installation and setup
```
go get -u github.com/malyg1n/sql-migrator
```
Create a file `.env` by copying from [`.env.example`](https://github.com/malyg1n/sql-migrator/tree/master/examples/.env.example) and specify your database settings.
In the same file, specify the folder for migrations (`migrations` in root of project by default). 
After that, tou need to run the command:
```bigquery
sql-migrator init
```
This command creates table and directory (if not exists) to store your migrations.
### Create migration files.
At the root of the project, you need to run a command with the following signature:
```bigquery
sql-migrator create migration-name
```
Like this:
```bigquery
sql-migrator create create-users-table
```
After that, two files will appear in the migrations' folder at the root of the project.
```bigquery
00001-create-users-table-up.sql
00001-create-users-table-down.sql
```
They need to write SQL code for rolling and rolling back migration, respectively.

See migrations' examples [here](https://github.com/malyg1n/sql-migrator/tree/master/examples).
### Migrations
To roll out migrations, use the command:
```bigquery
sql-migrator up
```
And to rollback:
```bigquery
sql-migrator down
```
Complete cleaning of all migrations and roll out them over again:
```bigquery
sql-migrator refresh
```
Complete cleaning of all migrations:
```bigquery
sql-migrator clean
```
