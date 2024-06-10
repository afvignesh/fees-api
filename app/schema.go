package billing

import ("time")

// Bill represents a billing record
type Bill struct {
    ID        string    `json:"id"`
    Currency  string    `json:"currency"`
    IsClosed  bool      `json:"is_closed"`
    CreatedAt time.Time `json:"created_at"`
    ClosedAt  time.Time `json:"closed_at"`
    LineItems []LineItem `json:"line_items"`
}

// LineItem represents a line item in a bill
type LineItem struct {
    ID          int     `json:"id"`
    BillID      string  `json:"bill_id"`
    Description string  `json:"description"`
    Amount      float64 `json:"amount"`
    Currency    string  `json:"currency"`
}
