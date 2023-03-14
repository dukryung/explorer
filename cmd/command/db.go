package command

import (
	"database/sql"
	"fmt"

	"github.com/hessegg/nikto-explorer/types/config"
	"github.com/hessegg/nikto-explorer/types/log"
	"github.com/spf13/cobra"
)

var tables map[string]string

func init() {
	tables = make(map[string]string)

	tables["block"] = `
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

ALTER TABLE public.block
    OWNER to %v;
`
	tables["transaction"] = `
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

ALTER TABLE public.transaction
    OWNER to %v;
`
	tables["validator"] = `
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

ALTER TABLE public.validator
    OWNER to %v;
`
	tables["uptime"] = `
-- Table: public.uptime

-- DROP TABLE public.uptime;

CREATE TABLE IF NOT EXISTS public.uptime
(
    idx integer NOT NULL GENERATED ALWAYS AS IDENTITY ( INCREMENT 1 START 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1 ),
    height integer,
    cons_address character varying(60) COLLATE pg_catalog."default",
    CONSTRAINT uptime_pkey PRIMARY KEY (idx)
)

TABLESPACE pg_default;

ALTER TABLE public.uptime
    OWNER to %v;
`
	tables["token"] = `
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
)

TABLESPACE pg_default;

ALTER TABLE public.token
    OWNER to %v;
`
	tables["nftoken_base"] = `
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
`
}

func SetUpDB() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "setup klaatoo-explorer database",
		RunE: func(cmd *cobra.Command, args []string) error {
			flags := ReadCommandFlags(cmd.Flags())

			appConfig := config.AppConfig{}
			err := appConfig.LoadConfig(flags.ConfigPath)
			if err != nil {
				return err
			}

			logger := log.NewLogger("db", config.DefaultLogConfig())
			logger.Info("setup database...")

			db, err := sql.Open(appConfig.Sync.DB.DriverName, appConfig.Sync.DB.GetDBInfo())
			if err != nil {
				return err
			}
			defer db.Close()

			for table, query := range tables {
				logger.Info("setup table", table)

				query = fmt.Sprintf(query, appConfig.Sync.DB.User)
				logger.Debug(query)
				_, err = db.Exec(query)
				if err != nil {
					logger.Error("Setup failed: table: ", table, ", error:", err.Error())
					return err
				}
			}

			logger.Info("setup complete!")
			return nil
		},
	}

	return cmd
}

func ResetDB() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "reset klaatoo-explorer database",
		RunE: func(cmd *cobra.Command, args []string) error {
			flags := ReadCommandFlags(cmd.Flags())

			appConfig := config.AppConfig{}
			err := appConfig.LoadConfig(flags.ConfigPath)
			if err != nil {
				return err
			}

			logger := log.NewLogger("db", config.DefaultLogConfig())
			logger.Info("reset database...")

			db, err := sql.Open(appConfig.Sync.DB.DriverName, appConfig.Sync.DB.GetDBInfo())
			if err != nil {
				return err
			}
			defer db.Close()

			for table, _ := range tables {
				logger.Info("reset table", table)

				query := fmt.Sprintf("TRUNCATE %v RESTART IDENTITY", table)
				_, err = db.Exec(query)
				if err != nil {
					return err
				}
			}

			logger.Info("reset complete!")
			return nil
		},
	}
	AddServerFlagsToCmd(cmd)
	return cmd
}
