package utils

import (
	"fmt"

	"github.com/bwmarrin/snowflake"
)

func GenerateSnowflake() int64 {
	// Create a new Node with a Node number of 1
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	// Generate a snowflake ID.
	return node.Generate().Int64()
}
