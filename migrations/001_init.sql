CREATE TABLE subscriptions (
  id SERIAL PRIMARY KEY,
  email TEXT NOT NULL UNIQUE,
  city TEXT NOT NULL,
  frequency TEXT NOT NULL,
  confirmed BOOLEAN NOT NULL DEFAULT FALSE,
  token TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT now(),
  confirmed_at TIMESTAMP
);
