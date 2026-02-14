package ledger

import (
	"strings"

	lkerr "github.com/victorotene80/LEDGERKIT/ledgerkit/errors"
	"github.com/victorotene80/LEDGERKIT/ledgerkit/money"
)

func ValidateJournalEntry(je JournalEntry) error {
	postings := je.Postings() // defensive copy

	if len(postings) < 2 {
		return lkerr.ErrEntryNeedsTwoPostings
	}

	// Canonical asset+scale from the first posting.
	ref := postings[0].Money()
	refAsset := ref.Asset()
	refScale := ref.Scale()

	debits := money.MustNew(refAsset, refScale, 0)
	credits := money.MustNew(refAsset, refScale, 0)

	for _, p := range postings {
		// Posting-level invariants
		if err := ValidatePosting(p); err != nil {
			return err
		}

		m := p.Money()

		// Mixed assets / scales
		if m.Asset() != refAsset {
			return lkerr.ErrAssetMismatch
		}
		if m.Scale() != refScale {
			return lkerr.ErrScaleMismatch
		}

		// Sum totals (uses Money.Add -> keeps overflow handling consistent)
		var err error
		switch p.Side() {
		case SideDebit:
			debits, err = debits.Add(m)
		case SideCredit:
			credits, err = credits.Add(m)
		default:
			return lkerr.ErrInvalidSide
		}
		if err != nil {
			return err
		}
	}

	if debits.Minor() != credits.Minor() {
		return lkerr.ErrUnbalancedEntry
	}

	return nil
}

func ValidatePosting(p Posting) error {
	if err := ValidateAccountRef(p.Account()); err != nil {
		return err
	}

	if !p.Side().Valid() {
		return lkerr.ErrInvalidSide
	}

	m := p.Money()

	if m.IsZero() {
		return lkerr.ErrZeroPosting
	}

	if !m.IsPositive() {
		return lkerr.ErrPostingNotPositive
	}

	return nil
}

func ValidateAccountRef(a AccountRef) error {
	id := strings.TrimSpace(a.ID())
	if id == "" {
		return lkerr.ErrInvalidAccountRef
	}

	k := a.Kind()
	if !k.Valid() {
		return lkerr.ErrInvalidAccountRef
	}

	code := strings.TrimSpace(a.Code())
	if k == KindOther {
		// REQUIRED when typ == OTHER
		if code == "" {
			return lkerr.ErrInvalidAccountRef
		}
	} else {
		// Optional strictness: if not OTHER, you can enforce code must be empty.
		// This keeps Key() stable and prevents weird "USER:abc:user1" keys.
		if code != "" {
			return lkerr.ErrInvalidAccountRef
		}
	}

	return nil
}
