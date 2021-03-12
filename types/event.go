package types

import "github.com/gofrs/uuid"

// Event is what an apply method accepts
type Event struct {
	Interaction *Interaction
	Session     *UserSession
	Entity      uuid.UUID
	Origin      uuid.UUID
	Endpoint    uuid.UUID
}
