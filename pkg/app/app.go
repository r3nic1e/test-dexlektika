package app

import (
	"net/http"

	"gorm.io/gorm"
)

type App struct {
	mux *http.ServeMux
	db  *gorm.DB
}

func NewApp() *App {
	return &App{
		mux: http.NewServeMux(),
	}
}
