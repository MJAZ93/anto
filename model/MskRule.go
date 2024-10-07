package model

type MskMessage struct {
	MaxLines        int
	MinLines        int
	AllowedRules    []string
	DisallowedRules []string
	Comment         string
}
