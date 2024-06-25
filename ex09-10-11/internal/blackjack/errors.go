package blackjack

var (
	errBust  error
	errSplit error
)

// type errType uint8
//
// const (
// 	bust errType = iota
// 	split
// )
//
// func NewError(typ errType, err ...error) error {
// 	switch typ {
// 	case bust:
// 		return errBust
// 	case split:
// 		return errSplit
// 	default:
// 		return errors.New("unexpected error")
// 	}
// }
//
// type ErrBust struct {
// 	Msg string
// 	Err error
// }
//
// func (e *ErrBust) Error() string {
// 	return fmt.Sprintf("%s: %v", e.Msg, e.Err)
// }
//
// type ErrSplit struct {
// 	Reason string
// }
//
// func (e *ErrSplit) Error() string {
// 	return fmt.Sprintf("split: %s", e.Reason)
// }
