package covid

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)
var nytURL = "https://raw.githubusercontent.com/nytimes/covid-19-data/master/us-counties.csv"
func TestRetrieveDataForLatest(t *testing.T) {
	expectedNewCases := "3 is the latest casecount for Snohomish county, with date 2020-01-23"
	lastRetrievedDate := "2020-01-22"
	latestFName := "test.txt"
	got, err := RetrieveDataForLatest(nytURL, lastRetrievedDate, 53061,latestFName)
	if err != nil{
		t.Errorf("Threw error: %v", err)
	}
	if got < expectedNewCases {
		t.Error("Retrieved", got, "Expected", expectedNewCases)
	}
}

func TestSendToTwilio(t *testing.T) {
	caseString := "Test string for Twilio"
	twilioApi := os.Getenv("TWILIOAPI")
	twilioSecret := os.Getenv("TWILIOSECRET")
	twilioNumber := os.Getenv("TWILIONUMBER")
	toNumbers := strings.Split(os.Getenv("TONUMBERS"),",")
	err := SendToTwilio(toNumbers, caseString, twilioApi, twilioSecret,twilioNumber)
	if err != nil{
		t.Error("Error thrown: ", err)
	}
}
func TestParseCSVForFips(t *testing.T) {
	testData := `date,county,state,fips,cases,deaths
2020-01-21,Snohomish,Washington,53061,1,0
2020-01-22,Snohomish,Washington,53061,1,0
2020-01-23,Snohomish,Washington,53061,3,0`
	expectedNewCases := "3 is the latest casecount, 0 deaths for Snohomish county, with date 2020-01-23"
	lastRetrievedDate := "2020-01-22"
	fakeBody := ioutil.NopCloser(bytes.NewReader([]byte(testData)))
	latestFName := "test.txt"
	got, err := ParseCSVForFips(lastRetrievedDate, 53061, fakeBody, latestFName)
	if err != nil{
		t.Errorf("Threw error: %v", err)
	}
	if got != expectedNewCases {
		t.Error("Retrieved", got, "Expected", expectedNewCases)
	}
	expectedNewCases = ""
	lastRetrievedDate = "2020-01-23"
	fakeBody = ioutil.NopCloser(bytes.NewReader([]byte(testData)))
	got, err = ParseCSVForFips(lastRetrievedDate, 53061, fakeBody, latestFName)
	if err != nil{
		t.Errorf("Threw error: %v", err)
	}
	if got != expectedNewCases {
		t.Error("Retrieved", got, "Expected", expectedNewCases)
	}

}