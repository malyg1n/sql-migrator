CREATE TABLE IF NOT EXISTS users
(
    id integer not null primary key autoincrement,
    name varchar (255) not null,
    email  varchar (255) not null unique,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
)