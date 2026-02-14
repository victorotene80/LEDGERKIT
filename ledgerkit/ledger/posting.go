package ledger

import (
	lkerr "github.com/Helen-projects/LEDGERKIT/ledgerkit/errors"
	"github.com/Helen-projects/LEDGERKIT/ledgerkit/money"
)

type Posting struct {
	account AccountRef
	side    Side
	amt     money.Money
}

func NewPosting(account AccountRef, side Side, amt money.Money) (Posting, error) {
	// account validation (minimal; full rules live in validate.go too)
	if account.ID() == "" || !account.Kind().ValidBasic() {
		return Posting{}, lkerr.ErrInvalidAccountRef
	}

	if err := side.Validate(); err != nil {
		return Posting{}, err
	}

	if amt.IsZero() {
		return Posting{}, lkerr.ErrZeroPosting
	}

	if !amt.IsPositive() {
		return Posting{}, lkerr.ErrPostingNotPositive
	}

	return Posting{account: account, side: side, amt: amt}, nil
}

func MustPosting(account AccountRef, side Side, amt money.Money) Posting {
	p, err := NewPosting(account, side, amt)
	if err != nil {
		panic(err)
	}
	return p
}

func (p Posting) Account() AccountRef { return p.account }
func (p Posting) Side() Side          { return p.side }
func (p Posting) Money() money.Money  { return p.amt }
