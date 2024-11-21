package dto

type SlackMessagePayload struct {
	Text   string              `json:"text"`
	Blocks []SlackMessageBlock `json:"blocks,omitempty"`
}

type SlackMessageBlock struct {
	Type string            `json:"type"`
	Text *SlackMessageText `json:"text,omitempty"`
}

type SlackMessageText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}
