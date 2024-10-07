package model

type FileRule struct {
	Pattern  string `json:"pattern"`
	MinCount int    `json:"min_count"`
	MaxCount int    `json:"max_count"`
}
