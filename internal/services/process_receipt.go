package services

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"math"
	"github.com/google/uuid"
	"github.com/saurabhag23/receipt-processor/internal/database"
	"github.com/saurabhag23/receipt-processor/internal/models"
)

// ProcessReceipt validates and processes a receipt, calculating points, generating a unique ID,
// and storing the result in the database.
func ProcessReceipt(receipt *models.Receipt, userID int) (*models.ProcessedReceipt, error) {
    ctx := context.Background()
    tx, err := database.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
    if err != nil {
        return nil, fmt.Errorf("could not start transaction: %v", err)
    }

    // First, make sure the user exists or create/update points if necessary.
    if err := upsertUserPoints(tx, userID, calculatePoints(receipt)); err != nil {
        tx.Rollback()
        return nil, fmt.Errorf("could not upsert user points: %v", err)
    }

    receiptHash := generateReceiptHash(receipt,userID)
    id := uuid.New().String()

    // Check for duplicate receipts before inserting
    exists, err := checkReceiptExists(tx, receiptHash)
    if err != nil {
        tx.Rollback()
        return nil, fmt.Errorf("error checking receipt existence: %v", err)
    }
    if exists {
        tx.Rollback()
        return nil, fmt.Errorf("duplicate receipt submission")
    }

    // Proceed to insert the receipt
    if err := insertReceipt(tx, id, userID, calculatePoints(receipt), receiptHash); err != nil {
        tx.Rollback()
        return nil, fmt.Errorf("could not insert receipt: %v", err)
    }

    // Commit the transaction
    if err := tx.Commit(); err != nil {
        return nil, fmt.Errorf("could not commit transaction: %v", err)
    }

    return &models.ProcessedReceipt{ID: id, Points: calculatePoints(receipt), Hash: receiptHash}, nil
}

func generateReceiptHash(receipt *models.Receipt, userID int) string {
    hash := sha256.New()
    // Include the userID in the hash to ensure uniqueness per user
    hashInput := fmt.Sprintf("%d-%s-%s-%s-%s", 
        userID, 
        receipt.Retailer, 
        receipt.PurchaseDate, 
        receipt.PurchaseTime, 
        receipt.Total)

    hash.Write([]byte(hashInput))

    for _, item := range receipt.Items {
        hash.Write([]byte(fmt.Sprintf("%s-%s", item.ShortDescription, item.Price)))
    }
    return hex.EncodeToString(hash.Sum(nil))
}


func checkReceiptExists(tx *sql.Tx, hash string) (bool, error) {
	var exists bool
	err := tx.QueryRow("SELECT EXISTS (SELECT 1 FROM receipts WHERE receipt_hash = $1)", hash).Scan(&exists)
	return exists, err
}

func insertReceipt(tx *sql.Tx, receiptID string, userID int, points int, hash string) error {
	_, err := tx.Exec("INSERT INTO receipts (id, user_id, points, receipt_hash) VALUES ($1, $2, $3, $4)", receiptID, userID, points, hash)
	return err
}

func upsertUserPoints(tx *sql.Tx, userID int, pointsToAdd int) error {
    // This SQL statement attempts to insert a new user or updates the points if the user already exists.
    _, err := tx.Exec(`
        INSERT INTO users (user_id, points)
        VALUES ($1, $2)
        ON CONFLICT (user_id) 
        DO UPDATE SET points = users.points + EXCLUDED.points`,
        userID, pointsToAdd)
    if err != nil {
        return fmt.Errorf("could not upsert user points: %v", err)
    }
    return nil
}

func validateReceipt(r *models.Receipt) error {
	if r.Retailer == "" || r.PurchaseDate == "" || r.PurchaseTime == "" || len(r.Items) == 0 || r.Total == "" {
		return fmt.Errorf("missing required receipt fields")
	}

	if !regexp.MustCompile(`^[\w\s\-&]+$`).MatchString(r.Retailer) {
		return fmt.Errorf("invalid retailer name format")
	}

	if _, err := time.Parse("2006-01-02", r.PurchaseDate); err != nil {
		return fmt.Errorf("invalid purchase date format")
	}

	if _, err := time.Parse("15:04", r.PurchaseTime); err != nil {
		return fmt.Errorf("invalid purchase time format")
	}

	if !regexp.MustCompile(`^\d+\.\d{2}$`).MatchString(r.Total) {
		return fmt.Errorf("invalid total format")
	}

	return nil
}

func calculatePoints(r *models.Receipt) int {
	points := 0

	points += countAlphanumeric(r.Retailer)

	if strings.HasSuffix(r.Total, ".00") {
		points += 50
	}

	total, _ := strconv.ParseFloat(r.Total, 64)
	if int(total*100)%25 == 0 {
		points += 25
	}

	points += (len(r.Items) / 2) * 5

	for _, item := range r.Items {
		trimmedDescription := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDescription)%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(price * 0.2))
		}
	}

	if day, _ := time.Parse("2006-01-02", r.PurchaseDate); day.Day()%2 != 0 {
		points += 6
	}

	if purchaseTime, _ := time.Parse("15:04", r.PurchaseTime); purchaseTime.Hour() >= 14 && purchaseTime.Hour() < 16 {
		points += 10
	}

	return points
}

func countAlphanumeric(s string) int {
	count := 0
	for _, char := range s {
		if unicode.IsLetter(char) || unicode.IsNumber(char) {
			count++
		}
	}
	return count
}
