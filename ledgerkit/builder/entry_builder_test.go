package builder

import (
	"testing"
	"time"

	lkerr "github.com/victorotene80/LEDGERKIT/ledgerkit/errors"
	"github.com/victorotene80/LEDGERKIT/ledgerkit/ledger"
	"github.com/victorotene80/LEDGERKIT/ledgerkit/money"
)

func TestEntryBuilder_BuildBalanced(t *testing.T) {
	b := MustNewEntryBuilder("e1", "ext1", time.Now())

	user := ledger.MustAccountRef(ledger.KindUser, "", "user_1")
	clearing := ledger.MustAccountRef(ledger.KindClearing, "", "main")

	amt := money.MustNew("NGN", 2, 10000)

	_, err := b.Debit(user, amt)
	if err != nil {
		t.Fatalf("debit unexpected error: %v", err)
	}
	_, err = b.Credit(clearing, amt)
	if err != nil {
		t.Fatalf("credit unexpected error: %v", err)
	}

	je, err := b.Build()
	if err != nil {
		t.Fatalf("expected valid entry, got %v", err)
	}

	if len(je.Postings()) != 2 {
		t.Fatalf("expected 2 postings")
	}
}

func TestEntryBuilder_MixedAssetsRejectedEarly(t *testing.T) {
	b := MustNewEntryBuilder("e1", "ext1", time.Now())

	user := ledger.MustAccountRef(ledger.KindUser, "", "user_1")

	ngn := money.MustNew("NGN", 2, 10000)
	usd := money.MustNew("USD", 2, 10000)

	_, err := b.Debit(user, ngn)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	_, err = b.Credit(user, usd)
	if err != lkerr.ErrAssetMismatch {
		t.Fatalf("expected ErrAssetMismatch, got %v", err)
	}
}

func TestEntryBuilder_UnbalancedRejectedOnBuild(t *testing.T) {
	b := MustNewEntryBuilder("e1", "ext1", time.Now())

	user := ledger.MustAccountRef(ledger.KindUser, "", "user_1")
	clearing := ledger.MustAccountRef(ledger.KindClearing, "", "main")

	debit := money.MustNew("NGN", 2, 10000)
	credit := money.MustNew("NGN", 2, 9000)

	_, _ = b.Debit(user, debit)
	_, _ = b.Credit(clearing, credit)

	_, err := b.Build()
	if err != lkerr.ErrUnbalancedEntry {
		t.Fatalf("expected ErrUnbalancedEntry, got %v", err)
	}
}

func TestEntryBuilder_NeedsTwoPostings(t *testing.T) {
	b := MustNewEntryBuilder("e1", "ext1", time.Now())

	user := ledger.MustAccountRef(ledger.KindUser, "", "user_1")
	amt := money.MustNew("NGN", 2, 10000)

	_, _ = b.Debit(user, amt)

	_, err := b.Build()
	if err != lkerr.ErrEntryNeedsTwoPostings {
		t.Fatalf("expected ErrEntryNeedsTwoPostings, got %v", err)
	}
}
