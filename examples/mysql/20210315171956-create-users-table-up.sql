CREATE TABLE IF NOT EXISTS users
(
    id bigserial not null primary key,
    name varchar (255) not null,
    email  varchar (255) not null unique,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
)