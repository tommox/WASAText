package api

import (
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
)

func (rt *_router) setMyPhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	user_id, err := strconv.Atoi(ps.ByName("User_id"))
	if err != nil {
		// Errore: L'ID utente fornito non è valido
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Error: invalid User_id")
		return
	}

	// Estrai lo User_id dal token Bearer
	requestingUser_id_str, err := extractBearerToken(r, w)
	if err != nil {
		// Errore: Token Bearer mancante o non valido
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("Error: unauthorized access, missing or invalid token")
		return
	}

	requestingUser_id, err := strconv.Atoi(requestingUser_id_str)
	if err != nil {
		// Errore: User_id non valido
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("Error: invalid User_id from token")
		return
	}

	if user_id != requestingUser_id {
		// Errore: L'utente corrente non è autorizzato a modificare la foto per questo ID
		w.WriteHeader(http.StatusUnauthorized)
		ctx.Logger.WithError(err).Error("setMyPhoto: unauthorized, user_id does not match token user")
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		// Errore: Problema nella lettura del body
		w.WriteHeader(http.StatusBadRequest)
		ctx.Logger.WithError(err).Error("setMyPhoto: error reading photo data from request body")
		return
	}

	err = rt.db.UpdateUserPhoto(user_id, data)
	if err != nil {
		// Errore: Problema durante l'aggiornamento della foto nel database
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("setMyPhoto: error updating photo in database")
		return
	}

	w.WriteHeader(http.StatusOK)
	ctx.Logger.Infof("Photo updated successfully for user_id: %d", user_id)
}
