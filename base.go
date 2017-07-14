package stats

type Sink interface {
	Push(x float64) error
	Name() string
}

type Variability interface {
	Min() float64
	Max() float64
	Mean() float64
	StandardDeviation() float64
}

const welford = "B. P. Welford's 1962 algorithm"

type WelfordSink struct {
	Sink
	Variability
}

func NewSink() *WelfordSink {
	return &WelfordSink{}
}

func (s *WelfordSink) Push(x float64) error {
	return nil
}

func (s WelfordSink) Name() string {
	return welford
}
