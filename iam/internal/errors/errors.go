package errs

import "errors"

var (
	// Ошибки пользователей
	ErrUserNotFound       = errors.New("пользователь не найден")
	ErrUserAlreadyExists  = errors.New("пользователь с таким логином уже существует")
	ErrInvalidCredentials = errors.New("неверный логин или пароль")

	// Ошибки сессий
	ErrSessionNotFound = errors.New("сессия не найдена или истекла")

	// Ошибки валидации
	ErrInvalidLogin    = errors.New("логин обязателен")
	ErrInvalidPassword = errors.New("пароль обязателен")
	ErrWeakPassword    = errors.New("пароль должен содержать минимум 8 символов")
	ErrEmptyCredential = errors.New("логин и пароль обязательны")
	ErrEmptySessionID  = errors.New("session_uuid обязателен")
	ErrInvalidUUID     = errors.New("неверный формат UUID")

	ErrEmptyLogin    = errors.New("логин пустой")
	ErrEmptyPassword = errors.New("пароль пустой")

	ErrInvalidRegister = errors.New("неверный формат запроса для регистрации")
)
