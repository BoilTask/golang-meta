package subsystem

type Subsystem struct {
}

func (s *Subsystem) GetName() string {
	return "Unknown"
}

func (s *Subsystem) Init() error {
	return nil
}

func (s *Subsystem) Start() error {
	return nil
}

func (s *Subsystem) Stop() error {
	return nil
}
