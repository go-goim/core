package snowflake

var (
	defaultNode *Node
)

func SetDefaultNode(nodeBit int64) {
	var err error
	defaultNode, err = NewNode(nodeBit)
	if err != nil {
		panic(err)
	}
}

func assertDefaultNode() {
	if defaultNode == nil {
		panic("default node is not set")
	}
}

// Generate returns a snowflake ID
func Generate() ID {
	assertDefaultNode()
	return defaultNode.Generate()
}
