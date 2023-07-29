package env

import (
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/lightninglabs/lndclient"
)

type photoBolt struct {
	RepoDirectory    string `envconfig:"REPO_DIR"`
	RembgFullPath    string `envconfig:"REMBG_FULL_PATH"`
	FfmpegFullPath   string `envconfig:"FFMPEG_FULL_PATH"`
	Automatic1111URL string `envconfig:"AUTOMATIC1111_URL"`

	LNDMacaroonHex string `envconfig:"LND_MACAROON_HEX"`
	LNDCertPath    string `envconfig:"LND_TLS_CERT_PATH"`
	LNDGrpcAddr    string `envconfig:"LND_GRPC_ADDR"`

	NostrPrivateKey string `envconfig:"NOSTR_PRIVATE_KEY"`
	NostrRelay      string `envconfig:"NOSTR_RELAY"`

	ImgbbSecret string `envconfig:"IMGBB_SECRET"`
	ServerPort  string `envconfig:"SERVER_PORT"`

	AppEnv string `envconfig:"APP_ENV"`
}

var (
	PhotoBolt = photoBolt{}
)

func Init(path string) error {
	if err := godotenv.Load(path); err != nil {
		return err
	}

	if err := envconfig.Process("", &PhotoBolt); err != nil {
		return err
	}

	return nil
}

func (x photoBolt) IsProd() bool {
	return x.AppEnv == "PROD"
}

func (x photoBolt) LndClientNetwork() lndclient.Network {
	if x.IsProd() {
		return lndclient.NetworkMainnet
	} else {
		return lndclient.NetworkRegtest
	}
}

func (x photoBolt) LnNetwork() *chaincfg.Params {
	if x.IsProd() {
		return &chaincfg.MainNetParams
	} else {
		return &chaincfg.RegressionNetParams
	}
}
