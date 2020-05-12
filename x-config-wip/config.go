package config

type Config struct {
	AppName string `env:"APP_NAME" envDefault:"go-gettingstarted"`
}

// Parse uses github.com/caarlos0/env or other lib to parse environment variables
func Parse() (Config, error) {
	panic("IMPLEMENT ME!")
}

// ParseManual uses os.LookupEnv or os.GetEnv to read environment variables
func ParseManual() (Config, error) {
	panic("IMPLEMENT ME!")
}
