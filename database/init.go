package database

// Init prepares the Neo4j database with constraints.
func Init() {
	n, err := NewNeo4j("bolt://localhost:7687", "", "", false)
	if err != nil {
		panic(err)
	}

	session, err := n.NewWriteSession()
	if err != nil {
		panic(err)
	}
	defer n.driver.Close()
	defer session.Close()

	_, err = session.Run("CREATE CONSTRAINT groupIDUnique IF NOT EXISTS ON (n:Group) ASSERT n.ID IS UNIQUE", map[string]interface{}{})
	if err != nil {
		panic(err)
	}

	_, err = session.Run("CREATE CONSTRAINT messageIDUnique IF NOT EXISTS ON (n:Message) ASSERT n.ID IS UNIQUE", map[string]interface{}{})
	if err != nil {
		panic(err)
	}

	_, err = session.Run("CREATE CONSTRAINT userIDUnique IF NOT EXISTS ON (n:Member) ASSERT n.UserID IS UNIQUE", map[string]interface{}{})
	if err != nil {
		panic(err)
	}
}
