# Power Consumption Metrics Tool

Queries HomeAssistant for a summary of your power usage over a period of time.
Home Assistant provides power usage data for each hour on its Energy dashboard, but does not have an API endpoint to query this data.
The websocket API does provide a way to do this, which is what the frontend uses.

I got fed up of trying to figure out how to get the same data that the Energy dashboard shows, so I wrote this tool to do it for me.

This tool queries the websocket API to get the power usage data for each hour over a period of days and outputs the data in various formats

## Installation

Go

```bash
go install github.com/poolski/powertracker@latest
```

[Releases](https://github.com/poolski/powertracker/releases)

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
  -c, --config string     config file (default "$HOME_DIR/.config/powertracker/config.yaml")
  -f, --csv-file string   the path of the CSV file to write to (default "results.csv")
  -d, --days int          number of days to compute power stats for (default 30)
  -h, --help              help for powertracker
  -i  --insecure          skip TLS verification
  -o, --output string     output format (text, table, csv)

```

## Example output

```bash
$ powertracker -d 7 # 7 days' worth of data

+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+
|    0     |    1     |    2     |    3     |    4     |    5     |    6     |    7     |    8     |    9     |    10    |    11    |    12    |    13    |    14    |    15    |    16    |    17    |    18    |    19    |    20    |    21    |    22    |    23    |
+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+
| 0.300000 | 0.326000 | 0.333000 | 0.298000 | 0.397000 | 0.554000 | 0.408000 | 0.519000 | 0.552000 | 0.761000 | 0.591000 | 0.564000 | 0.880000 | 0.584000 | 0.636000 | 0.540000 | 1.204000 | 1.272000 | 1.011000 | 0.991000 | 0.386000 | 0.420000 | 0.277000 | 0.376000 |
| 0.374000 | 0.338000 | 0.352000 | 0.361000 | 0.386000 | 0.596000 | 0.499000 | 0.662000 | 0.837000 | 0.643000 | 0.819000 | 0.865000 | 0.680000 | 0.612000 | 0.570000 | 0.793000 | 1.350000 | 1.141000 | 1.179000 | 1.048000 | 0.621000 | 0.422000 | 0.277000 | 0.361000 |
| 0.368000 | 0.442000 | 0.338000 | 0.451000 | 0.349000 | 0.663000 | 1.645000 | 0.655000 | 0.672000 | 0.793000 | 0.577000 | 0.790000 | 0.820000 | 0.529000 | 0.682000 | 0.485000 | 1.827000 | 0.929000 | 0.779000 | 0.973000 | 0.606000 | 0.928000 | 0.338000 | 0.374000 |
| 0.354000 | 0.432000 | 0.390000 | 0.390000 | 0.613000 | 0.827000 | 0.973000 | 0.824000 | 0.438000 | 0.762000 | 0.936000 | 0.830000 | 0.943000 | 0.873000 | 0.749000 | 1.452000 | 1.215000 | 0.729000 | 0.813000 | 0.683000 | 0.529000 | 0.389000 | 0.419000 | 0.404000 |
| 0.370000 | 0.449000 | 0.358000 | 0.400000 | 0.402000 | 0.625000 | 0.567000 | 1.175000 | 1.106000 | 0.448000 | 0.391000 | 0.723000 | 0.604000 | 0.754000 | 0.713000 | 0.830000 | 1.267000 | 1.237000 | 0.865000 | 0.790000 | 0.652000 | 0.649000 | 0.420000 | 0.489000 |
| 0.399000 | 0.372000 | 0.340000 | 0.371000 | 0.373000 | 0.591000 | 0.409000 | 0.744000 | 0.475000 | 0.649000 | 0.433000 | 0.536000 | 0.494000 | 0.561000 | 0.568000 | 0.583000 | 0.519000 | 0.543000 | 0.577000 | 0.483000 | 0.459000 | 0.440000 | 0.432000 | 0.432000 |
| 0.306000 | 0.394000 | 0.344000 | 0.352000 | 0.414000 | 0.617000 | 0.611000 | 0.861000 | 0.897000 | 0.971000 | 0.734000 | 0.552000 | 0.781000 | 0.465000 | 0.553000 | 0.621000 | 0.853000 | 0.776000 | 0.948000 | 0.507000 | 0.864000 | 0.348000 | 0.435000 | 0.331000 |
+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+
| 0.353000 | 0.393286 | 0.350714 | 0.374714 | 0.419143 | 0.639000 | 0.730286 | 0.777143 | 0.711000 | 0.718143 | 0.640143 | 0.694286 | 0.743143 | 0.625429 | 0.638714 | 0.757714 | 1.176429 | 0.946714 | 0.881714 | 0.782143 | 0.588143 | 0.513714 | 0.371143 | 0.395286 |
+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+----------+
```
