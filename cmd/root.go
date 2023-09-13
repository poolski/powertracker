package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "powertracker",
	Short: "Queries HomeAssistant for a summary of your power usage over a period of time",
	Long: `
	Home Assistant provides power usage data for each hour on its Energy dashboard, but does not have an API endpoint to query this data.
	The websocket API does provide a way to do this, which is what the frontend uses.
	This tool queries the websocket API to get the power usage data for each hour over a period of time, and then prints a summary of the data in a table.
	It also saves the data to a CSV file in the current directory.`,

	Run: func(cmd *cobra.Command, args []string) {
		c := Client{}
		if err := c.Connect(); err != nil {
			log.Fatal().Msgf("connecting to websocket: %s", err.Error())
		}
		c.computePowerStats()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		sep := string(filepath.Separator)
		confDir := home + sep + ".config"
		rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", confDir+"/powertracker/config.yaml", "config file")
	}
}

// Setup configuration
func initConfig() {
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Err(err).Msg("reading config file")
	}

	// If a config file doesn't exist, prompt the user for a first-time setup
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		fmt.Println("No config file found. Let's set one up.")

		fmt.Println("Please enter your API key:")
		var apiKey string
		fmt.Scanln(&apiKey)

		fmt.Println("Please enter the URL for your Home Assistant instance (e.g. http://localhost:8123):")
		var baseURL string
		fmt.Scanln(&baseURL)

		viper.Set("api_key", apiKey)
		viper.Set("url", baseURL)

		err := os.MkdirAll(filepath.Dir(cfgFile), 0755)
		if err != nil {
			log.Err(err).Msg("creating config directory")
		}

		f, err := os.Create(cfgFile)
		if err != nil {
			log.Err(err).Msg("creating file")
		}
		defer f.Close()

		if err := viper.WriteConfigAs(cfgFile); err != nil {
			log.Err(err).Msg("writing config file")
		}
	}
}
