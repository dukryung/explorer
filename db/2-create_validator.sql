CREATE TABLE IF NOT EXISTS public.validator
(
    val_address character varying(53) COLLATE pg_catalog."default",
    cons_address character varying(53) COLLATE pg_catalog."default",
    cons_pubkey bytea NOT NULL,
    moniker character varying(100) COLLATE pg_catalog."default",
    "raw" bytea,
    tokens character varying COLLATE pg_catalog."default",
    rank bigint,
    CONSTRAINT validator_pkey PRIMARY KEY (cons_pubkey)
)

TABLESPACE pg_default;
