package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

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
	database.ForensicDB.Order("created_at asc, id asc").Find(&ashes)

	// kalau belum ada abu sama sekali
	if len(ashes) == 0 {
		c.JSON(200, gin.H{
			"status":  "empty",
			"message": "The ash vault is silent. No records have been forged yet.",
		})
		return
	}

	prev := "GENESIS"

	for i, ash := range ashes {
		// bangun ulang hash dari data yang ada
		payload := append([]byte(prev), ash.Snapshot...)
		payload = append(payload, []byte(ash.Action)...)
		payload = append(payload, []byte(ash.EntityType)...)
		payload = append(payload, []byte(ash.EntityID)...)

		sum := sha256.Sum256(payload)
		expected := hex.EncodeToString(sum[:])

		if ash.Hash != expected {
			// ambil data user pelaku
			var user models.User
			database.ActiveDB.First(&user, "id = ?", ash.ActorID)

			detectedAt := utils.FormatIndoTime(time.Now())

			description := fmt.Sprintf(
				"Data terkorup. Terdapat ketidaksesuaian pencatatan oleh user '%s' dengan ID '%d'. "+
					"Ketidaksamaan hash terdeteksi pada %s WIB.",
				user.Username,
				user.ID,
				detectedAt,
			)

			c.JSON(409, gin.H{
				"status":      "corrupted",
				"description": description,
				"forensic_detail": gin.H{
					"user_id":     user.ID,
					"username":    user.Username,
					"record_id":   ash.ID,
					"detected_at": detectedAt,
					"chain_index": i,
					"prev_hash":   ash.PrevHash,
					"hash":        ash.Hash,
				},
			})
			return
		}

		prev = ash.Hash
	}

	c.JSON(200, gin.H{
		"status": "intact",
		"total":  len(ashes),
		"message": fmt.Sprintf(
			"Seluruh data forensik valid. %d catatan abu terikat sempurna dalam Seal of Ash.",
			len(ashes),
		),
	})
}
