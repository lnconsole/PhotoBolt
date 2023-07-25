package ln

import (
	"strings"

	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/zpay32"
)

type L402Data struct {
	AuthToken  string
	AmountMsat int
	Invoice    string
	Macaroon   string
}

func ParseL402(l402 string) (L402Data, error) {
	var (
		data L402Data
	)
	// Remove the 'L402 ' prefix from the string
	dataString := strings.TrimPrefix(l402, "LSAT ")

	// Split the string into individual key-value pairs
	pairs := strings.Split(dataString, ", ")

	// Loop through the pairs and extract the keys and values
	for _, pair := range pairs {
		partKey, partVal, _ := strings.Cut(pair, "=")
		key := strings.TrimSpace(partKey)
		value := strings.Trim(partVal, "\"")
		if key == "macaroon" {
			data.Macaroon = value
		}
		if key == "invoice" {
			data.Invoice = value
		}
	}
	zbolt11, err := zpay32.Decode(data.Invoice, lnNetwork)
	if err != nil {
		return L402Data{}, err
	}

	data.AmountMsat = int(*zbolt11.MilliSat)

	return data, nil
}

func HandleL402(l402 string) (L402Data, error) {
	data, err := ParseL402(l402)
	if err != nil {
		return L402Data{}, err
	}

	paymentData, err := PayInvoice(PaymentParams{
		Bolt11:     data.Invoice,
		MaxFeeMsat: 10000,
	})
	if err != nil {
		return L402Data{}, err
	}

	var (
		preimage  = ""
		authToken = ""
	)
	for update := range paymentData.Update {
		if update.PaymentStatus == lnrpc.Payment_SUCCEEDED.String() {
			preimage = update.Preimage.String()
			break
		}
	}

	// construct auth token
	authToken = "LSAT " + data.Macaroon + ":" + preimage

	data.AuthToken = authToken

	return data, nil
}
