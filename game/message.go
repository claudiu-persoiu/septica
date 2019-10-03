package game

type message struct {
	Action string `json:"action"`
	Data   string `json:"data,omitempty"`
}
