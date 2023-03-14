DROP TABLE IF EXISTS public.nftoken_base CASCADE;

CREATE TABLE IF NOT EXISTS public.nftoken_base
(
    idx integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    owner_address character varying COLLATE pg_catalog."default",
    issuer_address character varying COLLATE pg_catalog."default",
    delegated_to_address character varying COLLATE pg_catalog."default",
    block_height integer,
    burnt boolean DEFAULT false,

    id character varying COLLATE pg_catalog."default" NOT NULL UNIQUE,
    url character varying COLLATE pg_catalog."default", 
    hash character varying COLLATE pg_catalog."default", 
    info_url character varying COLLATE pg_catalog."default", 
    info_hash character varying COLLATE pg_catalog."default", 
    preview_url character varying COLLATE pg_catalog."default", 
    name character varying COLLATE pg_catalog."default", 
    collection_id  character varying COLLATE pg_catalog."default",
    CONSTRAINT nftoken_pkey PRIMARY KEY (id)
)
	
TABLESPACE pg_default;

ALTER TABLE public.nftoken_base
    OWNER to %v;

DROP VIEW IF EXISTS public.nftoken;

CREATE VIEW public.nftoken AS 
    SELECT * from public.nftoken_base 
    WHERE burnt = false;

DROP INDEX IF EXISTS nft_owner, nft_issuer, nft_delegated_to, nft_burnt;

CREATE INDEX nft_owner ON public.nftoken_base (owner_address) TABLESPACE pg_default;
CREATE INDEX nft_issuer ON public.nftoken_base (issuer_address) TABLESPACE pg_default;
CREATE INDEX nft_delegated_to ON public.nftoken_base (delegated_to_address) TABLESPACE pg_default;
CREATE INDEX nft_burnt ON public.nftoken_base (burnt) TABLESPACE pg_default;