package services

import (
	"database/sql"
	"fmt"

	"github.com/saurabhag23/receipt-processor/internal/database"
	"github.com/saurabhag23/receipt-processor/internal/models"
)

// GetPointsByID retrieves the points awarded to a receipt based on its unique ID.
func GetPointsByID(id string) (*models.ProcessedReceipt, error) {
	processedReceipt := &models.ProcessedReceipt{ID: id}

	// Prepare SQL statement for retrieving points
	stmt, err := database.DB.Prepare("SELECT points FROM receipts WHERE id = $1")
	if err != nil {
		return nil, fmt.Errorf("error preparing query: %v", err)
	}
	defer stmt.Close()

	// Execute the query
	err = stmt.QueryRow(id).Scan(&processedReceipt.Points)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no receipt found for ID: %s", id)
		}
		return nil, fmt.Errorf("error querying points for receipt: %v", err)
	}

	return processedReceipt, nil
}
