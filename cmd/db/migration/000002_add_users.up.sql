CREATE TABLE "users" (
  "email" text PRIMARY KEY,
  "hashed_password" text NOT NULL,
  "first_name" text NOT NULL,
  "last_name" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("email");