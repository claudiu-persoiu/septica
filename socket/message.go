package socket

type Message struct {
	Type     string
	Message []string `json:"commands,omitempty"`
}
