package api

import (
	"strings"
)

// verifica se l'id di un utente ha la lunghezza giusta
func validIdentifier(identifier string) bool {
	var nospace_id = strings.ReplaceAll(identifier, " ", "")
	// 5-25 come def in API
	return len(identifier) >= 3 && len(identifier) <= 16 && nospace_id != "" && !strings.ContainsAny(nospace_id, "?")
}
