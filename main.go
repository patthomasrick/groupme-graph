package main

import (
	"fmt"
	"log"

	"patrickwthomas.net/groupme-graph/groupme"
)

const groupIDTowel = "62858190"

func main() {
	// database.Init()
	// driver, err := database.NewNeo4j("bolt://localhost:7687", "", "", false)
	// if err != nil {
	// 	log.Panic(err)
	// }

	g, err := groupme.NewGroupMe()
	if err != nil {
		log.Panic(err)
	}
	groupIndex, err := g.GroupsIndex(1, 100, false)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Found %d groups.\n", len(groupIndex))

	// for _, group := range groupIndex {
	// 	group.SaveToNeo4j(driver)
	// }

	// groupTowel := g.GroupsShow(groupIDTowel)

	// groupTowelMessages := g.MessagesIndex(groupTowel.ID, "", "", "", 50)
	// for i := 0; i < len(groupTowelMessages); i++ {
	// 	m := groupTowelMessages[i]
	// 	m.SaveToNeo4j(driver)
	// }

	// groupme.Connect(driver)
}
