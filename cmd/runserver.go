package cmd

import (
	"fmt"

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
		app.Get("/", func(c iris.Context) {

			c.JSON(iris.Map{"lanagege": "golang "})
		})

		app.Listen(fmt.Sprintf("%s:%s", host, port))
	},
}

func init() {
	rootCmd.AddCommand(runserverCmd)

	runserverCmd.Flags().StringVarP(&port, "post", "P", "8000", "enter port number")
	runserverCmd.Flags().StringVarP(&host, "host", "H", "0.0.0.0", "enter host ")

}
