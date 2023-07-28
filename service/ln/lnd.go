package ln

import (
	"github.com/lightningnetwork/lnd/lntypes"
	"github.com/lightningnetwork/lnd/zpay32"
)

type PaymentParams struct {
	PaymentID  int
	Bolt11     string
	AmountMsat int
	MaxFeeMsat int
}

type PaymentUpdate struct {
	PaymentID     int
	Bolt11        string
	PaymentStatus string
	PaidMsat      int
	FailureReason string
	FeeMsat       int
	Preimage      *lntypes.Preimage
}

type PaymentData struct {
	Update chan PaymentUpdate
}

type InvoiceParams struct {
	AmountMsat    int
	Memo          string
	ExpirySeconds int
}

type InvoiceUpdate struct {
	Bolt11        string
	InvoiceStatus string
	ReceivedMsat  int
}

type InvoiceData struct {
	Bolt11   string
	Zbolt11  *zpay32.Invoice
	Preimage *lntypes.Preimage
	// NodeID   int
	Update chan InvoiceUpdate
}
