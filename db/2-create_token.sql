DROP TABLE IF EXISTS public.token;

CREATE TABLE IF NOT EXISTS public.token
(
    idx integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    owner_address character varying COLLATE pg_catalog."default",
    symbol character varying COLLATE pg_catalog."default",
    description character varying COLLATE pg_catalog."default",
    denom character varying COLLATE pg_catalog."default" NOT NULL,
    "precision" bigint,
    amount character varying COLLATE pg_catalog."default",
    CONSTRAINT token_pkey PRIMARY KEY (denom)
);

TABLESPACE pg_default;