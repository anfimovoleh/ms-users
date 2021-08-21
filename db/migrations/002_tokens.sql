-- +migrate Up

CREATE TABLE tokens(
  token varchar(128) PRIMARY KEY ,
  user_id integer,
  last_sent_at timestamp without time zone,
  FOREIGN KEY (user_id) REFERENCES users(id)
);

-- +migrate Down

DROP TABLE tokens;