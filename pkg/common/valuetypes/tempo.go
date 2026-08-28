package valuetypes

import "fmt"

const (
	minTempo = 40
	maxTempo = 300
)

var ErrInvalidTempo = fmt.Errorf("tempo must be greater than or equal to %d and less than or equal to %d", minTempo, maxTempo)

// Tempo represents track BPM (limited 40 - 300)
type Tempo struct {
	value int
}

func MakeTempo(value int) (Tempo, error) {
	err := assertTempoValid(value)
	if err != nil {
		return Tempo{}, err
	}
	return Tempo{value: value}, nil
}

func (t Tempo) Value() int {
	return t.value
}

func assertTempoValid(value int) error {
	if value < minTempo || value > maxTempo {
		return ErrInvalidTempo
	}
	return nil
}
