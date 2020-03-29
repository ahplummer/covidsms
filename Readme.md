# Instructions
This leverages the [NYT Github repo](https://raw.githubusercontent.com/nytimes/covid-19-data/master/us-counties.csv) csv file that gets updated periodically. It will scan for the latest date for a given locale (designated by a FIPS number).  These FIPS numbers are county/state specific.   

## Prerequisites
1. You need a place for this to run. AWS Lightsail, EC2, docker containers are all good.
2. You need a Twilio account.

## Building
* Do a `go build -o covidsms`, be sure to set GOOS if you need to target something other than your build machine.
Example: `env GOOS=linux GOARCH=amd64 go build covidsms`

## Testing
* Set envvars as below, then do a `go test ./...` from root. NOTE: You'll need everything but the FIPS number, as the test will send via Twilio.

## Running
1. You need the following envvars set; and you can use a `.env` file like so:

```.env
export TWILIOAPI=<REDACTED>
export TWILIOSECRET=<REDACTED>
export NYTURL=https://raw.githubusercontent.com/nytimes/covid-19-data/master/us-counties.csv
export TWILIONUMBER=<REDACTED>
export FIPS=<Retrieve your FIPS location by scanning the CSV file above, and locating your county/state.>
export TONUMBERS=<REDACTED>
```

2. Once you create the above `.env` file, simply do a `source .env`, followed by a `./covidsms`.