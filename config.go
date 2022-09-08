package foxmonger

const (
	MySQLType  = "mysql"
	PostgreSQL = "postgresql"
)

type Table struct {
	Name           string            `mapstructure:"name"`
	BaseMultiplier int               `mapstructure:"base_multiplier"`
	Data           map[string]string `mapstructure:"data"`
}

type Config struct {
	BaseCount   int     `mapstructure:"base_count"`
	Prevalidate bool    `mapstructure:"pre_validate"`
	DBType      string  `mapstructure:"db_type"`
	DBName      string  `mapstructure:"db_name"`
	DBHost      string  `mapstructure:"db_host"`
	DBUser      string  `mapstructure:"db_user"`
	DBPass      string  `mapstructure:"db_pass"`
	DBPort      string  `mapstructure:"db_port"`
	Tables      []Table `mapstructure:"tables"`
}
