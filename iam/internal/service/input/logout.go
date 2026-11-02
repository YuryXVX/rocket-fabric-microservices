package input

import "github.com/google/uuid"

type LogoutInput struct {
	SessionUUID uuid.UUID
}
