package env

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type photoBolt struct {
	RepoDirectory  string `envconfig:"REPO_DIR"`
	RembgFullPath  string `envconfig:"REMBG_FULL_PATH"`
	FfmpegFullPath string `envconfig:"FFMPEG_FULL_PATH"`
}

var (
	path      = "env/.env"
	PhotoBolt = photoBolt{}
)

func Init() error {
	if err := godotenv.Load(path); err != nil {
		return err
	}

	if err := envconfig.Process("", &PhotoBolt); err != nil {
		return err
	}

	return nil
}
