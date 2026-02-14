package ledger

import (
	"strings"

	lkerr "github.com/Helen-projects/LEDGERKIT/ledgerkit/errors"
)

// AccountRef is a logical identifier for an account (kind + optional code + id).
// It is NOT a database key.
type AccountRef struct {
	kind AccountKind
	code string // required when kind == KindOther
	id   string
}

func NewAccountRef(kind AccountKind, code string, id string) (AccountRef, error) {
	kind = AccountKind(strings.ToUpper(strings.TrimSpace(kind.String())))
	code = strings.ToUpper(strings.TrimSpace(code))
	id = strings.TrimSpace(id)

	if !kind.ValidBasic() {
		return AccountRef{}, lkerr.ErrInvalidAccountRef
	}

	// For non-OTHER kinds, code must be empty (keeps identity canonical).
	if kind != KindOther && code != "" {
		return AccountRef{}, lkerr.ErrInvalidAccountRef
	}

	// For OTHER, code is required (namespacing).
	if kind == KindOther {
		if code == "" || len(code) > 32 {
			return AccountRef{}, lkerr.ErrInvalidAccountRef
		}
		// reuse same char policy as kind
		ck := AccountKind(code)
		if !ck.ValidBasic() {
			return AccountRef{}, lkerr.ErrInvalidAccountRef
		}
	}

	if id == "" || len(id) > 128 {
		return AccountRef{}, lkerr.ErrInvalidAccountRef
	}

	return AccountRef{kind: kind, code: code, id: id}, nil
}

func MustAccountRef(kind AccountKind, code, id string) AccountRef {
	a, err := NewAccountRef(kind, code, id)
	if err != nil {
		panic(err)
	}
	return a
}

func (a AccountRef) Kind() AccountKind { return a.kind }
func (a AccountRef) Code() string      { return a.code }
func (a AccountRef) ID() string        { return a.id }

// Key is deterministic and collision-resistant for hashing/maps/fingerprints.
func (a AccountRef) Key() string {
	if a.kind == KindOther {
		return a.kind.String() + ":" + a.code + ":" + a.id
	}
	return a.kind.String() + ":" + a.id
}
