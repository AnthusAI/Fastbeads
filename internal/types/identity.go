package types

import (
	"crypto/rand"
	"fmt"
)

// EnsureIdentity populates UUID and DisplayID if missing.
// UUID is canonical; DisplayID defaults to ID when present.
func (i *Issue) EnsureIdentity() error {
	if err := i.ensureUUID(); err != nil {
		return err
	}
	if i.DisplayID == "" && i.ID != "" {
		i.DisplayID = i.ID
	}
	return nil
}

func (i *Issue) ensureUUID() error {
	if i.UUID != "" {
		return nil
	}
	uuid, err := newUUID()
	if err != nil {
		return err
	}
	i.UUID = uuid
	return nil
}

// newUUID generates an RFC 4122 UUIDv4.
func newUUID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	// Set version (4) and variant (RFC 4122).
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4],
		b[4:6],
		b[6:8],
		b[8:10],
		b[10:16],
	), nil
}
