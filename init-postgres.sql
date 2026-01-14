-- Разрешить подключение без пароля для пользователя postgres из Docker сети
ALTER USER postgres PASSWORD '0611';

CREATE SEQUENCE IF NOT EXISTS users_id_seq;
CREATE SEQUENCE IF NOT EXISTS box_id_seq;
CREATE SEQUENCE IF NOT EXISTS box_cars_id_seq;
CREATE SEQUENCE IF NOT EXISTS cars_id_seq;
CREATE SEQUENCE IF NOT EXISTS transactions_id_seq;
CREATE SEQUENCE IF NOT EXISTS users_cars_id_seq;


-- Создать базу данных если не существует (для PostgreSQL < 9.5 использовать DO $$ ... $$)
DO $$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'go-casino') THEN
      PERFORM dblink_exec('dbname=postgres', 'CREATE DATABASE "go-casino"');
   END IF;
END
$$;

-- Предоставить права
GRANT ALL PRIVILEGES ON DATABASE "go-casino" TO postgres;




CREATE TABLE IF NOT EXISTS public.users
(
    id integer NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    password character varying(255) COLLATE pg_catalog."default" NOT NULL,
    balance numeric NOT NULL,
    role character varying COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.boxes
(
    id integer NOT NULL DEFAULT nextval('box_id_seq'::regclass),
    mode character varying COLLATE pg_catalog."default" NOT NULL,
    price numeric NOT NULL,
    name character varying COLLATE pg_catalog."default",
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    CONSTRAINT box_pkey PRIMARY KEY (id)
);

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.boxes
    OWNER to postgres;


CREATE TABLE IF NOT EXISTS public.box_cars
(
    id integer NOT NULL DEFAULT nextval('box_cars_id_seq'::regclass),
    box_id integer NOT NULL,
    car_id integer NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    CONSTRAINT box_cars_pkey PRIMARY KEY (id),
    CONSTRAINT box_id_fk FOREIGN KEY (box_id)
        REFERENCES public.boxes (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT car_id_fk FOREIGN KEY (car_id)
        REFERENCES public.cars (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.box_cars
    OWNER to postgres;



CREATE TABLE IF NOT EXISTS public.cars
(
    id integer NOT NULL DEFAULT nextval('cars_id_seq'::regclass),
    name character varying COLLATE pg_catalog."default" NOT NULL,
    price numeric NOT NULL,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    CONSTRAINT cars_pkey PRIMARY KEY (id)
);

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.cars
    OWNER to postgres;



CREATE TABLE IF NOT EXISTS public.roles
(
    role character varying COLLATE pg_catalog."default" NOT NULL,
    description character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT roles_pkey PRIMARY KEY (role)
);

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.roles
    OWNER to postgres;



CREATE TABLE IF NOT EXISTS public.transactions
(
    id integer NOT NULL DEFAULT nextval('transactions_id_seq'::regclass),
    amount numeric NOT NULL,
    user_id integer NOT NULL,
    created_at time with time zone NOT NULL,
    updated_at time with time zone,
    CONSTRAINT transactions_pkey PRIMARY KEY (id),
    CONSTRAINT user_id_fk FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
);

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.transactions
    OWNER to postgres;




CREATE TABLE IF NOT EXISTS public.users
(
    id integer NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    name character varying(255) COLLATE pg_catalog."default" NOT NULL,
    password character varying(255) COLLATE pg_catalog."default" NOT NULL,
    balance numeric NOT NULL,
    role character varying COLLATE pg_catalog."default" NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.users
    OWNER to postgres;



CREATE TABLE IF NOT EXISTS public.users_cars
(
    id integer NOT NULL DEFAULT nextval('users_cars_id_seq'::regclass),
    user_id integer NOT NULL,
    car_id integer NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    CONSTRAINT users_cars_pkey PRIMARY KEY (id),
    CONSTRAINT car_id_fk FOREIGN KEY (car_id)
        REFERENCES public.cars (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT user_id_fk FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.users_cars
    OWNER to postgres;