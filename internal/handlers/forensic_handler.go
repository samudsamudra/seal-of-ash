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
	// ambil semua forensic ash
	var ashes []models.ForensicAsh
	database.ForensicDB.
		Order("created_at asc, id asc").
		Find(&ashes)

	// hitung transaksi di ledger
	var txCount int64
	database.ActiveDB.
		Model(&models.Transaction{}).
		Count(&txCount)

	// kalau vault kosong
	if len(ashes) == 0 {
		c.JSON(200, gin.H{
			"status":              "empty",
			"message":             "Forensic vault belum memiliki catatan. Sistem forensik belum mulai merekam.",
			"forensic_records":    0,
			"ledger_transactions": txCount,
			"note":                "Jumlah transaksi lebih banyak dari forensic record berarti ada transaksi yang dibuat sebelum sistem forensik aktif.",
		})
		return
	}

	prev := "GENESIS"

	for i, ash := range ashes {
		payload := append([]byte(prev), ash.Snapshot...)
		payload = append(payload, []byte(ash.Action)...)
		payload = append(payload, []byte(ash.EntityType)...)
		payload = append(payload, []byte(ash.EntityID)...)

		sum := sha256.Sum256(payload)
		expected := hex.EncodeToString(sum[:])

		if ash.Hash != expected {
			c.JSON(409, gin.H{
				"status":              "corrupted",
				"message":             "Integritas forensic vault rusak. Ada data yang dimodifikasi secara tidak sah.",
				"forensic_records":    len(ashes),
				"ledger_transactions": txCount,
				"broken_at_index":     i,
				"broken_record_id":    ash.ID,
			})
			return
		}

		prev = ash.Hash
	}

	// vault aman
	response := gin.H{
		"status":              "intact",
		"message":             "Seluruh forensic record valid dan terikat oleh hash chain.",
		"forensic_records":    len(ashes),
		"ledger_transactions": txCount,
	}

	// kasih warning kalau tidak sinkron
	if int64(len(ashes)) != txCount {
		response["warning"] = "Jumlah forensic record tidak sama dengan jumlah transaksi ledger. Tidak semua transaksi tercatat dalam sistem forensik."
	}

	c.JSON(200, response)
}
