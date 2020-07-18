package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/patilsuraj767/connection-manager/config"

	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete SSH connection",
	Long:  "Delete server from the database",
	Run: func(cmd *cobra.Command, args []string) {

		hostname, _ := cmd.Flags().GetString("hostname")

		if hostname != "" {
			config.DeleteServerFromDB(hostname)
		} else {
			servers := config.GetAllServers()
			prompt := promptui.Select{
				Label: "Delete Server From Database",
				Items: servers,
				Size:  20,
			}

			_, result, err := prompt.Run()

			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}
			config.DeleteServerFromDB(result)

		}

	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().String("hostname", "", "Hostname of the server (Required)")
}
