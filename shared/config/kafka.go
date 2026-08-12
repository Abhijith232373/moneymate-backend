package sharedconfig

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type KafkaConfig struct {
	Brokers  []string
	Username string
	Password string
	CACert   string 
}

func LoadKafkaConfig(v *viper.Viper) KafkaConfig {
	certPath := MustGet("KAFKA_CA_CERT_PATH")
	certBytes, err := os.ReadFile(certPath)
	if err != nil {
		panic(fmt.Sprintf("read kafka CA cert at %s: %v", certPath, err))
	}
	return KafkaConfig{
		Brokers:  []string{MustGet("KAFKA_BROKER")},
		Username: MustGet("KAFKA_USERNAME"),
		Password: MustGet("KAFKA_PASSWORD"),
		CACert:   string(certBytes),
	}
}