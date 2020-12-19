package utils

import (
	"github.com/rs/xid"
)

// GenerateID generate unique ID
func GenerateID() string {
	return xid.New().String()
}
