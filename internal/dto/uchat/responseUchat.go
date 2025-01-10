package dto

type ReponseUchatMid struct {
	Data []struct {
		ID      int    `json:"id"`
		Type    string `json:"type"`
		MsgType string `json:"msg_type"`
		Content string `json:"content"`
		Payload struct {
			Text string `json:"text"`
			URL  string `json:"url"`
		} `json:"payload"`
		SenderID string `json:"sender_id"`
		Ts       int    `json:"ts"`
	} `json:"data"`
}
