package types

import (
	"github.com/go-goim/core/pkg/types/snowflake"
)

var (
	defaultNode *snowflake.Node
)

func SetDefaultNode(nodeBit int64) {
	var err error
	defaultNode, err = snowflake.NewNode(nodeBit)
	if err != nil {
		panic(err)
	}
}

func assertDefaultNode() {
	if defaultNode == nil {
		panic("default node is not set")
	}
}
