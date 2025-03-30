package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/tommox/WASAText/service/api/reqcontext"
)

func (rt *_router) setMyNicknameHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	user_id, err := strconv.Atoi(ps.ByName("User_id"))
	if err != nil {
		ctx.Logger.WithError(err).Error("put-nickname: error converting id")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// ! Login
	requestingUserId_str, err := extractBearerToken(r, w)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		ctx.Logger.WithError(err).Error("put-nickname: Error")
		return
	}
	requestingUserId, err := strconv.Atoi(requestingUserId_str)
	if err != nil {
		ctx.Logger.WithError(err).Error("put-nickname: error converting requesingUserId")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if user_id != requestingUserId {
		w.WriteHeader(http.StatusUnauthorized)
		ctx.Logger.WithError(errors.New("you aren't allowed to use this operation")).Error("put-Nickname: Error")
		return
	}

	// Takes the body
	var newNickname Nickname
	err = json.NewDecoder(r.Body).Decode(&newNickname)
	if err != nil {
		ctx.Logger.WithError(err).Error("put-nickname: error decoding json")
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if !validIdentifier(newNickname.Nickname) {
		ctx.Logger.WithError(errors.New("nickname provided does not meet the necessary conditions")).Error("put-nickname: Error in Nickname section")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Check if nickname is free
	exist, err := rt.db.SearchNickname(newNickname.Nickname)
	if err != nil {
		ctx.Logger.WithError(err).Error("put-nickname: Error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if exist {
		ctx.Logger.WithError(errors.New("nickname already exists")).Error("put-nickname: Error in Nickname section")
		w.WriteHeader(http.StatusConflict)
		return
	}
	err = rt.db.ChangeNickname(User{User_id: user_id}.toDataBase(), newNickname.Nickname)
	if err != nil {
		ctx.Logger.WithError(err).Error("put-nickname: Error executing query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Nickname successfully updated",
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
