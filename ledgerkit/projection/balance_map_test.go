package projection

import (
	"testing"
	"time"

	"github.com/victorotene80/LEDGERKIT/ledgerkit/ledger"
	"github.com/victorotene80/LEDGERKIT/ledgerkit/money"
)

func TestBalanceMap_ApplyEntry_SingleEntry(t *testing.T) {
	bm := NewBalanceMap()

	user := ledger.MustAccountRef(ledger.KindUser, "", "user_1")
	clearing := ledger.MustAccountRef(ledger.KindClearing, "", "main")

	amt := money.MustNew("NGN", 2, 10000) // 100.00

	e, err := ledger.NewJournalEntry("e1", "ext1", time.Now(), []ledger.Posting{
		ledger.MustPosting(user, ledger.SideDebit, amt),
		ledger.MustPosting(clearing, ledger.SideCredit, amt),
	})
	if err != nil {
		t.Fatalf("entry should be valid: %v", err)
	}

	if err := bm.ApplyEntry(e); err != nil {
		t.Fatalf("apply should succeed: %v", err)
	}

	u, _ := bm.Get(user)
	c, _ := bm.Get(clearing)

	if u.Minor() != 10000 {
		t.Fatalf("expected user +10000, got %d", u.Minor())
	}
	if c.Minor() != -10000 {
		t.Fatalf("expected clearing -10000, got %d", c.Minor())
	}
}

func TestBalanceMap_ApplyEntries_MultipleEntries(t *testing.T) {
	bm := NewBalanceMap()

	user := ledger.MustAccountRef(ledger.KindUser, "", "user_1")
	clearing := ledger.MustAccountRef(ledger.KindClearing, "", "main")
	fee := ledger.MustAccountRef(ledger.KindFee, "", "fee_pool")

	amt100 := money.MustNew("NGN", 2, 10000)
	feeAmt := money.MustNew("NGN", 2, 200) // 2.00

	// Entry 1: user deposits 100.00
	e1, _ := ledger.NewJournalEntry("e1", "ext1", time.Now(), []ledger.Posting{
		ledger.MustPosting(user, ledger.SideDebit, amt100),
		ledger.MustPosting(clearing, ledger.SideCredit, amt100),
	})

	// Entry 2: fee charge 2.00 (user pays fee pool)
	e2, _ := ledger.NewJournalEntry("e2", "ext2", time.Now(), []ledger.Posting{
		ledger.MustPosting(user, ledger.SideCredit, feeAmt), // -2.00
		ledger.MustPosting(fee, ledger.SideDebit, feeAmt),   // +2.00
	})

	if err := bm.ApplyEntries([]ledger.JournalEntry{e1, e2}); err != nil {
		t.Fatalf("apply entries should succeed: %v", err)
	}

	u, _ := bm.Get(user)
	c, _ := bm.Get(clearing)
	f, _ := bm.Get(fee)

	if u.Minor() != 9800 {
		t.Fatalf("expected user 9800, got %d", u.Minor())
	}
	if c.Minor() != -10000 {
		t.Fatalf("expected clearing -10000, got %d", c.Minor())
	}
	if f.Minor() != 200 {
		t.Fatalf("expected fee +200, got %d", f.Minor())
	}
}

func TestBalanceMap_RejectsMixedAssetsAcrossEntries(t *testing.T) {
	bm := NewBalanceMap()

	user := ledger.MustAccountRef(ledger.KindUser, "", "user_1")
	clearing := ledger.MustAccountRef(ledger.KindClearing, "", "main")

	ngn := money.MustNew("NGN", 2, 10000)
	usd := money.MustNew("USD", 2, 10000)

	e1, _ := ledger.NewJournalEntry("e1", "ext1", time.Now(), []ledger.Posting{
		ledger.MustPosting(user, ledger.SideDebit, ngn),
		ledger.MustPosting(clearing, ledger.SideCredit, ngn),
	})
	e2, _ := ledger.NewJournalEntry("e2", "ext2", time.Now(), []ledger.Posting{
		ledger.MustPosting(user, ledger.SideDebit, usd),
		ledger.MustPosting(clearing, ledger.SideCredit, usd),
	})

	if err := bm.ApplyEntry(e1); err != nil {
		t.Fatalf("e1 should apply: %v", err)
	}
	if err := bm.ApplyEntry(e2); err == nil {
		t.Fatalf("expected mixed asset error")
	}
}
