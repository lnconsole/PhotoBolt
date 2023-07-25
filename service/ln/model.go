package ln

import (
	"context"
	"time"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/lightninglabs/lndclient"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/verrpc"
	"github.com/lightningnetwork/lnd/lntypes"
	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/lightningnetwork/lnd/routing/route"
	"github.com/lightningnetwork/lnd/zpay32"
	"github.com/lnconsole/photobolt/shared"
)

var (
	lndService *lndclient.GrpcLndServices
	lndNetwork lndclient.Network
	lnNetwork  *chaincfg.Params
)

func Init(
	LNDMacaroonHex string,
	LNDCertPath string,
	LNDGrpcAddr string,
	LNDNetwork lndclient.Network,
	LNNetwork *chaincfg.Params,
) error {
	lndNetwork = LNDNetwork
	lnNetwork = LNNetwork

	config := lndclient.LndServicesConfig{
		LndAddress:        LNDGrpcAddr,
		Network:           lndNetwork,
		CustomMacaroonHex: LNDMacaroonHex,
		TLSPath:           LNDCertPath,
		CheckVersion:      &verrpc.Version{Version: "v0.16.0-beta"},
	}

	var err error
	lndService, err = lndclient.NewLndServices(&config)
	if err != nil {
		return err
	}
	// check if alive
	_, err = lndService.Client.GetInfo(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func PayInvoice(params PaymentParams) (*PaymentData, error) {
	var (
		ctx = context.Background()
	)
	zbolt11, err := zpay32.Decode(params.Bolt11, lnNetwork)
	if err != nil {
		return nil, err
	}
	hash := lntypes.Hash(*zbolt11.PaymentHash)
	// log.Printf("%d %v %d %v",
	// 	lnwire.MilliSatoshi(params.MaxFeeMsat),
	// 	route.NewVertex(zbolt11.Destination),
	// 	lnwire.MilliSatoshi(params.AmountMsat),
	// 	hash,
	// )
	payload := lndclient.SendPaymentRequest{
		MaxFeeMsat:  lnwire.MilliSatoshi(params.MaxFeeMsat),
		Timeout:     time.Second * 10,
		Target:      route.NewVertex(zbolt11.Destination),
		PaymentHash: &hash,
		RouteHints:  zbolt11.RouteHints,
		Invoice:     params.Bolt11,
	}
	if shared.EmptyMsatIfNil(zbolt11.MilliSat) == 0 {
		payload.AmountMsat = lnwire.MilliSatoshi(params.AmountMsat)
	}

	updates, errs, err := lndService.Router.SendPayment(
		ctx,
		payload,
	)
	if err != nil {
		return nil, err
	}
	// track payment
	updatech := make(chan PaymentUpdate)
	go func() {
		for {
			select {
			case update := <-updates:
				updatech <- PaymentUpdate{
					// PaymentID:     params.PaymentID,
					Bolt11:        params.Bolt11,
					PaymentStatus: update.State.String(),
					PaidMsat:      int(update.Value),
					FailureReason: update.FailureReason.String(),
					FeeMsat:       int(update.Fee),
					Preimage:      &update.Preimage,
				}
			case err := <-errs:
				if err != nil {
					updatech <- PaymentUpdate{
						// PaymentID:     params.PaymentID,
						Bolt11:        params.Bolt11,
						PaymentStatus: lnrpc.Payment_FAILED.String(),
						FailureReason: lnrpc.PaymentFailureReason_FAILURE_REASON_ERROR.String(),
					}
				}
				return
			}
		}
	}()

	return &PaymentData{
		Update: updatech,
	}, nil
}
