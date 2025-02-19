package database

import (
	"errors"
	"strings"
)

var ErrDuplicateReaction = errors.New("reaction already exists")

var ErrReactionNotFound = errors.New("reaction not found")

var ErrMessageNotFound = errors.New("message not found")

var ErrConversationNotFound = errors.New("conversation not found")

var (
	ErrUserAlreadyInGroup = errors.New("user already in group")
	ErrUserNotInGroup     = errors.New("user not in group")
	ErrNotAuthorized      = errors.New("not authorized to perform this action")
	ErrGroupNotFound      = errors.New("group not found")
)

// IsUniqueConstraintError controlla se l'errore è dovuto a un vincolo di unicità fallito
func IsUniqueConstraintError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}
