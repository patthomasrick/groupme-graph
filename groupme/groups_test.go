package groupme

import (
	"log"
	"os"
	"testing"

	"patrickwthomas.net/groupme-graph/local"
)

var testAccessToken = os.Getenv("GROUPME_TEST_ACCESS_TOKEN")

func setupGroupMe() GroupMe {
	s, err := local.NewSettings("https://api.groupme.com/v3", testAccessToken)
	if err != nil {
		log.Panic(err)
	}

	g := GroupMe{
		settings: *s,
	}
	return g
}

func TestGroupsIndex(t *testing.T) {
	g := setupGroupMe()

	index, err := g.GroupsIndex(1, 10, false)
	if err != nil {
		log.Panic(err)
	}

	for i, group := range index {
		log.Printf("%d %s\n", i, group.Name)
	}

	if len(index) != 1 {
		t.Fail()
	}
	if index[0].Name != "Testing Group" {
		t.Fail()
	}
	if len(index[0].Members) != 1 {
		t.Fail()
	}
}

func TestGroupsIndexOmit(t *testing.T) {
	g := setupGroupMe()

	index, err := g.GroupsIndex(1, 10, true)
	if err != nil {
		log.Panic(err)
	}
	if len(index[0].Members) != 0 {
		t.Fail()
	}
}

func TestGroupsFormer(t *testing.T) {
	g := setupGroupMe()

	index, err := g.GroupsFormer()
	if err != nil {
		log.Panic(err)
	}
	if len(index) != 0 {
		t.Fail()
	}
}

func TestGroupsShow(t *testing.T) {
	g := setupGroupMe()

	group, err := g.GroupsShow("64420461")
	if err != nil {
		log.Panic(err)
	}

	if group == nil || group.Name != "Testing Group" {
		t.Fail()
	}
}

func TestGroupsShowBad(t *testing.T) {
	g := setupGroupMe()

	_, err := g.GroupsShow("0")
	if err == nil {
		t.Fail()
	}
}

func TestGroupsCreate(t *testing.T) {
	g := setupGroupMe()

	gName := "My New Test Group"
	gDesc := "My Group's Description"
	group, err := g.GroupsCreate(gName, gDesc, "", false)
	if err != nil {
		log.Panic(err)
	}
	if group.Name != gName || group.Description != gDesc {
		t.Fail()
	}

	// Find the group now.
	index, _ := g.GroupsIndex(1, 10, false)
	groupID := ""
	for _, group := range index {
		if group.Name == gName {
			groupID = group.ID
			break
		}
	}

	group, _ = g.GroupsShow(groupID)
	if group.Name != gName || group.Description != gDesc {
		t.Fail()
	}
}
