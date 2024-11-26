package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	// Login
	rt.router.POST("/session", rt.wrap(rt.loginHandler))

	// User
	rt.router.PUT("/username/:userId", rt.wrap(rt.updateUserNameHandler))
	rt.router.POST("/sMyPhoto/:userId", rt.wrap(rt.updateUserPhotoHandler))

	// Conversation
	rt.router.GET("/conversations", rt.wrap(rt.getConversationsHandler))
	rt.router.GET("/conversation/:conversationId", rt.wrap(rt.getConversationByIDHandler))

	// Message
	rt.router.POST("/sMessage", rt.wrap(rt.sendMessageHandler))
	rt.router.POST("/fMessage/:messageId", rt.wrap(rt.forwardMessageHandler))
	rt.router.POST("/cMessage/:messageId", rt.wrap(rt.addReactionHandler))
	rt.router.DELETE("/dMessage/:messageId", rt.wrap(rt.deleteMessageHandler))

	// Group
	rt.router.POST("/aGroup/:groupId/users/:userId", rt.wrap(rt.addUserToGroupHandler))
	rt.router.PATCH("/sGroupName/:groupId", rt.wrap(rt.updateGroupNameHandler))
	rt.router.POST("/sGroupPhoto/:groupId", rt.wrap(rt.updateGroupPhotoHandler))

	return rt.router
}
