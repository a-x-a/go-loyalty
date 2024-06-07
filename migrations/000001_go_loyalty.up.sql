-- "user"
-- create table
CREATE TABLE IF NOT EXISTS "user"(
    id SERIAL PRIMARY KEY,
    login varchar NOT NULL UNIQUE,
    password varchar NOT NULL
);

-- balance
-- create table
CREATE TABLE IF NOT EXISTS "balance"
(
    user_id int REFERENCES "user"(id) PRIMARY KEY,
    current numeric(15,2) NOT NULL DEFAULT 0,
    withdrawn numeric(15,2) NOT NULL DEFAULT 0
);

-- withdraw
-- create table
CREATE TABLE IF NOT EXISTS "withdraw"
(
    user_id int REFERENCES "user"(id) ,
    "order" varchar NOT NULL UNIQUE,
    sum numeric(15,2) NOT NULL DEFAULT 0,
    processed_at timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, "order")
);

-- "order"
-- create table
CREATE TABLE IF NOT EXISTS "order"
(
    "number" varchar PRIMARY KEY,
    user_id int REFERENCES "user"(id) ,
    status int NOT NULL DEFAULT 1,
    accrual numeric(15,2) NOT NULL DEFAULT 0,
    uploaded_at timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP
);