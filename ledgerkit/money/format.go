package money

import (
	"fmt"
)

// String implements fmt.Stringer.
// Output example: "NGN 100.00" for scale=2 minor=10000
func (m Money) String() string {
	sign := ""
	minor := m.amountMinor

	if minor < 0 {
		sign = "-"
		// careful with MinInt64 (you already handle overflow elsewhere, but do it safely here)
		if minor == -minor {
			// MinInt64 edge case: just print raw minor
			return fmt.Sprintf("%s%s %d", sign, m.asset.String(), m.amountMinor)
		}
		minor = -minor
	}

	scale := int(m.scale)
	if scale == 0 {
		return fmt.Sprintf("%s%s %d", sign, m.asset.String(), minor)
	}

	pow := int64(1)
	for i := 0; i < scale; i++ {
		pow *= 10
	}

	whole := minor / pow
	frac := minor % pow

	return fmt.Sprintf("%s%s %d.%0*d", sign, m.asset.String(), whole, scale, frac)
}
