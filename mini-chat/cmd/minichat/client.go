package cmd

import (
	"fmt"

	"github.com/shahin-bayat/mini-chat/internal/config"
	"github.com/spf13/cobra"
)

var clientCfg config.ClientConfig

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Run the interactive chat client",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := clientCfg.Validate(); err != nil {
			return err
		}
		fmt.Printf("Connecting to %s:%d as %s\n",
			clientCfg.Host, clientCfg.Port, clientCfg.User)
		// client := networking.NewClient(clientCfg.Host, clientCfg.Port)
		// if err := client.Connect(); err != nil {
		// 	return err
		// }
		// defer client.Close()
		// TODO: hand off to TUI or REPL
		return nil
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.Flags().StringVar(&clientCfg.Host, "host", "", "server host to connect")
	clientCmd.Flags().IntVar(&clientCfg.Port, "port", 0, "server port to connect")
	clientCmd.Flags().StringVar(&clientCfg.User, "user", "", "your chat username")
}
