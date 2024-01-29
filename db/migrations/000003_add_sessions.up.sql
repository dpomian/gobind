CREATE TABLE sessions (
    id UUID NOT NULL UNIQUE PRIMARY KEY,
    user_id UUID NOT NULL,
    refresh_token VARCHAR NOT NULL,
    user_agent VARCHAR NOT NULL,
    client_ip VARCHAR NOT NULL,
    is_blocked boolean NOT NULL default false,
    expires_at timestamp with time zone NOT NULL
);

ALTER TABLE sessions ADD FOREIGN KEY (user_id) REFERENCES users(id);