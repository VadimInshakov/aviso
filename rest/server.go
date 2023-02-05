package rest

import (
	"aviso/db/sqlite"
	"github.com/gin-gonic/gin"
	"net"
)

func Run(db *sqlite.DB, host, port string) {
	r := gin.Default()
	handlers := NewHandlers(db)
	r.LoadHTMLGlob("rest/templates/*")
	r.GET("/", handlers.rootHandler)
	r.Run(net.JoinHostPort(host, port))
}
