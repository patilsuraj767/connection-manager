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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
