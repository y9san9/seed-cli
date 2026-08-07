package session

type message struct {
	_       struct{} `cbor:",toarray"`
	Type    string   // message
	Version int
	Payload []byte
	Burn    bool
}

type burnRequest struct {
	_         struct{} `cbor:",toarray"`
	Type      string   // burn_request
	Version   int
	PublicKey []byte
}

type burnConfirmation struct {
	_         struct{} `cbor:",toarray"`
	Type      string   // burn_confirmation
	Version   int
	PublicKey []byte
}

type fileMetadata struct {
	_       struct{} `cbor:",toarray"`
	Type    string   // "encrypted"
	Version int
	Name    string
}

type burnFileMetadata struct {
	_       struct{} `cbor:",toarray"`
	Type    string   // "burn"
	Version int
}
