package ledger

import lkerr "github.com/victorotene80/LEDGERKIT/ledgerkit/errors"

type Side uint8

const (
	SideUnknown Side = iota
	SideDebit
	SideCredit
)

func (s Side) String() string {
	switch s {
	case SideDebit:
		return "DEBIT"
	case SideCredit:
		return "CREDIT"
	default:
		return "UNKNOWN"
	}
}

func (s Side) Valid() bool {
	return s == SideDebit || s == SideCredit
}

func (s Side) Validate() error {
	if !s.Valid() {
		return lkerr.ErrInvalidSide
	}
	return nil
}
