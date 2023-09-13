# Power Consumption Metrics Tool

Queries HomeAssistant for a summary of your power usage over a period of time.
Home Assistant provides power usage data for each hour on its Energy dashboard, but does not have an API endpoint to query this data.
The websocket API does provide a way to do this, which is what the frontend uses.

This tool queries the websocket API to get the power usage data for each hour over a period of days and outputs the data in various formats

## Installation

```bash
go install github.com/poolski/powertracker@latest
```

## Configuration

This tool requires a configuration file to be present at `~/.config/powertracker/config.yaml`. If one does not exist, it will ask for input and create it for you.
The only things this tool needs are the URL of your Home Assistant instance and a long-lived access token.

You can generate a long-lived access token by going to your Home Assistant instance, clicking on your profile picture in the bottom left, then clicking on "Long-Lived Access Tokens" at the bottom of the list and creating a new one.

## Usage

```bash
$ powertracker --help

Usage:
  powertracker [flags]

Flags:
  -c, --config string     config file (default "/Users/kyrill/.config/powertracker/config.yaml")
  -f, --csv-file string   the path of the CSV file to write to (default "results.csv")
  -d, --days int          number of days to compute power stats for (default 30)
  -h, --help              help for powertracker
  -o, --output string     output format (text, table, csv)
```
