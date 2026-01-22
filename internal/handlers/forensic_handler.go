package handlers

import (
	"net/http"
	"seal-of-ash/internal/database"
	"seal-of-ash/internal/models"
	"seal-of-ash/internal/utils"

	"github.com/gin-gonic/gin"
)

type SummonAshReq struct {
	Entity   string `json:"entity"`    // contoh: "transaction"
	EntityID string `json:"entity_id"` // UUID transaksi
}

func SummonAshes(c *gin.Context) {
	var req SummonAshReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var ash models.ForensicAsh
	err := database.ForensicDB.
		Where("entity_type = ? AND entity_id = ?", req.Entity, req.EntityID).
		Order("created_at desc").
		First(&ash).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "the ashes remain silent",
		})
		return
	}

	c.JSON(200, gin.H{
		"entity":    ash.EntityType,
		"entity_id": ash.EntityID,
		"action":    ash.Action,
		"snapshot":  string(ash.Snapshot),
		"actor_id":  ash.ActorID,
		"ip":        ash.IP,
		"ua":        ash.UserAgent,
		"hash":      ash.Hash,
		"created":   utils.FormatIndoTime(ash.CreatedAt),
	})
}
