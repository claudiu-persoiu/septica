package game

type Message struct {
	Action string `json:"action"`
	Data   string `json:"data,omitempty"`
}
