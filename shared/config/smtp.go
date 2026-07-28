package sharedconfig

import (
	// "strconv"

	"github.com/spf13/viper"
)

// type SMTPConfig struct {
// 	Host        string
// 	Port        int
// 	FromAddress string
// 	FromName    string
// 	Username    string
// 	Password    string
// }

// func LoadSMTPConfig(v *viper.Viper) SMTPConfig {
// 	port, _ := strconv.Atoi(Get("SMTP_PORT", v.GetString("smtp.port")))

// 	return SMTPConfig{
// 		Host:        Get("SMTP_HOST", v.GetString("smtp.host")),
// 		Port:        port,
// 		FromAddress: Get("SMTP_FROM_ADDRESS", v.GetString("smtp.from_address")),
// 		FromName:    Get("SMTP_FROM_NAME", v.GetString("smtp.from_name")),
// 		Username:    MustGet("SMTP_USERNAME"),
// 		Password:    MustGet("SMTP_PASSWORD"),
// 	}
// }

type EmailConfig struct {
	APIKey      string
	FromAddress string
	FromName    string
}

func LoadEmailConfig(v *viper.Viper) EmailConfig {
	return EmailConfig{
		APIKey:      MustGet("SENDGRID_API_KEY"),
		FromAddress: Get("EMAIL_FROM_ADDRESS", "onboarding@resend.dev"),
		FromName:    Get("EMAIL_FROM_NAME", "MoneyMate"),
	}
}