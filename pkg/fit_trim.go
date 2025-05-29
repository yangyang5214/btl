package pkg

type FitTrim struct {
}

func NewFitTrim(input string) *FitTrim {
	return &FitTrim{}
}

func (s *FitTrim) Run() error {
	return nil
}
