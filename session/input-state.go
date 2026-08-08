package session

type inputState struct {
	peer string
	key  []byte
	burn burnState
	text string
}

type burnState struct {
	key         []byte
	deletePaths []string
}

func (burn burnState) enabled() bool {
	return burn.key != nil
}
