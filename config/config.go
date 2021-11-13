package config

type TencentCredential struct {
	SecretId  string `yaml:"SecretId"`
	SecretKey string `yaml:"SecretKey"`
}

type Config struct {
	Email   string            `yaml:"Email"`
	Domains []string          `yaml:"Domains"`
	Tencent TencentCredential `yaml:"Tencent"`
	TTL     uint64            `yaml:"TTL"`
	KeyType string            `yaml:"KeyType"`
}
