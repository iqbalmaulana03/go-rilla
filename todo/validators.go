package todo

import "errors"

var (
	ErrTitleTooShort = errors.New("todo: title to short")
	ErrTitleTooLong  = errors.New("todo: title to long")
	ErrTitleEmpty    = errors.New("todo: title empty")
)

const maxTitle = 1000
const minTitle = 5

func validateTitle(title string) error {
	l := len(title)

	switch {
	case l == 0:
		return ErrTitleEmpty
	case l < minTitle:
		return ErrTitleTooShort
	case l > maxTitle:
		return ErrTitleTooLong
	default:
		return nil
	}
}
