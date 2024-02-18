package id_utils

import (
	"log"

	"github.com/bwmarrin/snowflake"
)

// snowFlake is a presentation of [github.com/bwmarrin/snowflake] id generator.
type snowFlake struct {
	generator *snowflake.Node
}

// NewSnowFlake returns an [snowFlake] that implements [IDGenerator].
func NewSnowFlake(nodeID int64) IDGenerator {
	generator, err := snowflake.NewNode(nodeID)
	if err != nil {
		log.Fatalf("unable to create snowflake: %v", err)
	}
	return &snowFlake{
		generator: generator,
	}
}

func (s *snowFlake) String() string {
	return s.generator.Generate().String()
}

func (s *snowFlake) Int64() int64 {
	return s.generator.Generate().Int64()
}
