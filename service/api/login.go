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

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if !validIdentifier(user.Nickname) {
		ctx.Logger.WithError(err).Error("session: Can't Create a User. Invalid nickname length")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	temp_user, err := rt.db.CheckUser(user.toDataBase())
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		ctx.Logger.WithError(err).Error("session: Error in CheckUser")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err == nil {
		user.User_id = temp_user.User_id
		user.Nickname = temp_user.Nickname

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			ctx.Logger.WithError(err).Error("session: Can't create response json")
		}
		return
	}

	err = rt.db.CreateUser(user.toDataBase())
	if err != nil {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(user); err != nil {
			ctx.Logger.WithError(err).Error("session: Can't create response json (new user)")
		}
		return
	}

	id, err_fu := rt.db.FindUserId(user.toDataBase())
	if err_fu != nil {
		ctx.Logger.WithError(err_fu).Error("session: Error in FindUser_id()")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user.User_id = id

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ctx.Logger.WithError(err).Error("session: can't create response json")
	}
}
