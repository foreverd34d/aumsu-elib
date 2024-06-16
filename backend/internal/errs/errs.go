// Пакет errs описывает ошибки, пригодные для переноса между слоями программы
// и для обработки в хэндлерах или middleware.
package errs

import "errors"

var (
	NotFound        = errors.New("not found")             // запрошенный ресурс не был найден
	Empty           = errors.New("the repo is empty")     // запрос в пустой репозиторий
	Internal        = errors.New("internal error")        // внутренняя ошибка репозитория или сервиса
	InvalidPassword = errors.New("invalid password")      // предоставленный пароль не совпадает с действительным
	InvalidLogin    = errors.New("invalid login")         // пользователь с таким именем не был найден
	RefreshExpired  = errors.New("refresh token expired") // токен для обновления истек
)
