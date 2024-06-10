package billing

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
)

// TestCreateBill tests the CreateBill API.
func TestCreateBill(t *testing.T) {
    req := &CreateBillRequest{
        BillID:   "test-bill-id",
        Currency: "USD",
    }

    err := CreateBill(context.Background(), req)
    assert.NoError(t, err)
}

// TestAddLineItem tests the AddLineItem API.
func TestAddLineItem(t *testing.T) {
    req := &AddLineItemRequest{
        BillID:      "test-bill-id",
        Description: "Test item",
        Amount:      10.0,
        Currency:    "USD",
    }

    err := AddLineItem(context.Background(), req)
    assert.NoError(t, err)
}

// TestCloseBill tests the CloseBill API.
func TestCloseBill(t *testing.T) {
    req := &CloseBillRequest{
        BillID: "test-bill-id",
    }

    err := CloseBill(context.Background(), req)
    assert.NoError(t, err)
}

// TestQueryOpenBills tests the QueryOpenBills API.
func TestQueryOpenBills(t *testing.T) {
    resp, err := QueryOpenBills(context.Background())
    assert.NoError(t, err)
    assert.NotNil(t, resp)
}

// TestQueryClosedBills tests the QueryClosedBills API.
func TestQueryClosedBills(t *testing.T) {
    resp, err := QueryClosedBills(context.Background())
    assert.NoError(t, err)
    assert.NotNil(t, resp)
}
