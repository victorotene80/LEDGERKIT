package ledger

import "strings"

// AccountKind is an extensible namespace label.
// The ledger does NOT interpret meaning; it’s for stable identity + fingerprinting.
type AccountKind string

const (
	KindUser     AccountKind = "USER"
	KindSystem   AccountKind = "SYSTEM"
	KindFee      AccountKind = "FEE"
	KindClearing AccountKind = "CLEARING"
	KindOther    AccountKind = "OTHER"
)

func (k AccountKind) String() string { return string(k) }

// ValidBasic checks formatting rules only.
// NOTE: we do NOT restrict to only the known constants; that would limit a library.
func (k AccountKind) ValidBasic() bool {
	s := strings.TrimSpace(string(k))
	if len(s) < 2 || len(s) > 32 {
		return false
	}
	// A-Z 0-9 _ -
	for _, r := range s {
		if (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '_' || r == '-' {
			continue
		}
		return false
	}
	return true
}

func (k AccountKind) Valid() bool {
	switch k {
	case KindUser,
		KindSystem,
		KindFee,
		KindClearing,
		KindOther:
		return true
	default:
		return false
	}
}
