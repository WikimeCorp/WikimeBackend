package main

import (
	"log"

	"github.com/WikimeCorp/WikimeBackend/db"
	. "github.com/WikimeCorp/WikimeBackend/types"
)

func testGetRating(id AnimeID) {
	ans, err := db.GetRating(id)
	if err != nil {
		log.Fatal("ERROR ", err)
	}
	log.Printf("%+v\n", ans)
}

func main() {
	id := AnimeID(0)
	err := db.Rate(id, 0, 3)
	log.Println(err)
	// testGetRating(id)

	// err := db.ChangeRating(0, 5, 1)
	// if err != nil {
	// 	log.Fatal("ERROR ", err)
	// }

	// testGetRating(id)
}
