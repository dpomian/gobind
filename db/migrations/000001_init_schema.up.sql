CREATE TABLE IF NOT EXISTS notebooks (
    id UUID NOT NULL,
    title character varying(255) NOT NULL,
    topic character varying(255) NOT NULL default('Misc'),
    content text NOT NULL default(''),
    deleted boolean NOT NULL default(false),
    last_modified timestamp with time zone NOT NULL default(now()),
    created_at timestamp with time zone NOT NULL
);

CREATE INDEX idx_notebook_id ON notebooks(id);