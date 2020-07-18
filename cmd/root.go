package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/atotto/clipboard"
	"github.com/manifoldco/promptui"
	"github.com/patilsuraj767/connection-manager/config"
	"github.com/rivo/tview"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "connection-manager",
	Short: "",
	Long: `Connection-manager is the CLI tool for managing ssh connections.
	It help in storing the servers ipaddress/hostname, username and password, 
	so that you can ssh the system in just one click.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
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

		decision := openview()
		copyOrssh(&host, &decision)
	},
}

func copyOrssh(host *config.Server, decision *string) {
	if *decision == "Copy" {
		err := clipboard.WriteAll(host.Address)
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	} else if *decision == "SSH" {
		takeSSH(host)
	}
}

func openview() string {
	var decision string
	app := tview.NewApplication()
	modal := tview.NewModal().
		SetText("Do you want to copy the hostname or ssh the system").
		AddButtons([]string{"Copy", "SSH"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			decision = buttonLabel
			app.Stop()
		})
	if err := app.SetRoot(modal, false).Run(); err != nil {
		panic(err)
	}
	return decision
}

func takeSSH(host *config.Server) {
	var count int = numberofsessions()
	for i := 1; i <= count; i++ {
		err := exec.Command("bash", "-c", "gnome-terminal --tab --active --maximize -- sshpass -p "+host.Password+" ssh -oStrictHostKeyChecking=no "+host.Username+"@"+host.Address+"").Start()
		if err != nil {
			fmt.Println(err)
			os.Exit(2)
		}
	}
}

func numberofsessions() int {
	var count string = "1"
	app := tview.NewApplication()
	form := tview.NewForm().
		AddInputField("Number of SSH Sessions?", count, 50, nil, func(text string) {
			fmt.Println(text)
			count = text
		}).
		AddButton("SSH", func() {
			app.Stop()
		})
	form.SetBorder(true).SetTitle("Session details").SetTitleAlign(tview.AlignLeft)
	if err := app.SetRoot(form, true).SetFocus(form).Run(); err != nil {
		panic(err)
	}

	i, err := strconv.Atoi(count)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	if i > 5 {
		fmt.Println("Error: Maximum 5 sessions can be opened")
		os.Exit(2)
	}
	return i
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
}
