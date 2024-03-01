package promptMessage

import "time"

type PromptMessage struct {
	ID            uint      `json:"id"`
	DateTime      time.Time `json:"dateTime"`
	InputMessage  string    `json:"inputMessage"`
	ResultMessage string    `json:"resultMessage"`
	UserId        uint      `json:"userId"`
	ToneId        uint      `json:"toneId"`
	FeatureId     uint      `json:"featureId"`
	ModelId       uint      `json:"modelId"`
}

type PromptMessageRequest struct {
	InputMessage string `json:"inputMessage"`
}
