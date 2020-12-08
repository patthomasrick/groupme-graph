package main

import (
	"fmt"
	"os"

	"patrickwthomas.net/groupme-graph/database"
	"patrickwthomas.net/groupme-graph/groupme"
)

const groupIDTowel = "62858190"

func main() {
	database.Init()
	driver, err := database.NewNeo4j("bolt://localhost:7687", "", "", false)
	if err != nil {
		panic(err)
	}

	g := groupme.NewGroupMe(os.Getenv("GROUPME_ACCESS_TOKEN"))
	groupIndex := g.GroupsIndex(1, 100, false)

	fmt.Printf("Found %d groups.\n", len(groupIndex))

	for _, group := range groupIndex {
		group.SaveToNeo4j(driver)
	}

	groupTowel := g.GroupsShow(groupIDTowel)

	groupTowelMessages := g.MessagesIndex(groupTowel.ID, "", "", "", 50)
	for i := 0; i < len(groupTowelMessages); i++ {
		m := groupTowelMessages[i]
		m.SaveToNeo4j(driver)
	}

	groupme.Connect(driver)
}
