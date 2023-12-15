package router

import (
	"github.com/claudiu-persoiu/septica/server/handler"
	"github.com/claudiu-persoiu/septica/server/templates"
	"net/http"
)

// Init generate new server router
func Init(th *templates.Handler, pubicHandler func(w http.ResponseWriter, r *http.Request)) http.Handler {

	router := http.NewServeMux()
	router.HandleFunc("/", th.ServeTemplate("template/index.html", nil))
	router.HandleFunc("/simulator", th.ServeTemplate("template/simulator.html", nil))
	router.HandleFunc("/public/", pubicHandler)
	router.HandleFunc("/ws", handler.HandleWebSocket)

	return router
}
