package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
)

func (rt *_router) doLoginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	w.Header().Set("Content-Type", "application/json")

	// Takes user r
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if !validIdentifier(user.Nickname) {
		ctx.Logger.WithError(err).Error("session: Can't Create a User. Invalid nickname lenght")
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
	if temp_user.User_id != 0 { // assignment of user values
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

	// Check if user is already in DB
	err = rt.db.CreateUser(user.toDataBase())
	if err != nil {
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(user)
		return
	}

	// Takes the ID saved in DB
	id, err_fu := rt.db.FindUserId(user.toDataBase())
	if err_fu != nil {
		ctx.Logger.WithError(err_fu).Error("session: Error in FindUser_id()")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user.User_id = id // Change the ID inside of the user variable

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("session: can't create response json")
		return
	}
}
