package config

import (
	"fmt"
	"strings"

	"github.com/hacker4257/pet_charity/pkg/logger"
	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig    `mapstructure:"server"`
	Database  DatabaseConfig  `mapstructure:"database"`
	Redis     RedisConfig     `mapstructure:"redis"`
	JWT       JWTConfig       `mapstructure:"jwt"`
	Aliyun    AliyunConfig    `mapstructure:"aliyun"`
	AliPay    AlipayConfig    `mapstructure:"alipay"`
	WechatPay WechatPayConfig `mapstructure:"wechatpay"`
	Kafka     KafkaConfig     `mapstructure:"kafka"`
	Log       LogConfig       `mapstructure:"log"` // ← 新增
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	Charset  string `mapstructure:"charset"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type JWTConfig struct {
	Secret        string `mapstructure:"secret"`
	AccessExpire  int    `mapstructure:"access_expire"`
	RefreshExpire int    `mapstructure:"refresh_expire"`
}

type AliyunConfig struct {
	AccessKeyID     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	SmsSignName     string `mapstructure:"sms_sign_name"`
	SmsTemplateCode string `mapstructure:"sms_template_code"`
}

type WechatPayConfig struct {
	AppID          string `mapstructure:"app_id"`
	MchID          string `mapstructure:"mch_id"`
	APIKeyV3       string `mapstructure:"api_key_v3"`
	SerialNo       string `mapstructure:"serial_no"`
	PrivateKeyPath string `mapstructure:"private_key_path"`
	NotifyURL      string `mapstructure:"notify_url"`
}

type AlipayConfig struct {
	AppID               string `mapstructure:"app_id"`
	PrivateKeyPath      string `mapstructure:"private_key_path"`
	AlipayPublicKeyPath string `mapstructure:"alipay_public_key_path"`
	NotifyURL           string `mapstructure:"notify_url"`
	ReturnURL           string `mapstructure:"return_url"`
}

type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
	GroupID string   `mapstructure:"group_id"`
}

type LogConfig = logger.Config

var Global *Config

func Load(path string) error {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("read config failed: %w", err)
	}

	// 支持环境变量覆盖，前缀 PET，如 PET_JWT_SECRET
	viper.SetEnvPrefix("PET")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	Global = &Config{}
	if err := viper.Unmarshal(Global); err != nil {
		return fmt.Errorf("unmarshal config failed: %w", err)
	}

	return nil
}
