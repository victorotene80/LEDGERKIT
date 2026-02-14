package main

import (
	"fmt"
	"time"

	"github.com/victorotene80/LEDGERKIT/ledgerkit/builder"
	"github.com/victorotene80/LEDGERKIT/ledgerkit/fingerprint"
	"github.com/victorotene80/LEDGERKIT/ledgerkit/ledger"
	"github.com/victorotene80/LEDGERKIT/ledgerkit/money"
)

func main() {
	user := ledger.MustAccountRef(ledger.KindUser, "", "user_1")
	clearing := ledger.MustAccountRef(ledger.KindClearing, "", "main")

	ngn100 := money.MustNew("NGN", 2, 10000)

	je, err := builder.MustNewEntryBuilder("e1", "ext_1", time.Now()).
		MustDebit(user, ngn100).
		MustCredit(clearing, ngn100).
		Build()
	must(err)

	fp := fingerprint.EntryFingerprint(je)
	fmt.Println("Fingerprint (hex):", fp.Hex())

	// or if you want the alias name:
	fp2 := fingerprint.FingerprintSHA256(je)
	fmt.Println("FingerprintSHA256 (hex):", fp2.Hex())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
