package groupme

import (
	"fmt"
	"log"

	"patrickwthomas.net/groupme-graph/database"
)

// Group is the JSON response format for a group.
type Group struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Type          string   `json:"type"`
	Description   string   `json:"description"`
	ImageURL      string   `json:"image_url"`
	CreatorUserID string   `json:"creator_user_id"`
	CreatedAt     int      `json:"created_at"`
	UpdatedAt     int      `json:"updated_at"`
	Members       []Member `json:"members"`
	ShareURL      string   `json:"share_url"`
	Messages      struct {
		Count                int    `json:"count"`
		LastMessageID        string `json:"last_message_id"`
		LastMessageCreatedAt int    `json:"last_message_created_at"`
		Preview              struct {
			Nickname    string `json:"nickname"`
			Text        string `json:"text"`
			ImageURL    string `json:"image_url"`
			Attachments []struct {
				Type        string  `json:"type"`
				URL         string  `json:"url,omitempty"`
				Lat         string  `json:"lat,omitempty"`
				Lng         string  `json:"lng,omitempty"`
				Name        string  `json:"name,omitempty"`
				Token       string  `json:"token,omitempty"`
				Placeholder string  `json:"placeholder,omitempty"`
				Charmap     [][]int `json:"charmap,omitempty"`
			} `json:"attachments"`
		} `json:"preview"`
	} `json:"messages"`
}

type changeOwner struct {
	GroupID string `json:"group_id"`
	OwnerID string `json:"owner_id"`
}

type createGroup struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Share       bool   `json:"share"`
	ImageURL    string `json:"image_url"`
}

// GroupsIndex gets the groups index from GroupMe.
func (g *GroupMe) GroupsIndex(page, perPage int, omitMemberships bool) ([]Group, error) {
	groups := &[]Group{}
	urlValues := map[string]string{
		"page":     fmt.Sprint(page),
		"per_page": fmt.Sprint(perPage),
	}
	if omitMemberships {
		urlValues["omit"] = "memberships"
	}
	_, err := g.groupMeRequest("GET", "/groups", urlValues, groups)
	if err != nil {
		return nil, err
	}
	return *groups, err
}

// GroupsFormer gets the groups index from GroupMe.
func (g *GroupMe) GroupsFormer() ([]Group, error) {
	groups := &[]Group{}
	_, err := g.groupMeRequest("GET", "/groups/former", nil, groups)
	if err != nil {
		return nil, err
	}
	return *groups, nil
}

// GroupsShow shows detail about a single group.
func (g *GroupMe) GroupsShow(groupID string) (*Group, error) {
	group := &Group{}
	_, err := g.groupMeRequest("GET", "/groups/"+groupID, nil, group)
	if err != nil {
		return nil, err
	}
	return group, nil
}

// GroupsCreate creates a new group.
func (g *GroupMe) GroupsCreate(name, description, imageURL string, share bool) (*Group, error) {
	group := &Group{}
	urlValues := createGroup{
		Name:        name,
		Description: description,
		ImageURL:    imageURL,
		Share:       share,
	}
	// _, err := g.groupMeRequest("POST", "/groups/create", urlValues, group)
	_, err := g.groupMeRequestPostObject("/groups", urlValues, group)
	if err != nil {
		return nil, err
	}
	return group, nil
}

// GroupsUpdate updates a group's information.
func (g *GroupMe) GroupsUpdate(name, description, imageURL string, officeMode, share bool) (*Group, error) {
	group := &Group{}
	urlValues := map[string]string{
		"name":        name,
		"description": description,
		"image_url":   imageURL,
		"office_mode": fmt.Sprint(officeMode),
		"share":       fmt.Sprint(share),
	}
	_, err := g.groupMeRequest("POST", "/groups/update", urlValues, group)
	if err != nil {
		log.Panic(err)
	}
	return group, nil
}

// GroupsDestroy deletes a group from GroupMe.
func (g *GroupMe) GroupsDestroy(groupID string) error {
	_, err := g.groupMeRequest("POST", "/groups/"+groupID+"/destroy", nil, nil)
	return err
}

// GroupsJoin updates a group's information.
func (g *GroupMe) GroupsJoin(groupID string, shareID string) (*Group, error) {
	group := &Group{}
	_, err := g.groupMeRequest("POST", "/groups/"+groupID+"/join/"+shareID, nil, group)
	if err != nil {
		log.Panic(err)
	}
	return group, nil
}

// GroupsRejoin rejoins a group that the user previously left.
func (g *GroupMe) GroupsRejoin(groupID string) (*Group, error) {
	group := &Group{}
	_, err := g.groupMeRequest("POST", "/groups/join", nil, group)
	if err != nil {
		log.Panic(err)
	}
	return group, nil
}

// GroupsChangeOwners rejoins a group that the user previously left.
func (g *GroupMe) GroupsChangeOwners(groupID, ownerID string) (*Group, error) {
	group := &Group{}
	urlValues := changeOwner{
		GroupID: groupID,
		OwnerID: ownerID,
	}
	_, err := g.groupMeRequestPostObject("/groups/change_owners", urlValues, group)
	if err != nil {
		log.Panic(err)
	}
	return group, nil
}

// SaveToNeo4j saves the current group into the database.
func (g *Group) SaveToNeo4j(driver *database.Neo4j) {
	session, err := driver.NewWriteSession()
	if err != nil {
		log.Panic(err)
	}
	defer session.Close()

	result, err := session.Run(fmt.Sprintf("MERGE (n:Group{%s})", Melt(*g)), map[string]interface{}{})
	e := result.Err()
	if err != nil {
		log.Panic(err)
	} else if e != nil {
		// panic(e)
	}

	session.Close()

	if len(g.Members) > 0 {
		for _, member := range g.Members {
			member.SaveToNeo4j(driver)
		}
	}
}
