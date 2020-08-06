CREATE TABLE categories (
    id bigserial NOT NULL primary key,
    alias varchar NOT NULL,
    name varchar NOT NULL
);

INSERT INTO categories (alias, name) VALUES ('other', 'Other');
INSERT INTO categories (alias, name) VALUES ('video_games', 'Video games');
INSERT INTO categories (alias, name) VALUES ('board_games', 'Board games');
INSERT INTO categories (alias, name) VALUES ('party', 'Party');
