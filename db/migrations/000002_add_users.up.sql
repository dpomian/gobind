CREATE TABLE users (
    id UUID NOT NULL UNIQUE PRIMARY KEY,
    username VARCHAR NOT NULL,
    email VARCHAR UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
    suspended BOOLEAN NOT NULL default false,
    created_at timestamp with time zone NOT NULL default(now()),
    password_changed_at timestamp with time zone NOT NULL default '0001-01-01 00:00:00Z'
);

ALTER TABLE notebooks ADD COLUMN user_id UUID NOT NULL;
ALTER TABLE notebooks ADD FOREIGN KEY (user_id) REFERENCES users(id);