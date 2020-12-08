package local

import (
	"log"
	"os"
	"testing"
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

func TestLoadSettingsNoInit(t *testing.T) {
	setupTestDir()

	_, err := LoadSettings()
	if err != nil && err.Error() != "new settings file created at ./settings.json, the access token needs to be configured" {
		t.Fail()
	}
}

func TestLoadSettingsInitNotSet(t *testing.T) {
	setupTestDir()

	err := initSettings()
	if err != nil {
		log.Panic(err)
	}
	_, err = LoadSettings()
	if err != nil && err.Error() != "access token needs to be configured at ./settings.json" {
		t.Fail()
	}
}

func TestLoadSettingsEmpty(t *testing.T) {
	setupTestDir()

	err := SaveSettings(&Settings{
		GroupMeAPI:  "",
		AccessToken: "",
	})
	if err != nil {
		log.Panic(err)
	}

	_, err = LoadSettings()
	if err != nil && err.Error() != "settings file is empty" {
		t.Fail()
	}
}

func TestLoadSettingsSuccess(t *testing.T) {
	setupTestDir()

	err := SaveSettings(&Settings{
		GroupMeAPI:  "api",
		AccessToken: "token",
	})
	if err != nil {
		log.Panic(err)
	}

	s, err := LoadSettings()
	if err != nil {
		t.Fail()
	} else if !(s != nil && s.GroupMeAPI == "api" && s.AccessToken == "token") {
		t.Fail()
	}
}

func TestInitSettingsTryToCreate(t *testing.T) {
	setupTestDir()

	// Begin the actual test.
	err := initSettings()
	if err != nil {
		log.Panic(err)
	}

	_, err = os.Stat("/tmp/test-groupme-graph/settings.json")
	if err != nil {
		t.Fail()
	}
}

func TestInitSettingsFileExists(t *testing.T) {
	setupTestDir()

	// Create the file in advance...
	file, _ := os.Create(settingsFileDir)
	file.Close()

	// Begin the actual test.
	err := initSettings()
	if err != nil {
		t.Fail()
	}
}

func TestSaveSettings(t *testing.T) {
	s := Settings{
		GroupMeAPI:  groupMeAPIUninit,
		AccessToken: "my access token",
	}

	err := SaveSettings(&s)
	if err != nil {
		log.Panic(err)
	}
}
