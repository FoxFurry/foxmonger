package foxmonger

type Table struct {
	Name           string            `mapstructure:"name"`
	BaseMultiplier int               `mapstructure:"base_multiplier"`
	Data           map[string]string `mapstructure:"data"`
	IsForeign      bool              `mapstructure:"-"` // Generated
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

func (c *Config) GetTableByName(targetName string) *Table {
	for idx := range c.Tables {
		if c.Tables[idx].Name == targetName {
			return &c.Tables[idx]
		}
	}

	return nil
}

func (t *Table) GetRowByName(targetName string) string {
	return t.Data[targetName]
}
