package security

import "github.com/google/uuid"

type Token struct {
	ID          uuid.UUID
	AccessLevel int
}

func (t Token) Valid() error {
	if t.AccessLevel < 0 {
		panic("Illegal accessLevel")
	}
	return nil
}
