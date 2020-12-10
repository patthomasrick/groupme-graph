package local

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

// Settings holds all of the relevant settings from the application.
type Settings struct {
	GroupMeAPI  string `json:"group_me_api"`
	AccessToken string `json:"access_token"`
}

const settingsFileDir = "./settings.json"
const groupMeAPIUninit = "https://api.groupme.com/v3"
const accessTokenUninit = "your API token here (get one here: https://dev.groupme.com/)"

var accessTokenPattern = regexp.MustCompile(`[a-f0-9]+`)

// GetGroupMeAPI get the set API.
func (s *Settings) GetGroupMeAPI() string {
	return s.GroupMeAPI
}

// SetGroupMeAPI set the API URL.
func (s *Settings) SetGroupMeAPI(v string) {
	s.GroupMeAPI = v
}

// GetAccessToken get the curent access token.
func (s *Settings) GetAccessToken() string {
	return s.AccessToken
}

// SetAccessToken set the current access token. Can error if not valid.
func (s *Settings) SetAccessToken(v string) error {
	tokenMatch := accessTokenPattern.FindString(v)
	if tokenMatch == "" {
		return fmt.Errorf("access token is invalid (is not a hexadecimal number, all lowercase)")
	}
	s.AccessToken = tokenMatch
	return nil
}

// NewSettings allows for the direct creation of a settings object.
func NewSettings(groupMeAPI, accessToken string) (*Settings, error) {
	s := new(Settings)
	s.GroupMeAPI = groupMeAPI
	s.AccessToken = accessTokenPattern.FindString(accessToken)
	if len(s.AccessToken) == 0 {
		return nil, fmt.Errorf("access token is invalid (is not a hexadecimal number, all lowercase)")
	}
	return s, nil
}

// LoadSettings loads configuration for the application from the harddisk.
func LoadSettings() (*Settings, error) {
	fileContents, err := ioutil.ReadFile(settingsFileDir)
	if os.IsNotExist(err) {
		err := initSettings()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("new settings file created at %s, the access token needs to be configured", settingsFileDir)
	} else if err != nil {
		return nil, err
	}

	s := new(Settings)
	err = json.Unmarshal(fileContents, s)

	if s.AccessToken == accessTokenUninit {
		return nil, fmt.Errorf("access token needs to be configured at %s", settingsFileDir)
	} else if s.AccessToken == "" && s.GroupMeAPI == "" {
		return nil, fmt.Errorf("settings file is empty")
	}

	// Validate access token.

	return s, nil
}

// SaveSettings Save settings to a file.
func SaveSettings(s *Settings) error {
	jsonMarshalled, _ := json.MarshalIndent(s, "", "\t")

	file, err := os.Create(settingsFileDir)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(jsonMarshalled)
	if err != nil {
		return err
	}
	return nil
}

func initSettings() error {
	s := Settings{
		GroupMeAPI:  groupMeAPIUninit,
		AccessToken: accessTokenUninit,
	}

	err := SaveSettings(&s)
	if err != nil {
		return err
	}
	return nil
}
