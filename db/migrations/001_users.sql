-- +migrate Up

CREATE TABLE users(
  id BIGSERIAL NOT NULL  PRIMARY KEY,
  name varchar(255) NOT NULL,
  email varchar(254) NOT NULL UNIQUE,
  date_of_birth varchar(255) NOT NULL,
  password varchar(255) NOT NULL,
  phone varchar(50) NOT NULL
);

-- +migrate Down

DROP TABLE users;