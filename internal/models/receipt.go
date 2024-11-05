// models.go
package models

// Receipt represents the main structure of a receipt submitted for processing.
// It includes information about the retailer, purchase date and time, items, and total amount.
type Receipt struct {
    Retailer     string `json:"retailer"`     // The name of the retailer or store
    PurchaseDate string `json:"purchaseDate"` // The date of purchase (expected format: YYYY-MM-DD)
    PurchaseTime string `json:"purchaseTime"` // The time of purchase (expected format: HH:MM in 24-hour format)
    Items        []Item `json:"items"`        // List of items in the receipt
    Total        string `json:"total"`        // Total amount paid, formatted as a string (expected format: 0.00)
}

// Item represents a single item on the receipt.
// It includes the item's description and price.
type Item struct {
    ShortDescription string `json:"shortDescription"` // A short description of the item
    Price            string `json:"price"`            // Price of the item, formatted as a string (expected format: 0.00)
}

// ProcessedReceipt represents a receipt after processing.
// It includes a unique ID and the total points awarded based on the receipt rules.
type ProcessedReceipt struct {
    ID     string `json:"id"`
    Points int    `json:"points"`
    Hash   string `json:"hash"` // Add this line
}