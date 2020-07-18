package cmd

import (
	"fmt"

	"github.com/manifoldco/promptui"
	"github.com/patilsuraj767/connection-manager/config"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit SSH connection",
	Long:  "Edit server details in the database",
	Run: func(cmd *cobra.Command, args []string) {
		servers := config.GetAllServers()
		prompt := promptui.Select{
			Label: "Select Server",
			Items: servers,
			Size:  20,
		}

		_, result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}
		host := config.GetDetailOfSpecificServer(result)

		app := tview.NewApplication()
		form := tview.NewForm().
			AddInputField("Name to profile", result, 50, nil, func(text string) {
				host.Name = text
			}).
			AddInputField("Server IPaddress or hostname", host.Address, 20, nil, func(text string) {
				host.Address = text
			}).
			AddInputField("Username", host.Username, 20, nil, func(text string) {
				host.Username = text
			}).
			AddInputField("Password", host.Password, 20, nil, func(text string) {
				host.Password = text
			}).
			AddButton("Save", func() {
				app.Stop()
				config.UpdateHost(host)
			}).
			AddButton("Quit", func() {
				app.Stop()
			})
		form.SetBorder(true).SetTitle("Edit Server details").SetTitleAlign(tview.AlignLeft)
		if err := app.SetRoot(form, true).SetFocus(form).Run(); err != nil {
			panic(err)
		}

	},
}

func init() {
	rootCmd.AddCommand(editCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// editCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// editCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
