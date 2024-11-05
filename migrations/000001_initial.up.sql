CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  username VARCHAR NOT NULL UNIQUE,
  telegram_id VARCHAR NOT NULL UNIQUE,
  password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS remembers (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  creater_id VARCHAR NOT NULL REFERENCES users(telegram_id),
  title VARCHAR NOT NULL,
  description TEXT DEFAULT '',
  birthday DATE NOT NULL
);
