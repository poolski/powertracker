package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/Songmu/prompter"
	"github.com/poolski/powertracker/cmd/client"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	days     int
	output   string
	csvFile  string
	insecure bool
)

var rootCmd = &cobra.Command{
	Use:   "powertracker",
	Short: "Queries HomeAssistant for a summary of your power usage over a period of time",
	Long: `
	Home Assistant provides power usage data for each hour on its Energy dashboard, but does not have an API endpoint to query this data.
	The websocket API does provide a way to do this, which is what the frontend uses.
	This tool queries the websocket API to get the power usage data for each hour over a period of time, and then prints a summary of the data in a table.
	It also saves the data to a CSV file in the current directory.`,

	Run: func(cmd *cobra.Command, args []string) {
		c := client.New(client.Config{
			Days:     days,
			Output:   output,
			FilePath: csvFile,
			Insecure: insecure,
		})
		if err := c.Connect(); err != nil {
			log.Fatal().Msgf("connecting to websocket: %s", err.Error())
		}
		c.ComputePowerStats()
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

		rootCmd.PersistentFlags().IntVarP(&days, "days", "d", 30, "number of days to compute power stats for")
		rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output format (text, table, csv)")
		rootCmd.PersistentFlags().StringVarP(&csvFile, "csv-file", "f", "results.csv", "the path of the CSV file to write to")
		rootCmd.PersistentFlags().BoolVarP(&insecure, "insecure", "i", false, "skip TLS verification")
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

		err := promtUserConfig()
		if err != nil {
			log.Fatal().Msgf("prompting user for config: %s", err.Error())
		}

		err = os.MkdirAll(filepath.Dir(cfgFile), 0755)
		if err != nil {
			log.Fatal().Msgf("creating config dir: %s", err.Error())
		}

		f, err := os.Create(cfgFile)
		if err != nil {
			log.Fatal().Msgf("creating config file: %s", err.Error())
		}
		defer f.Close()

		if err := viper.WriteConfigAs(cfgFile); err != nil {
			log.Fatal().Msgf("writing config file: %s", err.Error())
		}
	}
}

func promtUserConfig() error {
	urlPrompt := prompter.Prompt("Home Assistant URL - e.g. http://localhost:8123", "")
	token := prompter.Password("Home Assistant Long-Lived Access Token")
	sensorID := prompter.Prompt("Power sensor entity ID - e.g. sensor.power", "")

	haURL, err := url.Parse(urlPrompt)
	if haURL.Scheme == "" {
		haURL.Scheme = "http"
	}
	if err != nil {
		return fmt.Errorf("parsing URL: %w", err)
	}

	viper.Set("api_key", token)
	viper.Set("url", haURL.String())
	viper.Set("sensor_id", sensorID)
	return nil
}
