/*
Copyright © 2022 anton7r
*/
package cmd

import (
	"os"

	"github.com/anton7r/mgx/migrator"
	"github.com/jackc/pgx/v4"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mgx",
	Short: "mgx is used to run migrations to database",
	Long:  `mgx is a database migrations library for postgres databases`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "mgx create <name> - create migration file",
	Long:  "Optionally you can define subpaths that the migration files should be added under",
	Run: func(cmd *cobra.Command, args []string) {
		argumentLen := len(args)

		if argumentLen == 0 {
			cmd.PrintErr("No argument for the file name was passed")
			return
		}

		if argumentLen > 1 {
			cmd.PrintErr("Too many argumnets passed, only one is needed")
			return
		}

		filepath := args[0]
		migrator.CreateNewMigration(filepath)
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "mgx migrate <version> - migrates to the migration version defined",
	Long:  "Can be both used to upgrade upwards and downwards",
	Run: func(cmd *cobra.Command, args []string) {
		argumentLen := len(args)

		if argumentLen == 0 {
			cmd.PrintErr("No argument for the migration version was passed")
			return
		}

		if argumentLen > 1 {
			cmd.PrintErr("Too many arguments passed, only one is needed")
			return
		}

		dsnFlag := cmd.Flag("dsn")
		urlFlag := cmd.Flag("url")

		if dsnFlag != nil || urlFlag != nil {
			cmd.PrintErr("Both 'dsn' and 'url' were defined, only one of them is needed")
			return
		}

		var connection *pgx.Conn
		var err error

		if dsnFlag != nil {
			dsnValue := dsnFlag.Value.String()
			connection, err = migrator.ConnectDSN(dsnValue)

			if err != nil {
				cmd.PrintErr("Error while trying to connect to database with dsn '" + dsnValue + "'")
				return
			}
		}

		if urlFlag != nil {
			urlValue := urlFlag.Value.String()
			connection, err = migrator.ConnectURL(urlValue)

			if err != nil {
				cmd.PrintErr("Error while trying to connect to database with dsn '" + urlValue + "'")
				return
			}
		}

		migrationVer := args[0]

		err = migrator.Migrate(connection, migrationVer)
		if err != nil {
			cmd.PrintErr("Error while trying to migrate to '" + migrationVer + "'")
		}
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "mgx conifg - create configuration file",
	Long: `mgx config - creates configuration file
	
Not mandatory, but needed incase you want
to configure the folder that the migrations go
or if you want to configure any settings etc...,
but by default you should not need any configurations 
	`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(migrateCmd)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mgx.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	migrateCmd.Flags().StringP("dsn", "D", "",
		"dsn string contains relevant information regarding making database connections alternatively the 'url' flag can be used")
	migrateCmd.Flags().StringP("url", "U", "",
		"url contains relevant information regarding making database connections alternatively the 'dsn' flag can be used")
}
