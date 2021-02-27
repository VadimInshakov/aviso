package rest

import (
	"aviso/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

type handlers struct {
	db *db.DB
}

func NewHandlers(db *db.DB) *handlers {
	return &handlers{db: db}
}

func (h *handlers) rootHandler(c *gin.Context){
	links, err := h.db.QueryAll()
	if err != nil {
		c.Data(500, "application/text; charset=utf-8", []byte("internal parsing error"))
		return
	}
	c.HTML(http.StatusOK, "root.tmpl", gin.H{"Links":links})
}