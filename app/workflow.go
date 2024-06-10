package billing

import (
    "go.temporal.io/sdk/workflow"
)

func BillingWorkflow(ctx workflow.Context, billID string) error {
    var bill Bill
    workflow.GetLogger(ctx).Info("Billing workflow started", "BillID", billID)

    // Initialize the bill
    bill = Bill{
        ID:        billID,
        CreatedAt: workflow.Now(ctx),
    }

    // Store the initial bill state
    err := addBill(&bill)
    if err != nil {
        return err
    }

    // Listen for signals to add line items or close the bill
    signalChan := workflow.GetSignalChannel(ctx, "billing-signal")
    for !bill.IsClosed {
        var lineItem LineItem
        var closeBill bool
        selector := workflow.NewSelector(ctx)
        selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
            c.Receive(ctx, &lineItem)
            if lineItem.Description != "" {
                bill.LineItems = append(bill.LineItems, lineItem)
                err := addLineItem(bill.ID, lineItem)
                if err != nil {
                    workflow.GetLogger(ctx).Error("Failed to add line item", "Error", err)
                }
            } else {
                closeBill = true
            }
        })
        selector.Select(ctx)
        if closeBill {
            bill.IsClosed = true
            bill.ClosedAt = workflow.Now(ctx)
            err := updateBill(&bill)
            if err != nil {
                workflow.GetLogger(ctx).Error("Failed to update bill", "Error", err)
				return err
            }
        }
    }

    // Calculate total amount
    var totalAmount float64
    for _, item := range bill.LineItems {
        totalAmount += item.Amount
    }

    workflow.GetLogger(ctx).Info("Billing workflow completed", "BillID", billID, "TotalAmount", totalAmount)
    return nil
}
