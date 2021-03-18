package helpers

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
)

func GetCreateMigrationsTableSql(tableName string) string {
	return fmt.Sprintf(`
CREATE TABLE IF NOT EXISTS %s
(
    id integer not null primary key autoincrement,
    migration varchar(255) not null unique,
    version int not null,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);
`, tableName)
}

func GetCreateUsersTableSql() string {
	return `
CREATE TABLE users
(
    id integer not null primary key autoincrement,
    name varchar (255) not null,
    email  varchar (255) not null unique,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
)
`
}

func GetDropUsersTableSql() string {
	return `DROP TABLE users;`
}

func GetCreateListsTableSql() string {
	return `
CREATE TABLE lists
(
    id integer not null primary key autoincrement,
    label varchar (255) not null,
    description  varchar (255) not null unique,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
)
`
}

func GetDropListsTableSql() string {
	return `DROP TABLE lists;`
}
