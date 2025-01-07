package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
)

func (rt *_router) updateNicknameHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Imposta l’header di risposta come JSON
	w.Header().Set("Content-Type", "application/json")

	// 1) Recupera l’ID utente dalla path
	userIDParam := ps.ByName("User_id") // corrisponde a :User_id nella rotta
	if userIDParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "missing user_id in path"})
		return
	}

	// 2) Decodifica il JSON dal body: in questo esempio, ipotizziamo che il client
	//    ci invii un oggetto con la nuova Nickname, ad es. { "userName": "NuovoNome" }
	var body struct {
		Nickname string `json:"Nickname"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid JSON payload"})
		return
	}

	// 3) Verifica la validità del nickname con la stessa funzione usata nel loginHandler
	if !validIdentifier(body.Nickname) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "invalid nickname format"})
		return
	}

	// 4) Creiamo un oggetto User (del package api) con l’ID recuperato dalla path
	//    e convertiamolo a int, gestendo eventuali errori
	//    Se stai già gestendo l'ID come string altrove, adatta di conseguenza
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "user_id must be an integer"})
		return
	}

	// 5) Struttura "api.User" con l’ID. Il Nickname non è essenziale qui, perché stiamo per aggiornarlo
	user := User{
		User_id:  userID,
		Nickname: "", // verrà cambiato
	}

	// 6) Chiedi al DB di verificare se l’utente esiste, se vuoi controllare l'esistenza
	dbUser, err := rt.db.CheckUser(user.toDataBase())
	// Se l’utente non esiste, CheckUser potrebbe restituire sql.ErrNoRows
	// oppure un errore custom. Adatta in base alla tua implementazione:
	if err != nil {
		// Se l'errore è “no rows”, restituisci 404
		if errors.Is(err, sql.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "user not found"})
			return
		}
		// Altrimenti errore generico
		ctx.Logger.WithError(err).Error("updateNickname: CheckUser error")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "internal DB error"})
		return
	}

	// 7) Ora abbiamo dbUser (User del package database). Proviamo a cambiare Nickname
	//    User -> toDataBase() e poi chiami un metodo di update (es. rt.db.ChangeUserName(...))
	err = rt.db.ChangeNickname(dbUser, body.Nickname)
	if err != nil {
		// Se c’è un errore di UNIQUE constraint per un Nickname già esistente
		// potresti restituire 409 (Conflict), altrimenti 500. Dipende dal design.
		ctx.Logger.WithError(err).Error("updateNickname: ChangeNickname error")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "cannot update nickname"})
		return
	}

	// 8) Rispondi con un 200 (o 204) e magari un messaggio
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Nickname updated successfully",
	})
}
