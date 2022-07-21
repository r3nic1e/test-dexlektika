package cmd

import (
	"log"

	application "github.com/r3nic1e/test-dexlektika/pkg/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		addr := viper.GetString("listen-addr")
		dbAddr := viper.GetString("db-addr")

		app := application.NewApp()
		err := app.ConnectDB(dbAddr)
		if err != nil {
			log.Fatal(err)
		}

		log.Fatal(app.StartServer(addr))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().String("listen-addr", "localhost:8080", "HTTP listen address")
	serveCmd.Flags().String("db-addr", "", "PostgreSQL DSN (check out https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL for format)")
}
