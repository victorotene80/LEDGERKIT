package projection

import (
	lkerr "github.com/victorotene80/LEDGERKIT/ledgerkit/errors"
	"github.com/victorotene80/LEDGERKIT/ledgerkit/ledger"
	"github.com/victorotene80/LEDGERKIT/ledgerkit/money"
)

// BalanceMap holds projected balances per account.
// Key is ledger.AccountRef.Key() to stay stable and comparable.
type BalanceMap struct {
	asset money.AssetCode
	scale uint8
	set   bool

	m map[string]money.Money
}

func NewBalanceMap() *BalanceMap {
	return &BalanceMap{
		m: make(map[string]money.Money),
	}
}

func (bm *BalanceMap) Asset() money.AssetCode { return bm.asset }
func (bm *BalanceMap) Scale() uint8           { return bm.scale }
func (bm *BalanceMap) IsSet() bool            { return bm.set }

// Get returns balance for an account (zero if absent).
func (bm *BalanceMap) Get(a ledger.AccountRef) (money.Money, error) {
	if !bm.set {
		return money.Money{}, lkerr.ErrInvalidAmount
	}
	key := a.Key()
	if v, ok := bm.m[key]; ok {
		return v, nil
	}
	return money.MustNew(bm.asset, bm.scale, 0), nil
}

// ApplyEntry applies a validated journal entry into the projection.
func (bm *BalanceMap) ApplyEntry(e ledger.JournalEntry) error {
	// Enforce ledger invariants first.
	if err := ledger.ValidateJournalEntry(e); err != nil {
		return err
	}

	postings := e.Postings()
	ref := postings[0].Money()
	if !bm.set {
		bm.asset = ref.Asset()
		bm.scale = ref.Scale()
		bm.set = true
	} else {
		if ref.Asset() != bm.asset {
			return lkerr.ErrAssetMismatch
		}
		if ref.Scale() != bm.scale {
			return lkerr.ErrScaleMismatch
		}
	}

	for _, p := range postings {
		acctKey := p.Account().Key()

		cur, ok := bm.m[acctKey]
		if !ok {
			cur = money.MustNew(bm.asset, bm.scale, 0)
		}

		var next money.Money
		var err error

		switch p.Side() {
		case ledger.SideDebit:
			next, err = cur.Add(p.Money())
		case ledger.SideCredit:
			next, err = cur.Sub(p.Money())
		default:
			return lkerr.ErrInvalidSide
		}

		if err != nil {
			return err
		}
		bm.m[acctKey] = next
	}

	return nil
}

// ApplyEntries applies multiple entries in order.
func (bm *BalanceMap) ApplyEntries(entries []ledger.JournalEntry) error {
	for _, e := range entries {
		if err := bm.ApplyEntry(e); err != nil {
			return err
		}
	}
	return nil
}

func (bm *BalanceMap) MustGet(a ledger.AccountRef) money.Money {
	v, err := bm.Get(a)
	if err != nil {
		panic(err)
	}
	return v
}
