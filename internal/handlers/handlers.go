// handlers.go
package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/saurabhag23/receipt-processor/internal/models"
	"github.com/saurabhag23/receipt-processor/internal/services"
	"github.com/saurabhag23/receipt-processor/internal/utils"
)

// ProcessReceipt handles the POST request to process a receipt.
// It validates the JWT, decodes the receipt, processes it, and returns the receipt ID.
func ProcessReceipt(w http.ResponseWriter, r *http.Request) {
    var data struct {
        UserID int             `json:"userId"`
        Receipt models.Receipt `json:"receipt"`
    }
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, "Invalid JSON format", http.StatusBadRequest)
        return
    }

    processedReceipt, err := services.ProcessReceipt(&data.Receipt, data.UserID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    response := map[string]string{"id": processedReceipt.ID}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
// GetPoints handles the GET request to retrieve points for a specific receipt by ID.
// It validates the JWT, fetches the points for the receipt, and returns them.
func GetPoints(w http.ResponseWriter, r *http.Request) {
	if !utils.ValidateJWT(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]

	processedReceipt, err := services.GetPointsByID(id)
	if err != nil {
		if err.Error() == "no receipt found" {
			http.Error(w, "No receipt found for that ID", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	response := map[string]int{"points": processedReceipt.Points}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
