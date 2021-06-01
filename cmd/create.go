package cmd

import (
	"IdentityServer/config"
	"IdentityServer/models"
	"bufio"
	"context"
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
	createClientCmd.Flags().StringVarP(&RedirectUrl, "redirect_url", "r", "", "enter redirect url")

	createUserCmd.Flags().StringVarP(&user.FirstName, "name", "n", "", "enter user name")
	createUserCmd.Flags().StringVarP(&user.Email, "email", "e", "", "enter user name")
	createUserCmd.Flags().StringVarP(&user.Password, "password", "p", "", "enter user name")

}

var (
	client      models.OauthClient
	user        models.User
	RedirectUrl string
)

// createCmd represents the create command
var createClientCmd = &cobra.Command{
	Use:   "client",
	Short: "Create the oauth client",
	Run: func(cmd *cobra.Command, args []string) {

		var err error
		if client.Name == "" {
			fmt.Print("Enter client app name: ")
			reader := bufio.NewReader(os.Stdin)
			client.Name, _ = reader.ReadString('\n')
			client.Name = client.Name[:len(client.Name)-1]

		}

		if client.Description == "" {
			fmt.Print("Enter client app description: ")
			reader := bufio.NewReader(os.Stdin)

			client.Description, _ = reader.ReadString('\n')
			client.Description = client.Description[:len(client.Description)-1]

		}

		if RedirectUrl == "" {
			fmt.Print("Enter client app redirect url: ")
			reader := bufio.NewReader(os.Stdin)
			RedirectUrl, _ = reader.ReadString('\n')
			RedirectUrl = RedirectUrl[:len(RedirectUrl)-1]

		}

		if _, err := url.ParseRequestURI(RedirectUrl); err != nil {
			log.Fatal("url is not valid.")
		}
		client.RedirectURI = []string{RedirectUrl}

		db := config.DB()

		ctx := context.Background()
		_, err = client.Create(ctx, db)
		if err != nil {
			log.Println("Can not insert document")
			log.Fatal(err)
		}

		fmt.Println("Client is successfully created")
		fmt.Println("Client ID or Key is : ", client.Key)
		fmt.Println("Client Secert is : ", client.Secret)
		fmt.Println("Redirect url or callback url after the client login is : ", client.RedirectURI[0])

	},
}

var createUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Create the oauth2 dummy user",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("create called")
	},
}
