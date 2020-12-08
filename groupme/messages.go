package groupme

import (
	"fmt"

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
	Type string `json:"type"`
}

// AttachmentImage is a attachment containing a single image.
type AttachmentImage struct {
	Attachment
	URL string `json:"url"`
}

// AttachmentLocation contains a location.
type AttachmentLocation struct {
	Attachment
	Lat  string `json:"lat"`
	Lng  string `json:"lng"`
	Name string `json:"name"`
}

// AttachmentSplit is unknown.
type AttachmentSplit struct {
	Attachment
	Token string `json:"token"`
}

// AttachmentEmoji attaches a GroupMe emoji.
type AttachmentEmoji struct {
	Attachment
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
	if err != nil {
		panic(err)
	}
	return messages.Messages
}

// SaveToNeo4j saves the current message into the database.
func (m *Message) SaveToNeo4j(driver *database.Neo4j) {
	session, err := driver.NewWriteSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	result, err := session.Run(fmt.Sprintf(`MERGE (msg:Message{%s})`, Melt(*m)), map[string]interface{}{})
	e := result.Err()
	if err != nil {
		panic(err)
	} else if e != nil {
		// panic(e)
	}
}
