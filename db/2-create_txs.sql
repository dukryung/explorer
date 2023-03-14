CREATE TABLE IF NOT EXISTS public.transaction
(
    idx integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    height integer,
    txbody bytea,
    txhash character varying(64) COLLATE pg_catalog."default" NOT NULL,
    updated boolean NOT NULL DEFAULT false,
    code integer,
    "timestamp" timestamp with time zone,
    sender character varying(46) COLLATE pg_catalog."default",
    CONSTRAINT transaction_pkey PRIMARY KEY (txhash)
)

TABLESPACE pg_default;