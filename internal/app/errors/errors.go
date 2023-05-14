package errors

import "errors"

// ErrURLIsDeleted ошибка при попытке получения удаленного URL.
var ErrURLIsDeleted = errors.New("url is deleted")
