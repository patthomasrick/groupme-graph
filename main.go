package main

import (
	"fmt"
	"log"

	"patrickwthomas.net/groupme-graph/database"
	"patrickwthomas.net/groupme-graph/groupme"
)

const targetGroupID = "64420461"

func main() {
	database.Init()
	driver, err := database.NewNeo4j("bolt://localhost:7687", "", "", false)
	if err != nil {
		log.Panic(err)
	}

	g, err := groupme.NewGroupMe()
	if err != nil {
		log.Panic(err)
	}
	groupIndex, err := g.GroupsIndex(1, 100, false)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Found %d groups.\n", len(groupIndex))
	for _, group := range groupIndex {
		fmt.Printf("\t%s %s\n", group.Name, group.ID)
	}

	for _, group := range groupIndex {
		group.SaveToNeo4j(driver)
	}

	targetGroup, _ := g.GroupsShow(targetGroupID)

	targetGroupMessages := g.MessagesIndex(targetGroup.ID, "", "", "", 50)
	for i := 0; i < len(targetGroupMessages); i++ {
		m := targetGroupMessages[i]
		m.SaveToNeo4j(driver)
	}

	// groupme.Connect(driver)
}
