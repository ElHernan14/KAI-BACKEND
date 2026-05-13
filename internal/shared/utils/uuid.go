package utils

import "github.com/google/uuid"

func GenerateUUID() string {
	return uuid.NewString()
}

func ParseUUID(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}
