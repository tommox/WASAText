package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {

	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	rt.router.GET("/liveness", rt.liveness)

	// LOGIN
	rt.router.POST("/session", rt.wrap(rt.doLoginHandler))

	// USERS
	rt.router.PUT("/users/:User_id", rt.wrap(rt.setMyNicknameHandler))
	rt.router.PUT("/users/:User_id/photo", rt.wrap(rt.setMyPhotoHandler))

	// MESSAGES
	rt.router.POST("/messages", rt.wrap(rt.sendMessageHandler))
	rt.router.GET("/messages/:Message_id/details", rt.wrap(rt.getMessageHandler))
	rt.router.POST("/messages/:Message_id/forwards", rt.wrap(rt.forwardMessageHandler))
	rt.router.POST("/messages/:Message_id/reactions", rt.wrap(rt.manageReactionHandler))
	/*
		rt.router.DELETE("/messages/:Message_id", rt.wrap(rt.deleteMessageHandler))

		// GROUPS
		rt.router.POST("/groups/:Group_id/users", rt.wrap(rt.addUserToGroupHandler))
		rt.router.DELETE("/groups/:Group_id/users/:User_id", rt.wrap(rt.removeUserFromGroupHandler))
		rt.router.PATCH("/groups/:Group_id", rt.wrap(rt.updateGroupNameHandler))
		rt.router.PUT("/groups/:Group_id/photo", rt.wrap(rt.updateGroupPhotoHandler))

		// CONVERSATIONS
		rt.router.GET("/conversations", rt.wrap(rt.getConversationsHandler))
		rt.router.GET("/conversations/:Conversation_id", rt.wrap(rt.getConversationByIDHandler))
	*/
	return rt.router
}
