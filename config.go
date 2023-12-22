package foxmonger

type Table struct {
	Name           string            `mapstructure:"name"`
	BaseMultiplier int               `mapstructure:"base_multiplier"`
	BatchSize      int               `mapstructure:"batch_size"`
	Data           map[string]string `mapstructure:"data"`
	ExportQueries  bool              `mapstructure:"export_queries"`
	ExportPath     string            `mapstructure:"export_path"`
	Dummy          bool              `mapstructure:"dummy"`
}

type Config struct {
	DBType string `env:"DB_TYPE"`
	DBName string `env:"DB_NAME"`
	DBHost string `env:"DB_HOST"`
	DBUser string `env:"DB_USER"`
	DBPass string `env:"DB_PASS"`
	DBPort string `env:"DB_PORT"`
}
