#!/bin/sh

docker run --publish=7474:7474 --publish=7687:7687 --env=NEO4J_AUTH=none --volume="./groupme-graph/data:/var/lib/neo4j/data" neo4j
