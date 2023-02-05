package rest

import (
	"aviso/db/sqlite"
	"github.com/gin-gonic/gin"
	"net/http"
)

type handlers struct {
	db *sqlite.DB
}

func NewHandlers(db *sqlite.DB) *handlers {
	return &handlers{db: db}
}

func (h *handlers) rootHandler(c *gin.Context) {
	links, err := h.db.QueryAll()
	if err != nil {
		c.Data(500, "application/text; charset=utf-8", []byte("internal parsing error"))
		return
	}
	c.HTML(http.StatusOK, "root.tmpl", gin.H{"Links": links})
}
