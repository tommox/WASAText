package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Rotte "demo"
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	rt.router.GET("/liveness", rt.liveness)

	// LOGIN
	rt.router.POST("/session", rt.wrap(rt.doLoginHandler))

	// USERS
	rt.router.PUT("/users/:User_id", rt.wrap(rt.setMyUserNameHandler))
	rt.router.PUT("/photos/:User_id", rt.wrap(rt.setMyPhotoHandler))

	// CONVERSATIONS
	rt.router.GET("/conversations", rt.wrap(rt.getConversationsHandler))
	rt.router.GET("/conversations/:Conversation_id", rt.wrap(rt.getConversationByIDHandler))

	// MESSAGES
	rt.router.POST("/messages", rt.wrap(rt.sendMessageHandler))
	rt.router.POST("/messages/:Message_id/forwards", rt.wrap(rt.forwardMessageHandler))
	rt.router.POST("/messages/:Message_id/reactions", rt.wrap(rt.addReactionHandler))
	rt.router.DELETE("/messages/:Message_id/reactions", rt.wrap(rt.removeReactionHandler))
	rt.router.DELETE("/messages/:Message_id", rt.wrap(rt.deleteMessageHandler))

	// GROUPS
	rt.router.POST("/groups/:Group_id/users", rt.wrap(rt.addUserToGroupHandler))
	rt.router.DELETE("/groups/:Group_id/users/:User_id", rt.wrap(rt.removeUserFromGroupHandler))
	rt.router.PATCH("/groups/:Group_id", rt.wrap(rt.updateGroupNameHandler))
	rt.router.PUT("/groups/:Group_id/photo", rt.wrap(rt.updateGroupPhotoHandler))

	// PERMISSIONS
	rt.router.POST("/permissions/users/:User_id", rt.wrap(rt.checkUserPermissionsHandler))

	return rt.router
}
