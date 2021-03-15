CREATE TABLE IF NOT EXISTS schema_migrations
(
    id bigserial not null primary key,
    migration varchar(255) not null unique,
    version int not null,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);