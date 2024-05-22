-- Table: public.user

-- DROP TABLE IF EXISTS public."user";

CREATE TABLE "user"(
    id SERIAL PRIMARY KEY,
    login varchar(255) NOT NULL UNIQUE,
    password varchar(255) NOT NULL
);
