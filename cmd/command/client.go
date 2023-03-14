package command

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/hessegg/nikto-explorer/types/config"
	"github.com/hessegg/nikto-explorer/types/log"
	"github.com/spf13/cobra"
)

const (
	defaultClientPath      = "./client"
	ClientSourceConfigPath = "/src/config.json"
)

func GetClientCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "client",
		Short: "explorer frontend commands",
		Long:  `The commands use Node.js and yarn. Node.js version is  14.4+.`,
	}

	cmd.AddCommand(BuildClient())
	cmd.AddCommand(ServeClient())

	return cmd
}

func ServeClient() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "serve klaatoo-explorer client",
		RunE: func(cmd *cobra.Command, args []string) error {
			flags := ReadCommandFlags(cmd.Flags())
			logger := log.NewLogger("client", config.DefaultLogConfig())
			currentPath, _ := os.Getwd()

			port, err := cmd.Flags().GetString(flagPort)
			if err != nil {
				return err
			}

			if err := copyConfig(flags.ConfigPath, flags.Path); err != nil {
				return err
			}

			if err := os.Chdir(flags.Path); err != nil {
				return err
			}

			logger.Info("klaatoo-explorer: serve client")

			command := exec.Command("yarn")
			if err := execCommand(command); err != nil {
				return err
			}

			serve := []string{
				"run",
				"start",
				"--port",
				port,
			}
			command = exec.Command("yarn", serve...)
			if err := execCommand(command); err != nil {
				return err
			}

			if err := os.Chdir(currentPath); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().String(flagPort, "8080", "set serve port")
	AddClientFlagsToCmd(cmd)
	return cmd
}

func BuildClient() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "build klaatoo-explorer client",
		RunE: func(cmd *cobra.Command, args []string) error {
			flags := ReadCommandFlags(cmd.Flags())
			currentPath, _ := os.Getwd()
			logger := log.NewLogger("client", config.DefaultLogConfig())
			if err := copyConfig(flags.ConfigPath, flags.Path); err != nil {
				return err
			}

			if err := os.Chdir(flags.Path); err != nil {
				return err
			}

			logger.Info("klaatoo-explorer: build client")

			command := exec.Command("yarn")
			if err := execCommand(command); err != nil {
				return err
			}

			command = exec.Command("yarn", "run", "build")
			if err := execCommand(command); err != nil {
				return err
			}

			if err := os.Chdir(currentPath); err != nil {
				return err
			}
			return nil
		},
	}

	AddClientFlagsToCmd(cmd)
	return cmd
}

func execCommand(cmd *exec.Cmd) error {
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	reader := bufio.NewReader(pipe)
	line, err := reader.ReadString('\n')
	for err == nil {
		fmt.Print(line)
		line, err = reader.ReadString('\n')
	}

	return nil
}

func copyConfig(config, client string) error {
	src := config
	dest := path.Join(client, ClientSourceConfigPath)

	bz, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dest, bz, 0644)
	if err != nil {
		return err
	}

	return nil
}
