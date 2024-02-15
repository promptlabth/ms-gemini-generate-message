package promptMessage

import (
	"context"
	"database/sql"
	"log"
)

type promptMessageStorage struct {
	db *sql.DB
}

func NewPromptMessageStorage(db *sql.DB) *promptMessageStorage {
	return &promptMessageStorage{
		db: db,
	}
}

func (s promptMessageStorage) Save(ctx context.Context, promptMessage PromptMessage) (*PromptMessage, error) {
	rows, err := s.db.QueryContext(ctx,
		"INSERT INTO promptmessages (date_time, input_message, result_message, user_id, tone_id, feature_id, model_id) VALUES($1, $2, $3, $4, $5, $6 ,$7) RETURNING *",
		promptMessage.DateTime,
		promptMessage.InputMessage,
		promptMessage.ResultMessage,
		promptMessage.UserId,
		promptMessage.ToneId,
		promptMessage.FeatureId,
		promptMessage.ModelId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var promptResult PromptMessage
	for rows.Next() {
		if err := rows.Scan(
			&promptResult.ID,
			&promptResult.DateTime,
			&promptResult.InputMessage,
			&promptResult.ResultMessage,
			&promptResult.UserId,
			&promptResult.ToneId,
			&promptResult.FeatureId,
			&promptResult.ModelId,
		); err != nil {
			// Check for a scan error.
			// Query rows will be closed with defer.
			log.Fatal(err)
		}
	}

	rerr := rows.Close()
	if rerr != nil {
		return nil, rerr
	}

	// Rows.Err will report the last error encountered by Rows.Scan.
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &promptResult, nil

}
