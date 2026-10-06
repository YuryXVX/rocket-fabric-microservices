package application

import "github.com/google/uuid"

func validateUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
