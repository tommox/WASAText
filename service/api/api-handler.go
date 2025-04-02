package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// LOGIN
	rt.router.POST("/session", rt.wrap(rt.doLoginHandler))

	// USERS
	rt.router.PUT("/users/:User_id", rt.wrap(rt.setMyNicknameHandler))
	rt.router.PUT("/users/:User_id/photo", rt.wrap(rt.setMyPhotoHandler))
	rt.router.GET("/users", rt.wrap(rt.getUsersHandler))
	rt.router.GET("/users/:User_id/photo", rt.wrap(rt.getUserPhotoHandler))

	// MESSAGES
	rt.router.POST("/messages", rt.wrap(rt.sendMessageHandler))
	rt.router.GET("/messages/:Message_id", rt.wrap(rt.getMessageHandler))
	rt.router.DELETE("/messages/:Message_id", rt.wrap(rt.deleteMessageHandler))
	rt.router.POST("/messages/:Message_id/forwards", rt.wrap(rt.forwardMessageHandler))
	rt.router.POST("/messages/:Message_id/reactions", rt.wrap(rt.addReactionHandler))
	rt.router.DELETE("/messages/:Message_id/reactions", rt.wrap(rt.removeReactionHandler))
	rt.router.GET("/messages/:Message_id/reactions", rt.wrap(rt.getReactionHandler))

	// GROUPS
	rt.router.GET("/groups/:Group_id/users", rt.wrap(rt.getUsersOfGroupHandler))
	rt.router.POST("/groups", rt.wrap(rt.createGroupHandler))
	rt.router.DELETE("/groups/:Group_id", rt.wrap(rt.deleteGroupHandler))
	rt.router.POST("/groups/:Group_id/users/:User_id", rt.wrap(rt.addToGroupHandler))
	rt.router.DELETE("/groups/:Group_id/users/:User_id", rt.wrap(rt.leaveGroupHandler))
	rt.router.PATCH("/groups/:Group_id", rt.wrap(rt.changeGroupNameHandler))
	rt.router.PUT("/groups/:Group_id/photo", rt.wrap(rt.setGroupPhotoHandler))
	rt.router.GET("/groups/:Group_id/photo", rt.wrap(rt.getGroupPhotoHandler))
	rt.router.POST("/groups/:Group_id/messages", rt.wrap(rt.sendMessageToGroupHandler))

	// CONVERSATIONS
	rt.router.GET("/conversations", rt.wrap(rt.getMyConversationsHandler))
	rt.router.GET("/conversations/:Conversation_id", rt.wrap(rt.getConversationHandler))
	rt.router.POST("/conversations/conversation", rt.wrap(rt.checkOrCreateConversationHandler))
	rt.router.DELETE("/conversations/:Conversation_id", rt.wrap(rt.deleteConversationHandler))
	return rt.router
}
