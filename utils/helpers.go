/*
Helper functions for the service.Functions ar
*/
package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
	"github.com/tidwall/gjson"
)

type Result struct {
	score string
	link  string
}
type Fixture struct {
	league     string
	homeTeam   string
	awayTeam   string
	startTime  string
	result     Result
	streamLink string
}

func callApi(url string) (int, []byte) {

	method := "GET"
	var empty = []byte{}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {

		fmt.Println(err)

		return -1, empty
	}
	rapidApiKey, _ := os.LookupEnv("RAPID_API_KEY")
	// if exists {
	// 	println(rapidApiKey)
	// }
	req.Header.Add("x-rapidapi-key", rapidApiKey)
	req.Header.Add("x-rapidapi-host", "api-football-v1.p.rapidapi.com")

	res, err := client.Do(req)
	if err != nil {

		fmt.Println(err)
		fmt.Println('a')
		return res.StatusCode, empty
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return -1, empty
	}
	return res.StatusCode, body
}
func getDateTimeFromTimeStamp(timestamp string) string {
	i, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		panic(err)
	}
	datetime := time.Unix(i, 0)
	startTime := fmt.Sprintf("%s %d %s", datetime.Month().String(), datetime.Day(), datetime.Format(time.Kitchen))
	return startTime
}
const (
	YYYYMMDD="2006-01-02"
)
func GetDailyFixtures(leagues []int) []Fixture {
	currentDate:=time.Now().Format(YYYYMMDD)
	fmt.Println(currentDate)
	date := currentDate
	season := 2022
	fixtureList := []Fixture{}
	for _, league := range leagues {
		url := fmt.Sprintf("https://api-football-v1.p.rapidapi.com/v3/fixtures?date=%s&league=%d&season=%d&timezone=Asia/Calcutta", date, league, season)
		// fmt.Println(url)
		status, res := callApi(url)
		if status == 200 {
			result := gjson.GetBytes(res, "response")
			resultCount := gjson.GetBytes(res, "results").Int()
			fmt.Println(resultCount)
			if resultCount == 0 {
				continue
			}
			for _, r := range result.Array() {
				println(r.Get("fixture").String())
				timestamp := r.Get("fixture.timestamp").String()
				startTime := getDateTimeFromTimeStamp(timestamp)
				var result = Result{score: "ns", link: "ns"}
				var newFixture = Fixture{league: r.Get("league.name").String(),
					homeTeam:   r.Get("teams.home.name").String(),
					awayTeam:   r.Get("teams.away.name").String(),
					startTime:  startTime,
					result:     result,
					streamLink: "ns"}
				fixtureList = append(fixtureList, newFixture)
			}
		}else if status==-1{
			fmt.Println("Invalid Response or Some madness in urls (latter is more likely)")
		} else{
			fmt.Println("Bad request or API is down")
		}
		
	}
	return fixtureList
}
