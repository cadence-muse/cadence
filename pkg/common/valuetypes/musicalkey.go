package valuetypes

import (
	"errors"
	"regexp"
)

const keyRegexString = "^[A-G][b\\#]?(m|min)?$"

var (
	keyRegex      = regexp.MustCompile(keyRegexString)
	errInvalidKey = errors.New("invalid key")
)

// MusicalKey represents musical key (e.g. C, F#, Bb, G#m)
type MusicalKey struct {
	value string
}

func MakeKey(value string) (MusicalKey, error) {
	err := assertKeyValid(value)
	if err != nil {
		return MusicalKey{}, err
	}
	return MusicalKey{value: value}, nil
}

func (k MusicalKey) String() string {
	return k.value
}

func assertKeyValid(value string) error {
	if !keyRegex.MatchString(value) {
		return errInvalidKey
	}
	return nil
}
