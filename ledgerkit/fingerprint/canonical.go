package fingerprint

import (
	"bytes"
	"sort"
	"time"

	"github.com/victorotene80/LEDGERKIT/ledgerkit/ledger"
)

// CanonicalEntryBytes returns a deterministic, order-independent serialization of a JournalEntry.
// Format (versioned):
// v1|id=...|ext=...|ts=...|n=...|p=...|p=...
//
// Each posting is canonicalized and SORTED so fingerprint doesn't depend on input slice order.
func CanonicalEntryBytes(je ledger.JournalEntry) []byte {
	postings := je.Postings() // defensive copy already (per your JE design)

	type canonPosting struct {
		s string
	}

	cp := make([]canonPosting, 0, len(postings))
	for _, p := range postings {
		m := p.Money()

		// Posting encoding (stable, no spaces):
		// acct=<AccountRef.Key()>;side=<D|C>;asset=<ASSET>;scale=<S>;minor=<N>
		side := "?"
		switch p.Side() {
		case ledger.SideDebit:
			side = "D"
		case ledger.SideCredit:
			side = "C"
		}

		cp = append(cp, canonPosting{
			s: "acct=" + p.Account().Key() +
				";side=" + side +
				";asset=" + string(m.Asset()) +
				";scale=" + itoaU8(m.Scale()) +
				";minor=" + itoaI64(m.Minor()),
		})
	}

	sort.Slice(cp, func(i, j int) bool { return cp[i].s < cp[j].s })

	ts := je.Timestamp().UTC().Format(time.RFC3339Nano)

	var b bytes.Buffer
	b.Grow(128 + len(cp)*64)

	b.WriteString("v1")
	b.WriteString("|id=")
	b.WriteString(je.EntryID())
	b.WriteString("|ext=")
	b.WriteString(je.ExternalRef())
	b.WriteString("|ts=")
	b.WriteString(ts)
	b.WriteString("|n=")
	b.WriteString(itoaI(len(cp)))

	for _, p := range cp {
		b.WriteString("|p=")
		b.WriteString(p.s)
	}

	return b.Bytes()
}
