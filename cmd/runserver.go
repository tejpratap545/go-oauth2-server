package cmd

import (
	"IdentityServer/config"
	"IdentityServer/route"
	"fmt"

	"IdentityServer/middleware"

	"github.com/gorilla/securecookie"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/sessions"
	"github.com/spf13/cobra"
)

var port, host string

// runserverCmd represents the runserver command
var runserverCmd = &cobra.Command{
	Use:   "runserver",
	Short: "A brief description of your command",

	Run: func(cmd *cobra.Command, args []string) {

		hashKey := securecookie.GenerateRandomKey(64)
		blockKey := securecookie.GenerateRandomKey(32)
		s := securecookie.New(hashKey, blockKey)
		sess := sessions.New(sessions.Config{
			Cookie:          "_session_id",
			Expires:         0,
			AllowReclaim:    true,
			CookieSecureTLS: true,
			Encoding:        s,
		})

		app := iris.Default()
		app.Use(sess.Handler())

		sess.UseDatabase(config.Redis())

		app.Use(middleware.AddMongoToContext)

		app.RegisterView(iris.HTML("./web", ".html").Layout("layout.html").Reload(true))

		app.PartyFunc("/", route.Route)

		// app.Listen(iris.AutoTLS(fmt.Sprintf("%s:%s", host, port), "localhost", "tejpratapsingh545@gmail.com"))
		addr := fmt.Sprintf("%s:%s", host, port)
		app.Run(iris.TLS(addr, "./certificates/localhost.crt", "./certificates/localhost.key"))
	},
}

func init() {
	rootCmd.AddCommand(runserverCmd)

	runserverCmd.Flags().StringVarP(&port, "post", "P", "8000", "enter port number")
	runserverCmd.Flags().StringVarP(&host, "host", "H", "localhost", "enter host ")

}
