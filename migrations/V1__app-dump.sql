
SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;
SET default_tablespace = '';
SET default_table_access_method = heap;

CREATE TABLE public.word (
      id integer NOT NULL,
      name character varying(255),
      created_at timestamp without time zone,
      updated_at timestamp without time zone
);

ALTER TABLE public.word OWNER TO root;

CREATE SEQUENCE public.word_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE public.word_seq OWNER TO root;

ALTER SEQUENCE public.word_seq OWNED BY public.word.id;

ALTER TABLE ONLY public.word ALTER COLUMN id SET DEFAULT nextval('public.word_seq'::regclass);

ALTER TABLE ONLY public.word
    ADD CONSTRAINT word_pk PRIMARY KEY (id),
    ADD CONSTRAINT word_name UNIQUE (name);


CREATE TABLE public.user_statistic (
     id integer NOT NULL,
     user_id integer NOT NULL,
     word_id integer NOT NULL,
     requested integer NOT NULL,
     created_at timestamp without time zone,
     updated_at timestamp without time zone
);

ALTER TABLE public.user_statistic OWNER TO root;

CREATE SEQUENCE public.user_statistic_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER TABLE public.user_statistic_seq OWNER TO root;

ALTER SEQUENCE public.user_statistic_seq OWNED BY public.user_statistic.id;

ALTER TABLE ONLY public.user_statistic ALTER COLUMN id SET DEFAULT nextval('public.user_statistic_seq'::regclass);

ALTER TABLE ONLY public.user_statistic
    ADD CONSTRAINT user_statistic_pk PRIMARY KEY (id),
    ADD CONSTRAINT user_statistic_user_word UNIQUE (user_id,word_id);