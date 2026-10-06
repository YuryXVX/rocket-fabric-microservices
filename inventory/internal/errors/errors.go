package errs

import "errors"

var (
	ErrPartNotFound           = errors.New("деталь не найдена")
	ErrInvalidUUID            = errors.New("неверный формат UUID")
	ErrNothingToRelease       = errors.New("нет деталей для резервирования")
	ErrOutOfStock             = errors.New("детали закончились")
	ErrInvalidProperties      = errors.New("невалидные значения корпуса")
	ErrIncompatibleHull       = errors.New("корпус не выдерживает нагрузку двигателя")
	ErrRequirementDetailCheck = errors.New("корпус или двигатель отсутсвуют")
	ErrIncompatibleShield     = errors.New("щит нельзя использовать с этим типом оружия")
	ErrDuplicatePartUUID      = errors.New("обнаружены дубликаты идентификаторов деталей в запросе")
)
