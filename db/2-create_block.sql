CREATE TABLE IF NOT EXISTS public.block
(
    height integer NOT NULL,
    "raw" bytea,
    cons_pubkey bytea,
    block_time timestamp with time zone,
    diff_time integer,
    validator_set bytea,
    CONSTRAINT "Block_pkey" PRIMARY KEY (height)
)

TABLESPACE pg_default;
