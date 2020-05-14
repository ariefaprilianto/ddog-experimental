package err

import "errors"

// error list
var (
	ErrKeyExists    = errors.New("Key already exists")
	ErrKeyNotFound  = errors.New("Key not found")
	ErrKeyDuplicate = errors.New("Key already exists")
)
