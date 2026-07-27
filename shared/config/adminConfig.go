package sharedconfig

import "github.com/spf13/viper"

type AdminConfig struct {
    Email    string
    Password string
}

func LoadAdminConfig(v *viper.Viper) AdminConfig {
    return AdminConfig{
        Email:    Get("ADMIN_EMAIL", ""),
        Password: Get("ADMIN_PASSWORD", ""),
    }
}