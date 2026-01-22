package handlers

import (
	// "net/http"

	"seal-of-ash/internal/database"
	"seal-of-ash/internal/models"
	"seal-of-ash/internal/utils"

	"github.com/gin-gonic/gin"
)

type SummonCaseReq struct {
	TransactionID string `json:"transaction_id"` // root_id
}

func SummonCase(c *gin.Context) {
	var req SummonCaseReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	rootID := req.TransactionID

	var txs []models.Transaction
	database.ActiveDB.
		Where("root_id = ?", rootID).
		Order("created_at asc").
		Find(&txs)

	if len(txs) == 0 {
		c.JSON(404, gin.H{"error": "case not found"})
		return
	}

	timeline := []gin.H{}
	finalAmount := int64(0)

	for _, tx := range txs {
		var user models.User
		database.ActiveDB.First(&user, "id = ?", tx.CreatedBy)

		timeline = append(timeline, gin.H{
			"type":   tx.Type,
			"amount": tx.Amount,
			"by":     user.Username,
			"at":     utils.FormatIndoTime(tx.CreatedAt),
		})

		finalAmount += tx.Amount
	}

	status := "normal"
	if len(txs) > 1 {
		status = "corrected"
	}

	c.JSON(200, gin.H{
		"case_id":      rootID,
		"timeline":     timeline,
		"final_amount": finalAmount,
		"status":       status,
	})
}
