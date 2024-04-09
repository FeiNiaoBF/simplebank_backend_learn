package api

import (
	db "github.com/FeiNiaoBF/simplebank_backend_learn/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

// NewServer create a new Server instance
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	server.router = router
	return server
}

// Start start the HTTP server
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// errorResponse
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
