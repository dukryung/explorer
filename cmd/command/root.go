package command

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "klaatoo-explorer",
		Short: "klaatoo-explorer commands",
	}

	rootCmd.AddCommand(GetClientCMD())
	rootCmd.AddCommand(GetServerCMD())

	return rootCmd
}
