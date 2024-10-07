package model

type FolderRule struct {
	Pattern  string       `json:"pattern"`
	Files    []FileRule   `json:"files"`
	Folders  []FolderRule `json:"folders"`
	MinCount int          `json:"min_count"`
	MaxCount int          `json:"max_count"`
}
