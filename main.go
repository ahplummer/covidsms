package main

import (
	"covid/covid"
	"fmt"
	"os"
	"strconv"
	"strings"
)
func displayHelp() string{
	return `You need to have envvars set (TWILIOAPI, TWILIOSECRET,
NYTURL, TWILIONUMBER, TONUMBERS, FIPS`
}
func main() {
	twilioApi := os.Getenv("TWILIOAPI")
	twilioSecret := os.Getenv("TWILIOSECRET")
	nytURL := os.Getenv("NYTURL")
	twilioNumber := os.Getenv("TWILIONUMBER")
	toNumbers := strings.Split(os.Getenv("TONUMBERS"),",")
	fips := os.Getenv("FIPS")
	if twilioApi == "" || twilioSecret == "" ||
		nytURL == "" ||  twilioNumber == "" ||
		len(toNumbers) == 0 || fips == "" {
		panic(displayHelp())
	}
	fipsNumber, err := strconv.Atoi(fips)
	if err != nil{
		panic(err)
	}
	latestFName := "LatestDateChecked.txt"
	latestDate := covid.RetrieveLatestDate(latestFName)
	fmt.Println("Will get the latest date after: ", latestDate, "from FIPS: ", fipsNumber)
	caseString, err := covid.RetrieveDataForLatest(nytURL, latestDate, fipsNumber, latestFName)
	if caseString != "" {
		fmt.Println("Packaging up SMS: ", caseString)
		covid.SendToTwilio(toNumbers, caseString, twilioApi, twilioSecret,twilioNumber)

	} else {
		fmt.Println("Nothing new, bailing")
	}
}