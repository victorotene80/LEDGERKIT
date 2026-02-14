package fingerprint

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/victorotene80/LEDGERKIT/ledgerkit/ledger"
)

type Fingerprint [32]byte

func (f Fingerprint) Hex() string { return hex.EncodeToString(f[:]) }

func EntryFingerprint(je ledger.JournalEntry) Fingerprint {
	sum := sha256.Sum256(CanonicalEntryBytes(je))
	return Fingerprint(sum)
}
