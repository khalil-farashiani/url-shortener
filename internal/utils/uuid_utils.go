package utils

import "github.com/google/uuid"

func generateUniqueString() string {
	id := uuid.New()
	return id.String()
}
