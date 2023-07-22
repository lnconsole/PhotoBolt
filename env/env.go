package env

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type photoBolt struct {
	RepoDirectory    string `envconfig:"REPO_DIR"`
	RembgFullPath    string `envconfig:"REMBG_FULL_PATH"`
	FfmpegFullPath   string `envconfig:"FFMPEG_FULL_PATH"`
	Automatic1111URL string `envconfig:"AUTOMATIC1111_URL"`
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
