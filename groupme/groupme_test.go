package groupme

import (
	"log"
	"os"
	"testing"

	"patrickwthomas.net/groupme-graph/local"
)

func setupTestDir() {
	log.Printf("Making temp directory")
	err := os.Mkdir("/tmp/test-groupme-graph", os.ModePerm)
	if os.IsExist(err) {
		os.Remove("/tmp/test-groupme-graph/settings.json")
	} else if err != nil {
		log.Panic(err)
	}

	log.Printf("CDing into temp directory")
	err = os.Chdir("/tmp/test-groupme-graph")
	if err != nil {
		log.Panic(err)
	}
}

func TestNewGroupMe(t *testing.T) {
	setupTestDir()
	g, err := NewGroupMe()
	if err == nil || g != nil {
		t.Fail()
	}

	s, _ := local.NewSettings("api", "deadbeef")
	local.SaveSettings(s)

	g, err = NewGroupMe()
	if err != nil {
		t.Fail()
	}

	if g.settings.AccessToken != s.AccessToken {
		t.Fail()
	}
	if g.settings.GroupMeAPI != s.GroupMeAPI {
		t.Fail()
	}
}
