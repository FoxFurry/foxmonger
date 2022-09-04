package foxmonger

type FoxMonger interface {
	PopulateDatabase() error
}

type monger struct {
	keys memstore
	conf Config
}

func NewMonger(conf Config) FoxMonger {
	return &monger{
		keys: memstore{},
		conf: conf,
	}
}

func (m *monger) validateAndMarkForeigns() {
	for idx := range m.conf.Tables {

	}
}

func (m *monger) PopulateDatabase() error {

}
