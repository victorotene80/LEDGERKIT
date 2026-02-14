package errors

import "errors"

var (
	ErrInvalidAssetCode = errors.New("invalid asset code")
	ErrInvalidScale     = errors.New("invalid scale")
	ErrInvalidAmount    = errors.New("invalid amount")
	ErrAssetMismatch    = errors.New("asset mismatch")
	ErrScaleMismatch    = errors.New("scale mismatch")
	ErrOverflow         = errors.New("integer overflow")

	ErrInvalidAccountRef     = errors.New("invalid account ref")
	ErrInvalidSide           = errors.New("invalid side")
	ErrInvalidPosting        = errors.New("invalid posting")
	ErrEntryNeedsTwoPostings = errors.New("entry must have at least two postings")
	ErrZeroPosting           = errors.New("zero-value posting not allowed")
	ErrPostingNotPositive    = errors.New("posting amount must be positive")
	ErrUnbalancedEntry       = errors.New("unbalanced journal entry")
)
