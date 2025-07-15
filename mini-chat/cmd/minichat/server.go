package cmd

import (
	"fmt"

	"github.com/shahin-bayat/mini-chat/internal/config"
	"github.com/spf13/cobra"
)

var serverCfg config.ServerConfig

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the chat relay server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := serverCfg.Validate(); err != nil {
			return err
		}
		fmt.Printf("Starting server on %s:%d\n", clientCfg.Host, clientCfg.Port)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVar(&clientCfg.Host, "host", "0.0.0.0", "host to bind")
	serverCmd.Flags().IntVar(&clientCfg.Port, "port", 0, "port to listen on")
}
