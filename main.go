package main

import (
	"fmt"
	"log"
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
	var leagues = []int{1, 2, 3, 4, 5, 9, 39, 140, 143, 556}
	res := utils.GetDailyFixtures(leagues)
	fmt.Printf("%v", res)
}
