package fingerprint

import (
	"testing"
	"time"

	"github.com/victorotene80/LEDGERKIT/ledgerkit/ledger"
	"github.com/victorotene80/LEDGERKIT/ledgerkit/money"
)

func TestFingerprint_DeterministicSameEntry(t *testing.T) {
	user := ledger.MustAccountRef(ledger.KindUser, "", "user_1")
	clear := ledger.MustAccountRef(ledger.KindClearing, "", "main")
	amt := money.MustNew("NGN", 2, 10000)

	p1 := ledger.MustPosting(user, ledger.SideDebit, amt)
	p2 := ledger.MustPosting(clear, ledger.SideCredit, amt)

	ts := time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC)
	je, err := ledger.NewJournalEntry("e1", "ext_1", ts, []ledger.Posting{p1, p2})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	a := EntryFingerprint(je)
	b := EntryFingerprint(je)

	if a != b {
		t.Fatalf("expected same fingerprint, got %s vs %s", a.Hex(), b.Hex())
	}
}

func TestFingerprint_OrderIndependentPostings(t *testing.T) {
	user := ledger.MustAccountRef(ledger.KindUser, "", "user_1")
	clear := ledger.MustAccountRef(ledger.KindClearing, "", "main")
	amt := money.MustNew("NGN", 2, 10000)

	p1 := ledger.MustPosting(user, ledger.SideDebit, amt)
	p2 := ledger.MustPosting(clear, ledger.SideCredit, amt)

	ts := time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC)

	je1, err := ledger.NewJournalEntry("e1", "ext_1", ts, []ledger.Posting{p1, p2})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	je2, err := ledger.NewJournalEntry("e1", "ext_1", ts, []ledger.Posting{p2, p1}) // reversed
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	f1 := EntryFingerprint(je1)
	f2 := EntryFingerprint(je2)

	if f1 != f2 {
		t.Fatalf("expected same fingerprint despite posting order, got %s vs %s", f1.Hex(), f2.Hex())
	}
}

func TestFingerprint_ChangesIfAnyFieldChanges(t *testing.T) {
	user := ledger.MustAccountRef(ledger.KindUser, "", "user_1")
	clear := ledger.MustAccountRef(ledger.KindClearing, "", "main")
	amt := money.MustNew("NGN", 2, 10000)

	p1 := ledger.MustPosting(user, ledger.SideDebit, amt)
	p2 := ledger.MustPosting(clear, ledger.SideCredit, amt)

	ts := time.Date(2026, 2, 14, 0, 0, 0, 0, time.UTC)

	jeA, _ := ledger.NewJournalEntry("e1", "ext_1", ts, []ledger.Posting{p1, p2})
	jeB, _ := ledger.NewJournalEntry("e1", "ext_CHANGED", ts, []ledger.Posting{p1, p2})

	fA := EntryFingerprint(jeA)
	fB := EntryFingerprint(jeB)

	if fA == fB {
		t.Fatalf("expected different fingerprint when external ref changes")
	}
}
