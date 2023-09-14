package cmd

import (
	"crypto/tls"
	"encoding/csv"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/olekukonko/tablewriter"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Client struct {
	Conn      *websocket.Conn
	MessageID int // MessageID is the ID of the message sent to the websocket. These must be incremented with each subsequent request
}

// APIResponse represents the structure of the response received from the Home Assistant API.
type APIResponse struct {
	ID      int    `json:"id"`      // ID is the unique identifier of the response.
	Type    string `json:"type"`    // Type is the type of the response.
	Success bool   `json:"success"` // Success indicates whether the response was successful or not.
	Result  struct {
		SmartMeterElectricityImport []struct {
			Start  int     `json:"start"`
			End    int     `json:"end"`
			Change float64 `json:"change"` // Change is the amount of electricity imported.
		} `json:"sensor.smart_meter_electricity_import_2"`
	} `json:"result"` // Result contains the data returned by the API.
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

var (
	days     int
	output   string
	csvFile  string
	insecure bool
)

const hoursInADay = 24

func init() {
	rootCmd.PersistentFlags().IntVarP(&days, "days", "d", 30, "number of days to compute power stats for")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output format (text, table, csv)")
	rootCmd.PersistentFlags().StringVarP(&csvFile, "csv-file", "f", "results.csv", "the path of the CSV file to write to")
}

func (c *Client) Connect() error {
	c.MessageID = 1

	// Set up the websocket dialer
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	// Work out the URL to dial
	if viper.GetString("url") == "" {
		return fmt.Errorf("url is required")
	}
	dialURL, err := url.Parse(viper.GetString("url"))
	if err != nil {
		return err
	}
	if dialURL.Scheme == "http" {
		dialURL.Scheme = "ws"
	} else if dialURL.Scheme == "https" {
		dialURL.Scheme = "wss"
	}
	dialURL.Path = "/api/websocket"

	// Set up the TLS config
	// If Home Assistant is running with a self-signed cert, then we need to skip verification
	_ = &tls.Config{
		InsecureSkipVerify: true,
	}

	// Dial the websocket
	log.Info().Msgf("connecting to %s", dialURL.String())
	conn, _, err := dialer.Dial(dialURL.String(), nil)
	if err != nil {
		return err
	}

	// Read the initial message
	var initMsg map[string]any
	if err := conn.ReadJSON(&initMsg); err != nil {
		return err
	}

	// Send the authentication message
	if err := conn.WriteJSON(map[string]string{
		"type":         "auth",
		"access_token": viper.GetString("api_key"),
	}); err != nil {
		return err
	}

	// Read the authentication response
	var authResp map[string]any
	if err := conn.ReadJSON(&authResp); err != nil {
		return err
	}
	if authResp["type"] != "auth_ok" {
		return fmt.Errorf("authentication failed: %v", authResp["message"])
	}

	c.Conn = conn
	return nil
}

// computePowerStats computes the power statistics for a given number of days and hours.
// It prints a table to stdout where the rows are "days" and the columns are "hours".
// The function writes the results to a CSV file and prints the averages to the console.
func (c *Client) computePowerStats() {
	results, err := getResults(c)
	if err != nil {
		log.Error().Msg(fmt.Sprintf("getting results: %v", err))
		return
	}

	// Compute averages
	averages := make([]float64, hoursInADay)
	for i := range averages {
		sum := 0.0
		for j := range results {
			sum += results[j][i]
		}
		averages[i] = sum / float64(days)
	}

	// Generate column headers for table/CSV
	headers := make([]string, hoursInADay)
	for i := range headers {
		headers[i] = fmt.Sprintf("%d", i)
	}

	switch output {
	case "text":
		writePlainText(averages)
	case "table":
		printTable(results, averages, headers)
	case "csv":
		err = writeCSVFile(headers, results, averages)
		if err != nil {
			log.Error().Msg(fmt.Sprintf("writing CSV file: %v", err))
			return
		}
	default:
		printTable(results, averages, headers)
	}
}

// writePlainText prints the results to stdout in plain text.
// This is useful for using with something like https://garydoessolar.com/utilities/dailymodellingutility/
// You can copy and paste the results into the custom usage pattern section and it will generate more accurate predictions.
func writePlainText(averages []float64) {
	for _, v := range averages {
		fmt.Printf("%f,\n", v)
	}
}

func writeCSVFile(headers []string, results [][]float64, averages []float64) error {
	f, err := os.Create(csvFile)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	err = writer.Write(headers)
	if err != nil {
		return fmt.Errorf("writing headers: %w", err)
	}

	for _, row := range results {
		rowString := make([]string, len(row))
		for j, val := range row {
			rowString[j] = fmt.Sprintf("%f", val)
		}
		err = writer.Write(rowString)
		if err != nil {
			return fmt.Errorf("writing row: %w", err)
		}
	}

	averageString := make([]string, len(averages))
	for i, val := range averages {
		averageString[i] = fmt.Sprintf("%f", val)
	}
	err = writer.Write(averageString)
	if err != nil {
		return fmt.Errorf("writing averages: %w", err)
	}

	writer.Flush()

	return nil
}

func printTable(results [][]float64, averages []float64, headers []string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)

	for _, row := range results {
		rowString := make([]string, len(row))
		for j, val := range row {
			rowString[j] = fmt.Sprintf("%f", val)
		}
		table.Append(rowString)
	}

	averageString := make([]string, len(averages))
	for i, val := range averages {
		averageString[i] = fmt.Sprintf("%f", val)
	}
	table.SetFooter(averageString)
	table.Render()
}

func getResults(c *Client) ([][]float64, error) {
	// We're going to store the results in a slice of slices, where each slice is a day's worth of data
	// In other words, we're creating a table where the rows are "days" and the columns are "hours"
	// This is a bit of a hack, but it works.

	// What we're doing is creating an offset from the current *day* based on a multiple of
	// 24 hours, each time we iterate through the a "row" of the results slice.
	results := make([][]float64, days)

	for i := range results {
		c.MessageID++

		offset := time.Duration((i+1)*24) * time.Hour
		start := time.Now().Add(-offset).Truncate(24 * time.Hour).Format("2006-01-02T15:04:05.000Z")
		msg := map[string]interface{}{
			"id":            c.MessageID,
			"type":          "recorder/statistics_during_period",
			"start_time":    start,
			"end_time":      time.Now().Truncate(24 * time.Hour).Format("2006-01-02T15:04:05.000Z"),
			"statistic_ids": []string{"sensor.smart_meter_electricity_import_2"},
			"period":        "hour",
			"types":         []string{"change"},
			"units": map[string]string{
				"energy": "kWh",
			},
		}

		if err := c.write(msg); err != nil {
			return nil, fmt.Errorf("writing to websocket: %w", err)
		}

		var data APIResponse
		err := c.Conn.ReadJSON(&data)
		if err != nil {
			return nil, fmt.Errorf("reading from websocket: %w", err)
		}

		if !data.Success {
			return nil, fmt.Errorf("api response error: %v", data.Error)
		}
		changeSlice := make([]float64, hoursInADay)
		for j := range changeSlice {
			changeSlice[j] = data.Result.SmartMeterElectricityImport[j].Change
		}
		results[i] = changeSlice
	}
	return results, nil
}

func (c *Client) write(data map[string]interface{}) error {
	return c.Conn.WriteJSON(data)
}
