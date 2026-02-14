package main

import (
	"fmt"
	"time"

	"github.com/victorotene80/LEDGERKIT/ledgerkit/builder"
	"github.com/victorotene80/LEDGERKIT/ledgerkit/ledger"
	"github.com/victorotene80/LEDGERKIT/ledgerkit/money"
	"github.com/victorotene80/LEDGERKIT/ledgerkit/projection"
)

func main() {
	user := ledger.MustAccountRef(ledger.KindUser, "", "user_1")
	clearing := ledger.MustAccountRef(ledger.KindClearing, "", "main")

	ngn100 := money.MustNew("NGN", 2, 10000)

	e1, err := builder.MustNewEntryBuilder("e1", "ext_1", time.Now()).
		MustDebit(user, ngn100).
		MustCredit(clearing, ngn100).
		Build()
	must(err)

	e2, err := builder.MustNewEntryBuilder("e2", "ext_2", time.Now()).
		MustDebit(user, ngn100).
		MustCredit(clearing, ngn100).
		Build()
	must(err)

	bm := projection.NewBalanceMap()
	must(bm.ApplyEntries([]ledger.JournalEntry{e1, e2}))

	fmt.Println("User balance:", bm.MustGet(user).String())
	fmt.Println("Clearing balance:", bm.MustGet(clearing).String())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
