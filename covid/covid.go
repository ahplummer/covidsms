package covid

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)
func RetrieveLatestDate(latestFName string) string {
	latestDate := "2001-01-01"
	if _, err := os.Stat(latestFName); os.IsNotExist(err) {
		fmt.Println(latestFName, "doesn't exist.")
	} else {
		file, err := os.Open(latestFName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fLine := scanner.Text()[:10]
			layout := "2006-01-02"
			fDate, err := time.Parse(layout, fLine)
			if err != nil{
				log.Panicf("Error with lastDate parm in file, not in YYYY-MM-DD format: %v", fDate)
			}
			latestDate = fLine
			break;
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
	return latestDate
}
func SendToTwilio(toNumbers []string, caseString, twilioApi,
	twilioSecret, twilioNumber string) error {

	for i := 0; i < len(toNumbers); i++ {
		urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + twilioApi + "/Messages.json"
		msgData := url.Values{}
		msgData.Set("To", toNumbers[i])
		msgData.Set("From", twilioNumber)
		msgData.Set("Body", caseString)
		msgDataReader := *strings.NewReader(msgData.Encode())

		client := &http.Client{}
		req, err := http.NewRequest("POST", urlStr, &msgDataReader)
		if err != nil{
			return err
		}
		req.SetBasicAuth(twilioApi, twilioSecret)
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		resp, err := client.Do(req)
		if err != nil{
			return err
		}
		if (resp.StatusCode >= 200 && resp.StatusCode < 300) {
			var data map[string]interface{}
			decoder := json.NewDecoder(resp.Body)
			err := decoder.Decode(&data)
			if (err == nil) {
				fmt.Println(data["sid"])
			} else {
				return err
			}
		} else {
			fmt.Println(resp.Status);
			return errors.New(resp.Status)
		}
	}
	return nil

}
func WriteStreamToFile(targetFileName string, stream io.ReadCloser) error {
	// Create the file
	out, err := os.Create(targetFileName)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, stream)
	return err
}

func RetrieveDataForLatest(url string,lastDate string, fips int, latestFName string) (string, error){
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	//err = WriteStreamToFile("workingfile.csv", resp.Body)
	lastCases, err := ParseCSVForFips(lastDate, fips, resp.Body, latestFName)
	return lastCases, err
}

// ParseCSVForLatestFips: will return rows from CSV that match
//2020-01-21,Snohomish,Washington,53061,1,0
//2020-01-22,Snohomish,Washington,53061,1,0
//2020-01-23,Snohomish,Washington,53061,1,0
func ParseCSVForFips(lastDate string, fips int, stream io.ReadCloser, latestFName string)(string, error){
	layout := "2006-01-02"
	latestDate, err := time.Parse(layout, lastDate)
	if err != nil{
		fmt.Println("Error with lastDate parm, not in YYYY-MM-DD format: ", lastDate)
		return "", err
	}

	reader := csv.NewReader(stream)
	records, _ := reader.ReadAll()
	var newRecord []string
	var newDate string
	var county string
	for i := 0; i < len(records); i++ {
		recDate, dateErr := time.Parse(layout, records[i][0])
		if dateErr == nil {
			if records[i][3] == fmt.Sprintf("%d", fips) &&
				recDate.After(latestDate) {
				//found a good record.
				newRecord = records[i]
				newDate = records[i][0]
				county = records[i][1]
				err = ioutil.WriteFile(latestFName, []byte(newDate), 0644)
			}
		}
	}
	lastCases := -1
	lastDeaths := -1
	var message string
	if newRecord != nil{
		lastCases, err = strconv.Atoi(newRecord[4])
		lastDeaths, err = strconv.Atoi(newRecord[5])
		if err !=  nil {
			return "", err
		} else {
			message = fmt.Sprintf("%d is the latest casecount, %d deaths for %v county, with date %v", lastCases, lastDeaths, county, newDate)
		}
	}
	return message, nil
}