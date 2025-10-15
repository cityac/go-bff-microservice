package env

import (
	"github.com/caarlos0/env/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/openfga/go-sdk/client"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/auth"
	"gitlab.HIDDEN.com/HIDDEN/sms-core/logger"
	"os"
)

type EnvConfig struct {
	KafkaUrl            string `env:"KAFKA_BROKERS" envDefault:"localhost:9092"`
	SchemaRegistryUrl   string `env:"SCHEMA_REGISTRY_URL" envDefault:"http://localhost:8085"`
	HIDDENHost             string `env:"HIDDEN_SERVICE_HOST" envDefault:"http://localhost:6767"`
	HIDDENHost             string `env:"PRICE_VP_SERVICE_HOST" envDefault:"http://localhost:6868"`
	ClientHIDDENHost       string `env:"PRICE_CP_SERVICE_HOST" envDefault:"http://localhost:6870"`
	HIDDENProviderHost string `env:"HIDDEN_PROVIDER_HOST" envDefault:"http://localhost:6871"`
	UserService         string `env:"USER_SERVICE_URL" envDefault:"http://localhost:8765"`
	LogLevel            string `env:"LOG_LEVEL" envDefault:"info"`
	PRODUCTION          bool   `env:"PRODUCTION"`
	DashboardHost       string `env:"DASHBOARD_HOST" envDefault:"http://localhost:8000"`
	PartnerClientHost   string `env:"PARTNER_CLIENT_HOST" envDefault:"http://localhost:8002"`
	PriceClientHost     string `env:"PRICE_CLIENT_HOST" envDefault:"http://localhost:8001"`
	HIDDENClientHost       string `env:"HIDDEN_CLIENT_HOST" envDefault:"http://localhost:8003"`
	HIDDENOauthClientHost  string `env:"HIDDEN_CLIENT_OAUTH_HOST" envDefault:"http://localhost:8999"`
	HIDDENAdminHost        string `env:"HIDDEN_ADMIN_HOST" envDefault:"http://localhost:8004"`
	OrgAdminHost        string `env:"ORG_ADMIN_HOST" envDefault:"http://localhost:8005"`
	HIDDENBrokerHost    string `env:"HIDDEN_BROKER_HOST" envDefault:"http://localhost:6002"`
}

type OFGAConfig struct {
	ApiUrl   string `env:"OPENFGA_API_URL" envDefault:"https://openfga.HIDDEN.com"`
	ApiToken string `env:"OPENFGA_API_TOKEN" envDefault:"HIDDEN"`
	StoreId  string `env:"OPENFGA_STORE_ID" envDefault:"HIDDEN"`
	ModelId  string `env:"OPENFGA_MODEL_ID" envDefault:"HIDDEN"`
}

func NewEnvConfig() *EnvConfig {
	config := EnvConfig{}
	logger.InitLogger("info", os.Stdout)
	if err := env.Parse(&config); err != nil {
		logger.Fatal().Err(err).Msg("Error loading .env file")
	}
	logger.Info().Msgf("CONFIG \n%+v\n", config)
	return &config
}

func NewOFGAClientConfig() *client.ClientConfiguration {
	config := OFGAConfig{}
	if err := env.Parse(&config); err != nil {
		logger.Fatal().Err(err).Msg("Error loading .env file")
	}

	return auth.NewClientConfig(config.ApiUrl, config.StoreId, config.ModelId, config.ApiToken)
}
