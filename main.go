package main

import (
	"fmt"
	"log"

	wikimebackend "github.com/WikimeCorp/WikimeBackend/cmd/WikimeBackend"
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
	//id := AnimeID(0)
	// a, _ := primitive.ObjectIDFromHex("6334220ab1f26f806521a0e4")
	// err := db.DeleteCommentByID(a)
	// log.Println(a, " ", err)
	// testGetRating(id)

	// err := db.ChangeRating(0, 5, 1)
	// if err != nil {
	// 	log.Fatal("ERROR ", err)
	// }

	// testGetRating(id)
	fmt.Println("Start server")
	log.Fatal(wikimebackend.Start())
}
