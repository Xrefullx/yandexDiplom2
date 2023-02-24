package models

type Config struct {
	Address        string `env:"RUN_ADDRESS" envDefault:"localhost:8080"`
	DataBaseURI    string `env:"DATABASE_URI"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	SecretKey      string `env:"SECRET_KEY" envDefault:"Xrefullx"`
}
