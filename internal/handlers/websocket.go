package handlers

import (
	"github.com/fasthttp/websocket"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/isaee-xyz/whatomate/internal/middleware"
	ws "github.com/isaee-xyz/whatomate/internal/websocket"
	"github.com/valyala/fasthttp"
	"github.com/zerodha/fastglue"
)

var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
		return true // Allow all origins in development
	},
}

// WebSocketHandler handles WebSocket connections
func (a *App) WebSocketHandler(r *fastglue.Request) error {
	// Get token from query parameter
	token := string(r.RequestCtx.QueryArgs().Peek("token"))
	if token == "" {
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Missing token", nil, "")
	}

	// Validate JWT token
	userID, orgID, err := a.validateWSToken(token)
	if err != nil {
		a.Log.Error("WebSocket auth failed", "error", err)
		return r.SendErrorEnvelope(fasthttp.StatusUnauthorized, "Invalid token", nil, "")
	}

	// Upgrade to WebSocket
	err = upgrader.Upgrade(r.RequestCtx, func(conn *websocket.Conn) {
		client := ws.NewClient(a.WSHub, conn, userID, orgID)

		// Register client with hub
		a.WSHub.Register(client)

		// Start pumps in goroutines
		go client.WritePump()
		client.ReadPump() // Blocking - runs until connection closes
	})

	if err != nil {
		a.Log.Error("WebSocket upgrade failed", "error", err)
		return r.SendErrorEnvelope(fasthttp.StatusInternalServerError, "WebSocket upgrade failed", nil, "")
	}

	return nil
}

// validateWSToken validates a JWT token and returns user ID and organization ID
func (a *App) validateWSToken(tokenString string) (uuid.UUID, uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &middleware.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.Config.JWT.Secret), nil
	})

	if err != nil || !token.Valid {
		return uuid.Nil, uuid.Nil, err
	}

	claims, ok := token.Claims.(*middleware.JWTClaims)
	if !ok {
		return uuid.Nil, uuid.Nil, jwt.ErrTokenInvalidClaims
	}

	return claims.UserID, claims.OrganizationID, nil
}
