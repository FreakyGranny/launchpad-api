CREATE TABLE users (
    id bigserial NOT NULL primary key,
    username varchar NOT NULL,
    first_name varchar NOT NULL,
    last_name varchar NOT NULL,
    avatar varchar NOT NULL,
    email varchar NOT NULL unique,
    is_admin boolean NOT NULL DEFAULT FALSE,
    project_count int NOT NULL DEFAULT 0,
    success_rate numeric NOT NULL DEFAULT 0.0
);
