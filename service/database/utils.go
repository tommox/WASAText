package database

import (
	"errors"
	"strings"
)

// ErrDuplicateReaction indica che la reazione esiste già
var ErrDuplicateReaction = errors.New("reaction already exists")

// ErrReactionNotFound indica che la reazione non è stata trovata
var ErrReactionNotFound = errors.New("reaction not found")

// ErrMessageNotFound indica che il messaggio non è stato trovato
var ErrMessageNotFound = errors.New("message not found")

// IsUniqueConstraintError controlla se l'errore è dovuto a un vincolo di unicità fallito
func IsUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}
