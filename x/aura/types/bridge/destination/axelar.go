package destination

// Message ... TODO
// https://github.com/axelarnetwork/axelar-core/blob/c7afbd0703a1f9ee1af415490ba2e931fc2bc1f7/x/axelarnet/types/evm_translator.go#L55-L60
type Message struct {
	SourceChain   string `json:"source_chain"`
	SourceAddress string `json:"source_address"`
	Payload       []byte `json:"payload"`
	Type          int64  `json:"type"`
}
