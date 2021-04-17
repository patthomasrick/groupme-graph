package groupme

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"

	"patrickwthomas.net/groupme-graph/database"
)

// MessagesIndex is the index format returned by GroupMe.
type MessagesIndex struct {
	Count    int       `json:"count"`
	Messages []Message `json:"messages"`
}

// Message contains all information about a message.
type Message struct {
	ID          string       `json:"id"`
	SourceGUID  string       `json:"source_guid"`
	CreatedAt   int          `json:"created_at"`
	UserID      string       `json:"user_id"`
	GroupID     string       `json:"group_id"`
	Name        string       `json:"name"`
	AvatarURL   string       `json:"avatar_url"`
	Text        string       `json:"text"`
	System      bool         `json:"system"`
	FavoritedBy []string     `json:"favorited_by"`
	Attachments []Attachment `json:"attachments"`
}

// Attachment is the base struct for all attachments. Does nothing on its own.
type Attachment struct {
	Type        string  `json:"type"`
	URL         string  `json:"url"`
	Lat         string  `json:"lat"`
	Lng         string  `json:"lng"`
	Name        string  `json:"name"`
	Token       string  `json:"token"`
	Placeholder string  `json:"placeholder"`
	Charmap     [][]int `json:"charmap"`
}

// MessagesIndex gets the groups index from GroupMe.
func (g *GroupMe) MessagesIndex(groupID, beforeID, sinceID, afterID string, limit int) []Message {
	messages := &MessagesIndex{}
	urlValues := map[string]string{
		"limit": fmt.Sprint(limit),
	}

	if beforeID != "" {
		urlValues["before_id"] = beforeID
	} else if sinceID != "" {
		urlValues["since_id"] = sinceID
	} else if afterID != "" {
		urlValues["after_id"] = afterID
	}

	_, err := g.groupMeRequest("GET", fmt.Sprintf("/groups/%s/messages", groupID), urlValues, messages)
	if err != nil && errors.Is(err, err.(*json.SyntaxError)) {
		// Pass
	} else if err != nil {
		log.Panic(err)
	}

	sort.Slice(messages.Messages, func(i, j int) bool {
		return messages.Messages[i].CreatedAt < messages.Messages[j].CreatedAt
	})

	return messages.Messages
}

// SaveToNeo4j saves the current message into the database.
func (m *Message) SaveToNeo4j(driver *database.Neo4j) {
	session, err := driver.NewWriteSession()
	if err != nil {
		log.Panic(err)
	}
	defer session.Close()

	query := `
	MERGE (m:Message{ID: $ID})
	ON CREATE
		SET
			m.ID = $ID,
			m.SourceGUID = $SourceGUID,
			m.CreatedAt = $CreatedAt,
			m.UserID = $UserID,
			m.GroupID = $GroupID,
			m.Name = $Name,
			m.AvatarURL = $AvatarURL,
			m.Text = $Text,
			m.System = $System,
			m.FavoritedBy = $FavoritedBy
	ON MATCH
		SET
			m.FavoritedBy = $FavoritedBy
	WITH m
	MATCH (g:Group{ID: $GroupID}), (u:Member{UserID: $UserID})
	MERGE (g)-[:GROUP_OF]->(m)<-[:AUTHORED]-(u)
	WITH m
	UNWIND m.FavoritedBy as favoritedUserID
	MATCH (u:Member{UserID: favoritedUserID})
	MERGE (u)-[:FAVORITED]->(m)
	WITH m
	MATCH (m2:Message)
	WHERE m2.CreatedAt < m.CreatedAt AND m2.GroupID = m.GroupID
	WITH m, m2 ORDER BY m2.CreatedAt DESC LIMIT 1
	MERGE (m2)-[:REPLIED_BY]->(m)`
	result, err := session.Run(query, map[string]interface{}{
		"ID":          m.ID,
		"SourceGUID":  m.SourceGUID,
		"CreatedAt":   m.CreatedAt,
		"UserID":      m.UserID,
		"GroupID":     m.GroupID,
		"Name":        m.Name,
		"AvatarURL":   m.AvatarURL,
		"Text":        m.Text,
		"System":      m.System,
		"FavoritedBy": m.FavoritedBy,
	})
	e := result.Err()
	if err != nil {
		log.Panic(err)
	} else if e != nil {
		log.Panic(e)
	}

	for _, a := range m.Attachments {
		query = `
			MATCH (m:Message{ID: $ID})
			MERGE (m)-[:ATTACHMENT]->(a:Attachment{
				Type: $Type,
				URL: $URL,
				Lat: $Lat,
				Lng: $Lng,
				Name: $Name,
				Token: $Token,
				Placeholder: $Placeholder
			})
			`
		result, err := session.Run(query, map[string]interface{}{
			"ID":          m.ID,
			"Type":        a.Type,
			"URL":         a.URL,
			"Lat":         a.Lat,
			"Lng":         a.Lng,
			"Name":        a.Name,
			"Token":       a.Token,
			"Placeholder": a.Placeholder,
		})
		e := result.Err()
		if err != nil {
			log.Panic(err)
		} else if e != nil {
			log.Panic(e)
		}
	}
}
