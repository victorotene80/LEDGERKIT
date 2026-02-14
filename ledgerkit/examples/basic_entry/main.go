package main

import (
	"fmt"
	"time"

	"github.com/victorotene80/LEDGERKIT/ledgerkit/builder"
	"github.com/victorotene80/LEDGERKIT/ledgerkit/ledger"
	"github.com/victorotene80/LEDGERKIT/ledgerkit/money"
)

func main() {
	user := ledger.MustAccountRef(ledger.KindUser, "", "user_1")
	clearing := ledger.MustAccountRef(ledger.KindClearing, "", "main")

	amt := money.MustNew("NGN", 2, 10000) // 100.00 NGN

	b := builder.MustNewEntryBuilder("e1", "ext_1", time.Now()).
		MustDebit(user, amt).
		MustCredit(clearing, amt)

	je, err := b.Build()
	must(err)

	debits, credits, err := je.Totals()
	must(err)

	fmt.Println("EntryID:", je.EntryID())
	fmt.Println("ExternalRef:", je.ExternalRef())
	fmt.Println("Timestamp:", je.Timestamp().Format(time.RFC3339Nano))
	fmt.Println("Postings:", len(je.Postings()))
	fmt.Println("Debits:", debits.String())
	fmt.Println("Credits:", credits.String())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
