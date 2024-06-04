package source

// Message ... TODO
// https://github.com/axelarnetwork/axelar-core/blob/c7afbd0703a1f9ee1af415490ba2e931fc2bc1f7/x/axelarnet/message_handler.go#L86-L93
type Message struct {
	DestinationChain   string `json:"destination_chain"`
	DestinationAddress string `json:"destination_address"`
	Payload            []byte `json:"payload"`
	Type               int64  `json:"type"`
	Fee                *Fee   `json:"fee"` // Optional
}

// Fee ... TODO
// https://github.com/axelarnetwork/axelar-core/blob/c7afbd0703a1f9ee1af415490ba2e931fc2bc1f7/x/axelarnet/message_handler.go#L24-L29
type Fee struct {
	Amount          string  `json:"amount"`
	Recipient       string  `json:"recipient"`
	RefundRecipient *string `json:"refund_recipient"`
}
