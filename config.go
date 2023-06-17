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
	BaseCount        int     `mapstructure:"base_count"`
	Prevalidate      bool    `mapstructure:"pre_validate"`
	DefaultBatchSize int     `mapstructure:"default_batch_size"`
	DBType           string  `mapstructure:"db_type"`
	DBName           string  `mapstructure:"db_name"`
	DBHost           string  `mapstructure:"db_host"`
	DBUser           string  `mapstructure:"db_user"`
	DBPass           string  `mapstructure:"db_pass"`
	DBPort           string  `mapstructure:"db_port"`
	Tables           []Table `mapstructure:"tables"`
}
