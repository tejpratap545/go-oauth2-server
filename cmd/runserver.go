package cmd

import (
	"IdentityServer/route"

	"github.com/kataras/iris/v12"
	"github.com/spf13/cobra"
)

var port, host string

// runserverCmd represents the runserver command
var runserverCmd = &cobra.Command{
	Use:   "runserver",
	Short: "A brief description of your command",

	Run: func(cmd *cobra.Command, args []string) {
		app := iris.New()

		app.RegisterView(iris.HTML("./web", ".html").Layout("layout.html").Reload(true))

		app.PartyFunc("/", route.Route)

		// app.Listen(iris.AutoTLS(fmt.Sprintf("%s:%s", host, port), "localhost", "tejpratapsingh545@gmail.com"))
		addr := ":8000"
		app.Run(iris.TLS(addr, "./certificates/localhost.crt", "./certificates/localhost.key"))
	},
}

func init() {
	rootCmd.AddCommand(runserverCmd)

	runserverCmd.Flags().StringVarP(&port, "post", "P", "8000", "enter port number")
	runserverCmd.Flags().StringVarP(&host, "host", "H", "localhost", "enter host ")

}
