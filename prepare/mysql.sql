CREATE TABLE IF NOT EXISTS schema_migrations
(
    id integer not null primary key auto_increment,
    migration varchar(255) not null unique,
    version int not null,
    created_at timestamp DEFAULT CURRENT_TIMESTAMP
);