package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Sabretooth1405/footballbot/utils"
	"github.com/joho/godotenv"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}
func main() {
	start:=time.Now()
	var leagues = []int{39,1, 2, 3, 4, 5, 9,  140, 143, 556}
	res := utils.GetDailyFixtures(leagues)
	fmt.Printf("%v", res)
	end:=time.Now()
	fmt.Printf("%d",end.Unix()-start.Unix())

}
