package ledger

import (
	lkerr "github.com/victorotene80/LEDGERKIT/ledgerkit/errors"
)

// ValidateJournalEntry enforces LedgerKit invariants:
//
// 1) An entry must contain at least 2 postings.
// 2) All postings must use the same AssetCode (+ same Scale).
// 3) All posting amounts must be positive.
// 4) Sum(Debits) == Sum(Credits).
// 5) No zero-value postings.
func ValidateJournalEntry(je JournalEntry) error {
	postings := je.Postings() // defensive copy

	if len(postings) < 2 {
		return lkerr.ErrEntryNeedsTwoPostings
	}

	// Canonical asset+scale from the first posting.
	ref := postings[0].Money()
	refAsset := ref.Asset()
	refScale := ref.Scale()

	var debitTotal int64
	var creditTotal int64

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

		// Sum totals
		switch p.Side() {
		case SideDebit:
			debitTotal += m.Minor()
		case SideCredit:
			creditTotal += m.Minor()
		default:
			return lkerr.ErrInvalidSide
		}
	}

	if debitTotal != creditTotal {
		return lkerr.ErrUnbalancedEntry
	}

	return nil
}

// ValidatePosting enforces posting-level invariants:
// - account ref must be valid
// - side must be valid
// - money must be positive and non-zero
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

// ValidateAccountRef is intentionally conservative until you finalize AccountRef design.
// Once you switch to AccountKind string + OTHER requires code, tighten this.
func ValidateAccountRef(a AccountRef) error {
	// At minimum: require non-empty ID.
	if a.ID() == "" {
		return lkerr.ErrInvalidAccountRef
	}
	// TODO (after you paste account_ref.go): validate kind + OTHER code rules.
	return nil
}
