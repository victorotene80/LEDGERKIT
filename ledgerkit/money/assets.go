package money

import (
	"strings"

	"github.com/Helen-projects/LEDGERKIT/ledgerkit/errors"
)

type AssetCode string

func (a AssetCode) String() string { return string(a) }

func (a AssetCode) Valid() bool {
	s := strings.TrimSpace(string(a))
	if len(s) < 2 || len(s) > 16 {
		return false
	}
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

func MustAsset(code string) AssetCode {
	a := AssetCode(strings.ToUpper(strings.TrimSpace(code)))
	if !a.Valid() {
		panic(errors.ErrInvalidAssetCode)
	}
	return a
}
