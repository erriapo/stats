package stats

type Sink interface {
	Push(x float64) error
	Name() string
}

type WelfordSink struct {
}

func NewSink() *WelfordSink {
	return &WelfordSink{}
}

func (s *WelfordSink) Push(x float64) error {
	return nil
}
