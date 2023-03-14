package command

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hessegg/nikto-explorer/server/app"
	"github.com/spf13/cobra"
)

func GetServerCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "klaatoo-explorer server commands",
	}

	cmd.AddCommand(StartServer())
	cmd.AddCommand(SetUpDB())
	cmd.AddCommand(ResetDB())

	return cmd
}

func StartServer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "start klaatoo-explorer server",
		RunE: func(cmd *cobra.Command, args []string) error {
			flags := ReadCommandFlags(cmd.Flags())

			quit := make(chan os.Signal)
			signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

			server := app.NewApp(flags.ConfigPath)
			err := server.RunServers()
			if err != nil {
				panic(err)
			}

			<-quit

			server.CloseServers()
			return nil
		},
	}
	AddServerFlagsToCmd(cmd)
	return cmd
}
