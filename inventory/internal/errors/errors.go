package errs

import "errors"

var (
	ErrPartNotFound      = errors.New("деталь не найдена")
	ErrInvalidUUID       = errors.New("неверный формат UUID")
	ErrNothingToRelease  = errors.New("нет деталей для резервирования")
	ErrOutOfStock        = errors.New("детали закончились")
	ErrInvalidProperties = errors.New("невалидные значения корпуса")
)
