package builder

import (
	"strings"
	"time"

	lkerr "github.com/victorotene80/LEDGERKIT/ledgerkit/errors"
	"github.com/victorotene80/LEDGERKIT/ledgerkit/ledger"
	"github.com/victorotene80/LEDGERKIT/ledgerkit/money"
)

// EntryBuilder builds a valid ledger.JournalEntry using a fluent API.
// It does not bypass ledger invariants; Build() returns ledger.NewJournalEntry(...)
// which validates the final postings.
type EntryBuilder struct {
	entryID     string
	externalRef string
	ts          time.Time

	// Canonical asset/scale chosen from first Money.
	asset money.AssetCode
	scale uint8
	set   bool

	postings []ledger.Posting

	requirePostings bool
}

func NewEntryBuilder(entryID, externalRef string, ts time.Time) (*EntryBuilder, error) {
	entryID = strings.TrimSpace(entryID)
	externalRef = strings.TrimSpace(externalRef)

	if entryID == "" || len(entryID) > 128 {
		return nil, lkerr.ErrInvalidAmount // optional: ErrInvalidEntryID later
	}
	if externalRef == "" || len(externalRef) > 128 {
		return nil, lkerr.ErrInvalidAmount // optional: ErrInvalidExternalRef later
	}

	return &EntryBuilder{
		entryID:         entryID,
		externalRef:     externalRef,
		ts:              ts.UTC(),
		postings:        make([]ledger.Posting, 0, 4),
		requirePostings: true,
	}, nil
}

func MustNewEntryBuilder(entryID, externalRef string, ts time.Time) *EntryBuilder {
	b, err := NewEntryBuilder(entryID, externalRef, ts)
	if err != nil {
		panic(err)
	}
	return b
}

// Debit adds a debit posting.
func (b *EntryBuilder) Debit(acct ledger.AccountRef, amt money.Money) (*EntryBuilder, error) {
	return b.add(acct, ledger.SideDebit, amt)
}

// Credit adds a credit posting.
func (b *EntryBuilder) Credit(acct ledger.AccountRef, amt money.Money) (*EntryBuilder, error) {
	return b.add(acct, ledger.SideCredit, amt)
}

func (b *EntryBuilder) add(acct ledger.AccountRef, side ledger.Side, amt money.Money) (*EntryBuilder, error) {
	// Set canonical asset+scale from first posting.
	if !b.set {
		b.asset = amt.Asset()
		b.scale = amt.Scale()
		b.set = true
	} else {
		// Prevent mixed assets/scales early (fail fast).
		if amt.Asset() != b.asset {
			return nil, lkerr.ErrAssetMismatch
		}
		if amt.Scale() != b.scale {
			return nil, lkerr.ErrScaleMismatch
		}
	}

	p, err := ledger.NewPosting(acct, side, amt)
	if err != nil {
		return nil, err
	}

	b.postings = append(b.postings, p)
	return b, nil
}

// Build finalizes the entry and enforces invariants via ledger.NewJournalEntry.
func (b *EntryBuilder) Build() (ledger.JournalEntry, error) {
	if b.requirePostings && len(b.postings) == 0 {
		return ledger.JournalEntry{}, lkerr.ErrEntryNeedsTwoPostings
	}
	return ledger.NewJournalEntry(b.entryID, b.externalRef, b.ts, b.postings)
}

func (b *EntryBuilder) MustDebit(acct ledger.AccountRef, amt money.Money) *EntryBuilder {
	nb, err := b.Debit(acct, amt)
	if err != nil {
		panic(err)
	}
	return nb
}

func (b *EntryBuilder) MustCredit(acct ledger.AccountRef, amt money.Money) *EntryBuilder {
	nb, err := b.Credit(acct, amt)
	if err != nil {
		panic(err)
	}
	return nb
}
