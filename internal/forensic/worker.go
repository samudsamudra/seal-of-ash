package forensic

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"seal-of-ash/internal/database"
	"seal-of-ash/internal/events"
	"seal-of-ash/internal/models"
	"time"

	"github.com/google/uuid"
)

func StartWorker() {
	go func() {
		for e := range events.EventBus {
			processEvent(e)
		}
	}()
}

func processEvent(e events.Event) {
	var snapshot any

	switch e.Entity {
	case "transaction":
		var tx models.Transaction
		database.ActiveDB.First(&tx, e.EntityID)
		snapshot = tx
	default:
		return
	}

	raw, _ := json.Marshal(snapshot)

	hash := sha256.Sum256(raw)

	record := models.ForensicAsh{
		ID:         uuid.New().String(),
		EntityType: e.Entity,
		EntityID:   e.EntityID,
		Action:     e.Type,
		Snapshot:   raw,
		ActorID:    e.ActorID,
		IP:         e.IP,
		UserAgent:  e.UA,
		Hash:       hex.EncodeToString(hash[:]),
		CreatedAt:  time.Now(),
	}

	database.ForensicDB.Create(&record)
}
