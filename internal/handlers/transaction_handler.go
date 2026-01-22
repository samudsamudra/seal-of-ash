package handlers

import (
	"net/http"

	"seal-of-ash/internal/database"
	"seal-of-ash/internal/events"
	"seal-of-ash/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateTransactionReq struct {
	Amount int64 `json:"amount"`
}

func CreateTransaction(c *gin.Context) {
	var req CreateTransactionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil user ID dari JWT (AuthMiddleware)
	uid := c.GetUint("user_id")

	txID := uuid.New().String()

	// Buat transaksi
	tx := models.Transaction{
		ID:        txID,
		RootID:    txID, // root = dirinya sendiri
		Amount:    req.Amount,
		Type:      "normal",
		CreatedBy: uid,
	}

	if err := database.ActiveDB.Create(&tx).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create transaction"})
		return
	}

	// Kirim event ke forensic worker
	events.EventBus <- events.Event{
		Type:     "CREATE",
		Entity:   "transaction",
		EntityID: tx.ID,
		ActorID:  uid,
		IP:       c.ClientIP(),
		UA:       c.Request.UserAgent(),
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":       tx.ID,
		"root_id":  tx.RootID,
		"actor_id": uid,
	})
}
