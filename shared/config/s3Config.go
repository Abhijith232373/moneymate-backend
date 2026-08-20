package sharedconfig

import "github.com/spf13/viper"

type S3Config struct {
	Bucket     string
	Region     string
	PublicBase string
}

func LoadS3Config(v *viper.Viper) S3Config {
	return S3Config{
		Bucket:     MustGet("S3_BUCKET"),
		Region:     MustGet("S3_REGION"),
		PublicBase: MustGet("S3_PUBLIC_BASE"),
	}
}