package rest

import (
	"github.com/bskracic/sizif/runtime"
	"gorm.io/gorm"
)

type Handler struct {
	Db  *gorm.DB
	Rtm runtime.Runtime
}

func NewHandler(db *gorm.DB, rtm runtime.Runtime) *Handler {
	return &Handler{Db: db, Rtm: rtm}
}
