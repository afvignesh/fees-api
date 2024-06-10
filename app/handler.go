package billing

import (
    "context"
    "time"
    "log"
    "go.temporal.io/sdk/client"
    "go.temporal.io/sdk/worker"
)

var temporalClient client.Client

func init() {
    var err error
    temporalClient, err = client.Dial(client.Options{})
    if err != nil {
        panic(err)
    }
    
    go startWorker()
}

func startWorker() {
    c, err := client.Dial(client.Options{})
    if err != nil {
        log.Fatalln("Unable to create Temporal client", err)
    }
    defer c.Close()

    w := worker.New(c, "billing", worker.Options{})

    w.RegisterWorkflow(BillingWorkflow)

    err = w.Run(worker.InterruptCh())
    if err != nil {
        log.Fatalln("Unable to start worker", err)
    }
}

type CreateBillRequest struct {
    BillID   string `json:"bill_id"`
    Currency string `json:"currency"`
}

type AddLineItemRequest struct {
    BillID      string  `json:"bill_id"`
    Description string  `json:"description"`
    Amount      float64 `json:"amount"`
    Currency    string  `json:"currency"`
}

type CloseBillRequest struct {
    BillID string `json:"bill_id"`
}

type BillResponse struct {
    ID          string     `json:"id"`
    Currency    string     `json:"currency"`
    LineItems   []LineItem `json:"line_items"`
    IsClosed    bool       `json:"is_closed"`
    CreatedAt   time.Time  `json:"created_at"`
    ClosedAt    time.Time  `json:"closed_at"`
}

// CreateBill creates a new bill
// encore:api public method=POST path=/bills/create
func CreateBill(ctx context.Context, req *CreateBillRequest) error {
    options := client.StartWorkflowOptions{
        ID:        req.BillID,
        TaskQueue: "billing",
    }
    _, err := temporalClient.ExecuteWorkflow(ctx, options, BillingWorkflow, req.BillID)
    return err
}

// AddLineItem adds a line item to an existing open bill
// encore:api public method=POST path=/bills/add-line-item
func AddLineItem(ctx context.Context, req *AddLineItemRequest) error {
    signal := LineItem{
        BillID:      req.BillID,
        Description: req.Description,
        Amount:      req.Amount,
        Currency:    req.Currency,
    }
    return temporalClient.SignalWorkflow(ctx, req.BillID, "", "billing-signal", signal)
}

// CloseBill closes an active bill
// encore:api public method=POST path=/bills/close
func CloseBill(ctx context.Context, req *CloseBillRequest) error {
    return temporalClient.SignalWorkflow(ctx, req.BillID, "", "billing-signal", LineItem{})
}

// QueryOpenBills returns all open bills
// encore:api public method=GET path=/bills/open
func QueryOpenBills(ctx context.Context) (BillsResponse, error) {
    openBills, err := getBillsByStatus(ctx, false)
    if err != nil {
        return BillsResponse{}, err
    }
    return BillsResponse{Bills: convertToBillResponses(openBills)}, nil
}

// QueryClosedBills returns all closed bills
// encore:api public method=GET path=/bills/closed
func QueryClosedBills(ctx context.Context) (BillsResponse, error) {
    closedBills, err := getBillsByStatus(ctx, true)
    if err != nil {
        return BillsResponse{}, err
    }
    return BillsResponse{Bills: convertToBillResponses(closedBills)}, nil
}

func convertToBillResponses(bills []*Bill) []BillResponse {
    var responses []BillResponse
    for _, bill := range bills {
        responses = append(responses, BillResponse{
            ID:        bill.ID,
            Currency:  bill.Currency,
            LineItems: bill.LineItems,
            IsClosed:  bill.IsClosed,
            CreatedAt: bill.CreatedAt,
            ClosedAt:  bill.ClosedAt,
        })
    }
    return responses
}
