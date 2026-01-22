package handlers

import (
	"net/http"
	"seal-of-ash/internal/database"
	"seal-of-ash/internal/events"
	"seal-of-ash/internal/models"
	"github.com/gin-gonic/gin"
)

type CreateTransactionReq struct {
	Amount int64 `json:"amount"`
	UserID uint  `json:"user_id"` // nanti dari JWT
}

func CreateTransaction(c *gin.Context) {
	var req CreateTransactionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := models.Transaction{
		Amount:    req.Amount,
		Type:      "normal",
		CreatedBy: req.UserID,
	}

	database.ActiveDB.Create(&tx)

	// kirim event ke forensic worker
	events.EventBus <- events.Event{
		Type:      "CREATE",
		Entity:    "transaction",
		EntityID:  tx.ID,
		ActorID:  req.UserID,
		IP:       c.ClientIP(),
		UA:       c.Request.UserAgent(),
	}

	c.JSON(201, gin.H{
		"message": "Transaction forged into the ledger",
		"id":      tx.ID,
	})
}
