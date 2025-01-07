package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
)

func (rt *_router) loginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Takes user r
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if !validIdentifier(user.Nickname) {
		ctx.Logger.WithError(err).Error("session: Can't Create a User. User nickname not Valid. <<")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if user exists
	temp_user, err := rt.db.CheckUser(user.toDataBase())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		ctx.Logger.WithError(err).Error("session: Error in CheckUser")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if temp_user.User_id != 0 { // ? Ho trovato l'user quindi assegno i valori
		user.User_id = temp_user.User_id
		user.Nickname = temp_user.Nickname

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(user)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("session: Can't create response json")
			return
		}
		return
	}

	// ! Controlla se l'utente è già presebnte nel DB
	err = rt.db.CreateUser(user.toDataBase())
	if err != nil {
		w.WriteHeader(http.StatusOK) // L'utente esisteva già nel DB
		_ = json.NewEncoder(w).Encode(user)
		return
	}

	// ! prendo l'ID Salvato nel DB
	id, err_fu := rt.db.FindUserId(user.toDataBase())
	if err_fu != nil {
		ctx.Logger.WithError(err_fu).Error("session: Error in FindUser_id()")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user.User_id = id // Cambio l'ID dentro la variabile user
}

/* ! Creo la cartella del nuovo Utente
CreateFolder(strconv.Itoa(id), ctx)

w.WriteHeader(http.StatusCreated)
err = json.NewEncoder(w).Encode(user)
if err != nil {
	w.WriteHeader(http.StatusInternalServerError)
	ctx.Logger.WithError(err).Error("session: can't create response json")
	return
}
*/
