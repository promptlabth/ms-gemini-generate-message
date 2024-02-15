package promptMessage

import "time"

type PromptMessage struct {
	ID            uint      `json:"id" db:"id"`
	DateTime      time.Time `json:"dateTime" db:"date_time"`
	InputMessage  string    `json:"inputMessage" db:"input_message"`
	ResultMessage string    `json:"resultMessage" db:"result_message"`
	UserId        uint      `json:"userId" db:"user_id"`
	ToneId        uint      `json:"toneId" db:"tone_id"`
	FeatureId     uint      `json:"featureId" db:"feature_id"`
	ModelId       uint      `json:"modelId" db:"model_id"`
}
