package controller

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/mbarolo/webp_backend/services"
	"github.com/mbarolo/webp_backend/services/sessions"
	"github.com/mbarolo/webp_backend/utils"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return r.Header.Get("Origin") == os.Getenv("FRONTEND_URL")
	},
}

func CreateGameHandler(c *gin.Context) {
	game, ok := services.GameRegistry[c.Param("game")]
	if !ok {
		utils.JsonResponse(c, http.StatusNotFound, "Game not found", nil)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return //error already returned by upgrader
	}
	defer conn.Close()

	session := sessions.Registry.Create(game())
	player, err := session.Join(conn)
	if err != nil {
		return
	}

	runPlayer(session, player, conn)
}

func JoinGameHandler(c *gin.Context) {
	session, ok := sessions.Registry.Find(c.Param("code"))
	if !ok {
		utils.JsonResponse(c, http.StatusNotFound, "Session not found", nil)
		return
	}
	if err := session.CanJoin(); err != nil {
		utils.JsonResponse(c, http.StatusConflict, err.Error(), nil)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.JsonResponse(c, http.StatusInternalServerError, "error upgrading request "+err.Error(), nil)
		return
	}

	player, err := session.Join(conn)
	if err != nil {
		conn.WriteJSON(gin.H{"type": "error", "message": err.Error()})
		return
	}

	runPlayer(session, player, conn)
}

func runPlayer(s *sessions.Session, p *sessions.Player, conn *websocket.Conn) {
	conn.SetReadDeadline(time.Now().Add(sessions.PONG_WAIT))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(sessions.PONG_WAIT))
	})

	done := make(chan struct{})
	defer close(done)
	go func() {
		ticker := time.NewTicker(sessions.PING_PERIOD)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if s.Ping(p, conn) != nil {
					return
				}
			case <-done:
				return
			}
		}
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Disconnect(p, conn)
			return
		}
		s.HandleMessage(p, msg)
	}
}
