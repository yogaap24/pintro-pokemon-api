CREATE TABLE IF NOT EXISTS users (
    id         bytea       NOT NULL,
    username   text        NOT NULL UNIQUE,
    password   text        NOT NULL,
    role       text        NOT NULL,
    created_at timestamptz NOT NULL,
    updated_at timestamptz,
    deleted_at timestamptz,

    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS users_pokemons (
    id          bytea       NOT NULL,
    user_id     bytea       NOT NULL,
    pokemon_id  int         NOT NULL,
    nickname    char(50)    NOT NULL,
    captured_at timestamptz NOT NULL,
    released   boolean     NOT NULL DEFAULT FALSE,

    PRIMARY KEY(id),
    FOREIGN KEY(user_id) REFERENCES users(id)
);