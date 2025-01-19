package api

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
)

// verifica se l'id di un utente ha la lunghezza giusta
func validIdentifier(identifier string) bool {
	var nospace_id = strings.ReplaceAll(identifier, " ", "")
	// 3-16 come def in API
	return len(identifier) >= 3 && len(identifier) <= 16 && nospace_id != "" && !strings.ContainsAny(nospace_id, "?")
}

// verifica se il messaggio di un utente è accettabile
var msgRegex = regexp.MustCompile(`^[a-zA-Z0-9 .,!?'\-]+$`)

func validMessage(message string) bool {
	// Controllo lunghezza
	if len(message) < 1 || len(message) > 1000 {
		return false
	}
	// Verifica il pattern
	return msgRegex.MatchString(message)
}

func extractBearerToken(req *http.Request, w http.ResponseWriter) (string, error) {

	// Ottieni l'intestazione Authorization dalla richiesta
	authHeader := req.Header.Get("Authorization")
	// Verifica se l'intestazione Authorization è vuota o mancante
	if authHeader == "" {
		w.WriteHeader(http.StatusForbidden)
		return "", errors.New("intestazione Authorization mancante")
	}

	list := strings.Split(authHeader, " ")
	token := list[len(list)-1]

	return token, nil
}
