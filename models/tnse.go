package models

type TnseEmbeddingRecord struct {
	RecordId        int64   `json:"id"`
	Label           int64   `json:"label"`
	PredictionLabel int64   `json:"pred_label"`
	Correct         string  `json:"correct"`
	TnseX           float64 `json:"tnse_x"`
	TnseY           float64 `json:"tnse_y"`
}
