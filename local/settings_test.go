package local

import (
	"log"
	"os"
	"reflect"
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
	setupTestDir()

	s := Settings{
		GroupMeAPI:  groupMeAPIUninit,
		AccessToken: "my access token",
	}

	err := SaveSettings(&s)
	if err != nil {
		log.Panic(err)
	}
}

func TestNewSettings(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *Settings
		wantErr bool
	}{
		{
			name:  "Valid",
			input: "a",
			want: &Settings{
				"a",
				"a",
			},
			wantErr: false,
		},
		{
			name:    "Invalid",
			input:   "{}!!*@#*(#*(#*(@)))}",
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewSettings("a", tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSettings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSettings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSettings_GetGroupMeAPI_SetGroupMeAPI(t *testing.T) {
	s := new(Settings)
	s.SetGroupMeAPI("a")
	if s.GetGroupMeAPI() != "a" {
		t.Fail()
	}
	s.SetGroupMeAPI("b")
	if s.GetGroupMeAPI() != "b" {
		t.Fail()
	}
}

func TestSettings_GetAccessToken_SetAccessToken(t *testing.T) {
	s, _ := NewSettings("a", "deadbeef")
	if s.GetAccessToken() != "deadbeef" {
		t.Fail()
	}

	err := s.SetAccessToken("1010101")
	if err != nil {
		t.Fail()
	}

	if s.GetAccessToken() != "1010101" {
		t.Fail()
	}

	err = s.SetAccessToken("{}{}{}@#$@#$(@$)%")
	if err == nil {
		t.Fail()
	}
}
