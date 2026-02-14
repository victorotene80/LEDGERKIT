package ledger

import (
	"testing"
	"time"

	"github.com/victorotene80/LEDGERKIT/ledgerkit/money"
)

func TestJournalEntry_Balanced(t *testing.T) {
	user := MustAccountRef(KindUser, "", "user_1")
	clearing := MustAccountRef(KindClearing, "", "main")

	amt := money.MustNew("NGN", 2, 10000) // 100.00 NGN

	p1 := MustPosting(user, SideDebit, amt)
	p2 := MustPosting(clearing, SideCredit, amt)

	je, err := NewJournalEntry("e1", "ext_1", time.Now(), []Posting{p1, p2})
	if err != nil {
		t.Fatalf("expected valid entry, got error: %v", err)
	}
	if len(je.Postings()) != 2 {
		t.Fatalf("expected 2 postings")
	}
}

func TestJournalEntry_Unbalanced(t *testing.T) {
	user := MustAccountRef(KindUser, "", "user_1")
	clearing := MustAccountRef(KindClearing, "", "main")

	debit := money.MustNew("NGN", 2, 10000)
	credit := money.MustNew("NGN", 2, 9000)

	p1 := MustPosting(user, SideDebit, debit)
	p2 := MustPosting(clearing, SideCredit, credit)

	_, err := NewJournalEntry("e1", "ext_1", time.Now(), []Posting{p1, p2})
	if err == nil {
		t.Fatalf("expected error for unbalanced entry")
	}
}

func TestJournalEntry_AssetMismatch(t *testing.T) {
	user := MustAccountRef(KindUser, "", "user_1")
	clearing := MustAccountRef(KindClearing, "", "main")

	ngn := money.MustNew("NGN", 2, 10000)
	usd := money.MustNew("USD", 2, 10000)

	p1 := MustPosting(user, SideDebit, ngn)
	p2 := MustPosting(clearing, SideCredit, usd)

	_, err := NewJournalEntry("e1", "ext_1", time.Now(), []Posting{p1, p2})
	if err == nil {
		t.Fatalf("expected error for asset mismatch")
	}
}

func TestJournalEntry_NeedsAtLeastTwoPostings(t *testing.T) {
	user := MustAccountRef(KindUser, "", "user_1")
	amt := money.MustNew("NGN", 2, 10000)

	p1 := MustPosting(user, SideDebit, amt)

	_, err := NewJournalEntry("e1", "ext_1", time.Now(), []Posting{p1})
	if err == nil {
		t.Fatalf("expected error for <2 postings")
	}
}

func TestAccountRef_OtherRequiresCode(t *testing.T) {
	defer func() { _ = recover() }()

	// should panic because OTHER requires code
	_ = MustAccountRef(KindOther, "", "id_1")

	t.Fatalf("expected panic for OTHER without code")
}
