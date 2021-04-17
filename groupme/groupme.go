package groupme

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"

	"patrickwthomas.net/groupme-graph/database"
	"patrickwthomas.net/groupme-graph/local"
)

// GroupMe holds data pertaining to the current GroupMe instance.
type GroupMe struct {
	settings local.Settings
}

type meta struct {
	Code   int      `json:"code"`
	Errors []string `json:"errors"`
}

type envelope struct {
	Meta     meta        `json:"meta"`
	Response interface{} `json:"response"`
}

var quotePattern = regexp.MustCompile(`[{,]"[\w_]+":`)

// NewGroupMe creates a new GroupMe manager.
func NewGroupMe() (*GroupMe, error) {
	g := new(GroupMe)
	s, err := local.LoadSettings()
	if err != nil {
		return nil, err
	}
	g.settings = *s
	return g, nil
}

func letterOpener(responseBody []byte, dest interface{}) (meta, error) {
	meta := meta{}
	unmarshalledResponse := envelope{Meta: meta, Response: dest}
	err := json.Unmarshal(responseBody, &unmarshalledResponse)
	if err != nil {
		// log.Println(string(responseBody))
		return meta, err
	} else if unmarshalledResponse.Meta.Code/100 != 2 {
		return meta, fmt.Errorf("%d: %v", unmarshalledResponse.Meta.Code, unmarshalledResponse.Meta.Errors)
	}
	return meta, err
}

func (g *GroupMe) groupMeRequest(method, requestSubDir string, values map[string]string, dest interface{}) (meta, error) {
	request, err := http.NewRequest(method, g.settings.GroupMeAPI+requestSubDir, nil)
	if err != nil {
		return meta{}, err
	}

	// Build the queries and get a response. Could be either GET or POST
	var response *http.Response
	query := request.URL.Query()
	query.Add("token", g.settings.AccessToken)
	if method == "GET" {
		if values != nil {
			for k, v := range values {
				query.Add(k, v)
			}
		}
		request.URL.RawQuery = query.Encode()
		response, err = http.Get(request.URL.String())
	} else if method == "POST" {
		postData := url.Values{}
		if values != nil {
			for k, v := range values {
				postData.Add(k, v)
			}
		}
		request.URL.RawQuery = query.Encode()
		response, err = http.PostForm(request.URL.String(), postData)
	}
	if err != nil {
		return meta{}, err
	}
	defer response.Body.Close()

	// Get the message body out of the response.
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return meta{}, err
	}

	// Extract the response from the GroupMe envelope.
	return letterOpener(body, dest)
}

func (g *GroupMe) groupMeRequestPostObject(requestSubDir string, values interface{}, dest interface{}) (meta, error) {
	request, err := http.NewRequest("POST", g.settings.GroupMeAPI+requestSubDir, nil)
	if err != nil {
		return meta{}, err
	}

	// Marshall the input object.
	marshalled, err := json.Marshal(values)
	if err != nil {
		return meta{}, err
	}

	// Build the queries and get a response. Could be either GET or POST
	var response *http.Response
	query := request.URL.Query()
	query.Add("token", g.settings.AccessToken)
	request.URL.RawQuery = query.Encode()
	response, err = http.Post(request.URL.String(), "application/json", bytes.NewBuffer(marshalled))
	if err != nil {
		return meta{}, err
	}
	defer response.Body.Close()

	// Get the message body out of the response.
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return meta{}, err
	}

	// Extract the response from the GroupMe envelope.
	return letterOpener(body, dest)
}

func noQuotes(s string) string {
	return strings.ReplaceAll(s, "\"", "")
}

func quoteEscape(s string) string {
	return strings.ReplaceAll(s, "\"", "\\\"")
}

// Melt returns a modified JSON format with no quotes around keys.
func Melt(v interface{}) string {
	cypherParts := []string{}

	value := reflect.ValueOf(v)
	typeOfS := value.Type()
	for i := 0; i < value.NumField(); i++ {
		if typeOfS.Field(i).Type.Name() == "string" {
			s := strings.ReplaceAll(fmt.Sprintf(`%v`, value.Field(i).Interface()), "\"", "\\\"")
			cypherParts = append(cypherParts, fmt.Sprintf(`%s:"%v"`, typeOfS.Field(i).Name, s))
		} else if typeOfS.Field(i).Type.Name() == "int" {
			cypherParts = append(cypherParts, fmt.Sprintf(`%s:%v`, typeOfS.Field(i).Name, value.Field(i).Interface()))
		} else if typeOfS.Field(i).Type.Name() == "bool" {
			cypherParts = append(cypherParts, fmt.Sprintf(`%s:%v`, typeOfS.Field(i).Name, value.Field(i).Interface()))
		}
	}

	return strings.Join(cypherParts, ", ")
}

// ConnectData connects the data in the graph database as best it can.
func ConnectData(driver *database.Neo4j) error {
	session, err := driver.NewWriteSession()
	if err != nil {
		return err
	}
	defer session.Close()

	result, err := session.Run(`MATCH (m:Member), (n:Message) WHERE n.UserID = m.UserID
	MERGE (m)-[:AUTHORED]->(n)`, map[string]interface{}{})
	if err != nil {
		return err
	} else if result.Err() != nil {
		return result.Err()
	}

	result, err = session.Run(`MATCH (m:Group), (n:Message) WHERE n.GroupID = m.ID
	MERGE (m)-[:HAS_MESSAGE]->(n)`, map[string]interface{}{})
	if err != nil {
		return err
	} else if result.Err() != nil {
		return result.Err()
	}

	return nil
}
