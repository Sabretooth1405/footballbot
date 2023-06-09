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
	"sync"
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
	YYYYMMDD = "2006-01-02"
)

// func leagueSliceToChannel(leagues []int) <-chan int {
// 	out := make(chan int, cap(leagues))
// 	go func() {
// 		for _, league := range leagues {
// 			out <- league
// 		}
// 	}()

//		close(out)
//		return out
//	}
var wg sync.WaitGroup
var mutex = sync.Mutex{}

func GetDailyFixtures(leagues []int) []Fixture {
	currentDate := time.Now().Format(YYYYMMDD)
	fmt.Println(currentDate)
	date := currentDate
	season := 2022
	fixtureList := []Fixture{}
	for _, league := range leagues {
		leaugeCopy := league
		wg.Add(1)
		go func() {
			defer wg.Done()
			url := fmt.Sprintf("https://api-football-v1.p.rapidapi.com/v3/fixtures?date=%s&league=%d&season=%d&timezone=Asia/Calcutta", date, leaugeCopy, season)
			// fmt.Println(url)
			status, res := callApi(url)
			if status == 200 {
				result := gjson.GetBytes(res, "response")
				resultCount := gjson.GetBytes(res, "results").Int()
				fmt.Println(resultCount)
				if resultCount == 0 {
					return
				}
				for _, r := range result.Array() {
					println(r.Get("fixture").String())
					timestamp := r.Get("fixture.timestamp").String()
					startTime := getDateTimeFromTimeStamp(timestamp)
					var result = Result{}
					fmt.Println(r.Get("fixture.status.short").String())
					if r.Get("fixture.status.short").String() == "FT" {
						result.score = fmt.Sprintf("%s %s-%s %s", r.Get("teams.home.name").String(),
							r.Get("score.fulltime.home").String(), r.Get("score.fulltime.away").String(),
							r.Get("teams.away.name").String())
						result.link = "NA"
					} else {
						result.score = "ns or pt"
						result.link = "NA"
					}
					var newFixture = Fixture{league: r.Get("league.name").String(),
						homeTeam:   r.Get("teams.home.name").String(),
						awayTeam:   r.Get("teams.away.name").String(),
						startTime:  startTime,
						result:     result,
						streamLink: "NA"}
					mutex.Lock()
					fixtureList = append(fixtureList, newFixture)
					mutex.Unlock()
				}
			} else if status == -1 {
				fmt.Println("Invalid Response or Some madness in urls (latter is more likely)")
				return
			} else {
				fmt.Println("Bad request or API is down")
				return
			}

		}()

	}
	wg.Wait()
	return fixtureList
}
