package main

type Signal struct {
	ID         string `json:"id"`
	Block      string `json:"block"`
	Aspect     string `json:"aspect"`
	Inspection string `json:"inspection"`
}

var inspectionTransitions = map[string]map[string]bool{
	"pending":         {"clear": true, "needs_attention": true},
	"needs_attention": {"cleared": true},
	"clear":           {},
	"cleared":         {},
}
