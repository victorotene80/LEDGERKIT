package ledger

import (
	"time"

	lkerr "github.com/Helen-projects/LEDGERKIT/ledgerkit/errors"
	"github.com/Helen-projects/LEDGERKIT/ledgerkit/money"
)

type JournalEntry struct {
	entryID     string
	externalRef string
	ts          time.Time
	postings    []Posting
}

func NewJournalEntry(entryID, externalRef string, ts time.Time, postings []Posting) (JournalEntry, error) {
	je := JournalEntry{
		entryID:     entryID,
		externalRef: externalRef,
		ts:          ts.UTC(),
		postings:    clonePostings(postings),
	}
	if err := je.Validate(); err != nil {
		return JournalEntry{}, err
	}
	return je, nil
}

func (je JournalEntry) EntryID() string      { return je.entryID }
func (je JournalEntry) ExternalRef() string  { return je.externalRef }
func (je JournalEntry) Timestamp() time.Time { return je.ts }

func (je JournalEntry) Postings() []Posting { return clonePostings(je.postings) }

// Keep this convenience method.
func (je JournalEntry) Validate() error { return ValidateJournalEntry(je) }

func (je JournalEntry) Totals() (debits money.Money, credits money.Money, err error) {
	if len(je.postings) == 0 {
		return money.Money{}, money.Money{}, lkerr.ErrEntryNeedsTwoPostings
	}

	ref := je.postings[0].Money()
	d := money.MustNew(ref.Asset(), ref.Scale(), 0)
	c := money.MustNew(ref.Asset(), ref.Scale(), 0)

	for _, p := range je.postings {
		switch p.Side() {
		case SideDebit:
			d, err = d.Add(p.Money())
		case SideCredit:
			c, err = c.Add(p.Money())
		default:
			return money.Money{}, money.Money{}, lkerr.ErrInvalidSide
		}
		if err != nil {
			return money.Money{}, money.Money{}, err
		}
	}
	return d, c, nil
}

func clonePostings(in []Posting) []Posting {
	if len(in) == 0 {
		return nil
	}
	out := make([]Posting, len(in))
	copy(out, in)
	return out
}
