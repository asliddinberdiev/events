CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  telegram_id VARCHAR NOT NULL UNIQUE,
  username VARCHAR NOT NULL UNIQUE,
  password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS events (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  creater_id VARCHAR NOT NULL REFERENCES users(telegram_id),
  title VARCHAR NOT NULL,
  description TEXT DEFAULT '',
  date DATE NOT NULL
);
