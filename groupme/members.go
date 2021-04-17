package groupme

import (
	"fmt"
	"log"

	"patrickwthomas.net/groupme-graph/database"
)

// AddMember is returned after adding members.
type AddMember struct {
	Nickname    string `json:"nickname"`
	UserID      string `json:"user_id,omitempty"`
	GUID        string `json:"guid"`
	PhoneNumber string `json:"phone_number,omitempty"`
	Email       string `json:"email,omitempty"`
}

// Member contains information about a GroupMe user.
type Member struct {
	ID           string `json:"id"`
	UserID       string `json:"user_id"`
	Nickname     string `json:"nickname"`
	Muted        bool   `json:"muted"`
	ImageURL     string `json:"image_url"`
	Autokicked   bool   `json:"autokicked"`
	AppInstalled bool   `json:"app_installed"`
	GUID         string `json:"guid"`
}

type membersAddResult struct {
	ResultsID string `json:"results_id"`
}

type changeNickname struct {
	Membership struct {
		Nickname string `json:"nickname"`
	} `json:"membership"`
}

// MembersAdd adds a member to a group.
func (g *GroupMe) MembersAdd(groupID string, m []AddMember) string {
	result := &membersAddResult{}
	_, err := g.groupMeRequestPostObject(fmt.Sprintf("/groups/%s/members/add", groupID), m, result)
	if err != nil {
		log.Panic(err)
	}
	return result.ResultsID
}

// MembersResults gets the results of a MembersAdd operation.
func (g *GroupMe) MembersResults(groupID string, addResultGUID string) []Member {
	result := &[]Member{}
	_, err := g.groupMeRequest("GET", fmt.Sprintf("/groups/%s/members/results/%s", groupID, addResultGUID), nil, result)
	if err != nil {
		log.Panic(err)
	}
	return *result
}

// MembersRemove removes a single member from a group.
func (g *GroupMe) MembersRemove(groupID string, membershipID string) {
	_, err := g.groupMeRequest("POST", fmt.Sprintf("/groups/%s/members/%s/remove", groupID, membershipID), nil, nil)
	if err != nil {
		log.Panic(err)
	}
}

// MembersUpdate updates your nickname in a group.
func (g *GroupMe) MembersUpdate(groupID, nickname string) Member {
	values := changeNickname{}
	values.Membership.Nickname = nickname
	result := &Member{}
	_, err := g.groupMeRequestPostObject(fmt.Sprintf("/groups/%s/memberships/update", groupID), values, result)
	if err != nil {
		log.Panic(err)
	}
	return *result
}

// SaveToNeo4j saves the current member into the database.
func (m *Member) SaveToNeo4j(driver *database.Neo4j, group_id string) {
	session, err := driver.NewWriteSession()
	if err != nil {
		log.Panic(err)
	}
	defer session.Close()

	query := `
	MERGE (m:Member{UserID: $UserID})
	ON CREATE
		SET
			m.ID = $ID,
			m.UserID = $UserID,
			m.Nickname = $Nickname,
			m.Muted = $Muted,
			m.ImageURL = $ImageURL,
			m.Autokicked = $Autokicked,
			m.AppInstalled = $AppInstalled,
			m.GUID = $GUID
	ON MATCH
		SET
			m.Nickname = $Nickname,
			m.Muted = $Muted,
			m.ImageURL = $ImageURL,
			m.Autokicked = $Autokicked,
			m.AppInstalled = $AppInstalled
	WITH m
	MATCH (g:Group{ID:$GroupID})
	MERGE (g)-[:MEMBER]->(m)
	`
	result, err := session.Run(query, map[string]interface{}{
		"ID":           m.ID,
		"UserID":       m.UserID,
		"Nickname":     m.Nickname,
		"Muted":        m.Muted,
		"ImageURL":     m.ImageURL,
		"Autokicked":   m.Autokicked,
		"AppInstalled": m.AppInstalled,
		"GUID":         m.GUID,
		"GroupID":      group_id,
	})
	e := result.Err()
	if err != nil {
		log.Panic(err)
	} else if e != nil {
		log.Panic(e)
	}
}
