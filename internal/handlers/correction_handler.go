package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"seal-of-ash/internal/database"
	"seal-of-ash/internal/models"
)

type RequestCorrectionReq struct {
	TransactionID string `json:"transaction_id"`
	Reason        string `json:"reason"`
}

func RequestCorrection(c *gin.Context) {
	var req RequestCorrectionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid := c.GetUint("user_id")

	cr := models.CorrectionRequest{
		ID:            uuid.New().String(),
		TransactionID: req.TransactionID,
		Reason:        req.Reason,
		Status:        "pending",
		RequestedBy:   uid,
		CreatedAt:     time.Now(),
	}

	database.ActiveDB.Create(&cr)

	c.JSON(201, gin.H{
		"message": "Correction request has been submitted and awaits approval.",
		"id":      cr.ID,
	})
}
func ApproveCorrection(c *gin.Context) {
	id := c.Param("id")

	var cr models.CorrectionRequest
	if err := database.ActiveDB.First(&cr, "id = ?", id).Error; err != nil {
		c.JSON(404, gin.H{"error": "correction request not found"})
		return
	}

	if cr.Status != "pending" {
		c.JSON(400, gin.H{"error": "already processed"})
		return
	}

	var oldTx models.Transaction
	if err := database.ActiveDB.First(&oldTx, "id = ?", cr.TransactionID).Error; err != nil {
		c.JSON(404, gin.H{"error": "original transaction not found"})
		return
	}

	supervisorID := c.GetUint("user_id")

	// fallback kalau data lama belum punya root_id
	rootID := oldTx.RootID
	if rootID == "" {
		rootID = oldTx.ID
	}

	// reversal
	reversal := models.Transaction{
		ID:             uuid.New().String(),
		RootID:         rootID,
		Amount:         -oldTx.Amount,
		Type:           "reversal",
		RefTransaction: &oldTx.ID,
		CreatedBy:      supervisorID,
		CreatedAt:      time.Now(),
	}
	database.ActiveDB.Create(&reversal)

	// adjustment (sementara masih sama, nanti bisa diganti nilai baru)
	adjustment := models.Transaction{
		ID:             uuid.New().String(),
		RootID:         rootID,
		Amount:         oldTx.Amount,
		Type:           "adjustment",
		RefTransaction: &oldTx.ID,
		CreatedBy:      supervisorID,
		CreatedAt:      time.Now(),
	}
	database.ActiveDB.Create(&adjustment)

	cr.Status = "approved"
	database.ActiveDB.Save(&cr)

	c.JSON(200, gin.H{
		"message":       "Correction approved",
		"root_id":       rootID,
		"reversal_id":   reversal.ID,
		"adjustment_id": adjustment.ID,
	})
}
