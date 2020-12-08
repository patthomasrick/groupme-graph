package database

import "github.com/neo4j/neo4j-go-driver/neo4j"

// Neo4j is the struct the manages the current connection to Neo4j.
type Neo4j struct {
	driver neo4j.Driver
}

// NewNeo4j constructs a new Neo4j state instance.
func NewNeo4j(uri, username, password string, encrypted bool) (*Neo4j, error) {
	n := new(Neo4j)
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""), func(c *neo4j.Config) {
		c.Encrypted = encrypted
	})
	if err != nil {
		return nil, err
	}
	n.driver = driver

	return n, nil
}

// NewReadSession gets a new read session from the Neo4j driver.
func (n *Neo4j) NewReadSession() (neo4j.Session, error) {
	session, err := n.driver.Session(neo4j.AccessModeRead)
	if err != nil {
		return nil, err
	}
	return session, nil
}

// NewWriteSession gets a new write session from the Neo4j driver.
func (n *Neo4j) NewWriteSession() (neo4j.Session, error) {
	session, err := n.driver.Session(neo4j.AccessModeWrite)
	if err != nil {
		return nil, err
	}
	return session, nil
}
