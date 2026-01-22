package handlers

import (
	"crypto/sha256"
	"encoding/hex"
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
		"prev_hash": ash.PrevHash,
		"created":   utils.FormatIndoTime(ash.CreatedAt),
	})
}

func VerifyAshChain(c *gin.Context) {
	var ashes []models.ForensicAsh
	database.ForensicDB.Order("created_at asc").Find(&ashes)

	prev := "GENESIS"

	for i, ash := range ashes {
		payload := append([]byte(prev), ash.Snapshot...)
		payload = append(payload, []byte(ash.Action)...)
		payload = append(payload, []byte(ash.EntityType)...)
		payload = append(payload, []byte(ash.EntityID)...)

		sum := sha256.Sum256(payload)
		expected := hex.EncodeToString(sum[:])

		if ash.Hash != expected {
			c.JSON(500, gin.H{
				"status": "corrupted",
				"at":     i,
				"id":     ash.ID,
			})
			return
		}

		prev = ash.Hash
	}

	c.JSON(200, gin.H{
		"status": "intact",
		"total":  len(ashes),
	})
}
