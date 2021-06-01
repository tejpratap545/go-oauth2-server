package cmd

import (
	"IdentityServer/config"
	"IdentityServer/models"
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create",
}

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.AddCommand(createClientCmd)
	createCmd.AddCommand(createUserCmd)

	createClientCmd.Flags().StringVarP(&client.Name, "name", "n", "", "enter client name")
	createClientCmd.Flags().StringVarP(&client.Description, "description", "d", "", "enter description")
	createClientCmd.Flags().StringVarP(&client.RedirectURI, "redirect_url", "r", "", "enter redirect url")

}

var client models.OauthClient

// createCmd represents the create command
var createClientCmd = &cobra.Command{
	Use:   "client",
	Short: "Create the oauth client",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)
		var err error
		if client.Name == "" {

			fmt.Print("Enter client app name: ")
			client.Name, err = reader.ReadString('\n')
			client.Name = client.Name[:len(client.Name)-1]
			if err != nil {
				log.Fatal("Can not read name")
			}

		}

		if client.Description == "" {

			fmt.Print("Enter client app description: ")
			client.Description, err = reader.ReadString('\n')
			client.Description = client.Description[:len(client.Description)-1]
			if err != nil {
				log.Fatal("Can not read description")
			}

		}

		if client.RedirectURI == "" {

			fmt.Print("Enter client app redirect url: ")
			client.RedirectURI, err = reader.ReadString('\n')
			client.RedirectURI = client.RedirectURI[:len(client.RedirectURI)-1]
			if err != nil {
				log.Fatal("Can not read redirect url")
			}

		}

		if _, err := url.ParseRequestURI(client.RedirectURI); err != nil {
			log.Fatal("url is not valid.")
		}

		db := config.DB()
		_, err = client.Create(db)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Client is successfully created")
		fmt.Println("Client ID or Key is : ", client.Key)
		fmt.Println("Client Secert is : ", client.Secret)
		fmt.Println("Redirect url or callback url after the client login is : ", client.RedirectURI)

	},
}

var createUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Create the oauth2 dummy user",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
	},
}
