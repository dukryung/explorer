package command

import (
	"github.com/hessegg/nikto-explorer/types/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

const (
	flagPath       = "path"
	flagConfigPath = "config"
	flagPort       = "port"
)

type Flags struct {
	Path       string
	ConfigPath string
}

func AddClientFlagsToCmd(command *cobra.Command) {
	command.Flags().StringP(flagPath, "p", defaultClientPath, "set explorer client root directory.")
	command.Flags().StringP(flagConfigPath, "c", config.ClientConfigPath, "set config path")
}

func AddServerFlagsToCmd(command *cobra.Command) {
	command.Flags().StringP(flagConfigPath, "c", config.ServerConfigPath, "set config path")
}

func ReadCommandFlags(flagSet *pflag.FlagSet) Flags {
	flags := Flags{}

	flags.Path, _ = flagSet.GetString(flagPath)
	flags.ConfigPath, _ = flagSet.GetString(flagConfigPath)

	return flags
}
