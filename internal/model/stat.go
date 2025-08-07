package model

type MatchResult struct {
	ID       string `bson:"_id,omitempty"`
	Player   string `bson:"player"`
	Opponent string `bson:"opponent"`
	Result   string `bson:"result"`
}
