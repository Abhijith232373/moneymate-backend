// shared/config/razorpay_config.go
package sharedconfig

import "github.com/spf13/viper"

type RazorpayConfig struct {
	KeyID     string
	KeySecret string
}

func LoadRazorpayConfig(v *viper.Viper) RazorpayConfig {
	return RazorpayConfig{
		KeyID:     MustGet("RAZORPAY_KEY_ID"),
		KeySecret: MustGet("RAZORPAY_KEY_SECRET"),
	}
}