package utils

import (
	"fmt"
	"io"
	"net/http"
     "os"
	"github.com/tidwall/gjson"
)

type Result struct {
	score string
	link  string
}
type Fixture struct {
	leauge     string
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
	rapidApiKey, exists := os.LookupEnv("RAPID_API_KEY")
    if exists{
		println(rapidApiKey)
	}
	req.Header.Add("x-rapidapi-key",rapidApiKey)
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

func GetDailyFixtures(leagues []int) []Fixture {
	date := "2023-05-07"
	season := 2022
	fixtureList := []Fixture{}
	for i, league := range leagues {
		if i==1{
			break
		}
		league = 39
		url := fmt.Sprintf("https://api-football-v1.p.rapidapi.com/v3/fixtures?date=%s&league=%d&season=%d", date, league, season)
		// fmt.Println(url)
		status, res := callApi(url)
		if status == 200 {
			result := gjson.GetBytes(res, "response")
			for _, r := range result.Array() {
				
				println(r.Get("league.name").String())
				println()
			}
			
			fmt.Println(status)
		
		}
	}
	return fixtureList
}
