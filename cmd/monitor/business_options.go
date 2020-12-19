package main

import (
	"strings"

	"github.com/stehrn/hpc-poc/internal/utils"
)

var businessNames []string

func init() {
	businessNames = strings.Split(utils.Env("BUSINESS_NAMES"), ",")
}

// BusinessNameOptions used for populating business option widget
type BusinessNameOptions struct {
	Name     string
	Selected bool
}

// NewBusinessNameOptions create new slice of BusinessNameOptions
func NewBusinessNameOptions(selected string) []BusinessNameOptions {
	options := make([]BusinessNameOptions, len(businessNames))
	for i, name := range businessNames {
		var isOptSelected bool
		if name == selected {
			isOptSelected = true
		}
		options[i] = BusinessNameOptions{name, isOptSelected}
	}
	return options
}
