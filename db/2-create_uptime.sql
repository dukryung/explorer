CREATE TABLE IF NOT EXISTS public.uptime
(
    idx integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    height integer,
    cons_address character varying(60) COLLATE pg_catalog."default",
    CONSTRAINT uptime_pkey PRIMARY KEY (idx)
)

TABLESPACE pg_default;