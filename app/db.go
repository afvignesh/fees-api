package billing

import (
    "context"
    "encore.dev/storage/sqldb"
    "log"
)

var db = sqldb.NewDatabase("billsdb", sqldb.DatabaseConfig{
    Migrations: "./migrations",
})

func addBill(bill *Bill) error {
    ctx := context.Background()
    _, err := db.Exec(ctx, `
        INSERT INTO bills (id, currency, is_closed, created_at, closed_at)
        VALUES ($1, $2, $3, $4, $5)`,
        bill.ID, bill.Currency, bill.IsClosed, bill.CreatedAt, bill.ClosedAt)
    if err != nil {
        log.Println("Error adding bill:", err)
    }
    return err
}

func updateBill(bill *Bill) error {
    ctx := context.Background()
    _, err := db.Exec(ctx, `
        UPDATE bills
        SET is_closed = $1, closed_at = $2
        WHERE id = $3`,
        bill.IsClosed, bill.ClosedAt, bill.ID)
    if err != nil {
        log.Println("Error updating bill:", err)
    }
    return err
}

func addLineItem(billID string, item LineItem) error {
    ctx := context.Background()
    _, err := db.Exec(ctx, `
        INSERT INTO line_items (bill_id, description, amount, currency)
        VALUES ($1, $2, $3, $4)`,
        billID, item.Description, item.Amount, item.Currency)
    if err != nil {
        log.Println("Error adding line item:", err)
    }
    return err
}

func getBillsByStatus(ctx context.Context, isClosed bool) ([]*Bill, error) {
    rows, err := db.Query(ctx, `
        SELECT id, currency, is_closed, created_at, closed_at
        FROM bills
        WHERE is_closed = $1`, isClosed)
    if err != nil {
        log.Println("Error querying bills by status:", err)
        return nil, err
    }
    defer rows.Close()

    var bills []*Bill
    for rows.Next() {
        var bill Bill
        err := rows.Scan(&bill.ID, &bill.Currency, &bill.IsClosed, &bill.CreatedAt, &bill.ClosedAt)
        if err != nil {
            log.Println("Error scanning bill:", err)
            return nil, err
        }

        lineItems, err := getLineItems(ctx, bill.ID)
        if err != nil {
            log.Println("Error getting line items for bill:", err)
            return nil, err
        }
        bill.LineItems = lineItems
        bills = append(bills, &bill)
    }
    return bills, nil
}

func getLineItems(ctx context.Context, billID string) ([]LineItem, error) {
    rows, err := db.Query(ctx, `
        SELECT id, bill_id, description, amount, currency
        FROM line_items
        WHERE bill_id = $1`, billID)
    if err != nil {
        log.Println("Error querying line items:", err)
        return nil, err
    }
    defer rows.Close()

    var items []LineItem
    for rows.Next() {
        var item LineItem
        err := rows.Scan(&item.ID, &item.BillID, &item.Description, &item.Amount, &item.Currency)
        if err != nil {
            log.Println("Error scanning line item:", err)
            return nil, err
        }
        items = append(items, item)
    }
    return items, nil
}
