package cmd

import (
	"github.com/patilsuraj767/connection-manager/config"

	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add SSH connection",
	Long:  `Add server to the database`,
	Run: func(cmd *cobra.Command, args []string) {

		var host config.Server
		host.Name, _ = cmd.Flags().GetString("name")
		host.Address, _ = cmd.Flags().GetString("address")
		host.Username, _ = cmd.Flags().GetString("username")
		host.Password, _ = cmd.Flags().GetString("password")

		if host.Name != "" && host.Address != "" && host.Username != "" {
			config.AddServerToDB(host)
		} else {
			app := tview.NewApplication()
			form := tview.NewForm().
				AddInputField("Name to profile", "", 50, nil, func(text string) {
					host.Name = text
				}).
				AddInputField("Server IPaddress or hostname", "", 20, nil, func(text string) {
					host.Address = text
				}).
				AddInputField("Username", "", 20, nil, func(text string) {
					host.Username = text
				}).
				AddInputField("Password", "", 20, nil, func(text string) {
					host.Password = text
				}).
				AddButton("Save", func() {
					app.Stop()
					config.AddServerToDB(host)
				}).
				AddButton("Quit", func() {
					app.Stop()
				})
			form.SetBorder(true).SetTitle("Add Server details").SetTitleAlign(tview.AlignLeft)
			if err := app.SetRoot(form, true).SetFocus(form).Run(); err != nil {
				panic(err)
			}

		}

	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().String("name", "", "Name to profile ")
	addCmd.Flags().String("address", "", "IP or hostname of the server")
	addCmd.Flags().String("username", "", "SSH username ")
	addCmd.Flags().String("password", "", "SSH password ")

}
