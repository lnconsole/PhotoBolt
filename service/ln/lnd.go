package ln

import "github.com/lightningnetwork/lnd/lntypes"

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
