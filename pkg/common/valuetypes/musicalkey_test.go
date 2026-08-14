package valuetypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeKey(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr error
	}{
		{name: "natural major key", value: "C"},
		{name: "natural major key at end of alphabet range", value: "G"},
		{name: "sharp major key", value: "F#"},
		{name: "flat major key", value: "Bb"},
		{name: "minor key with m suffix", value: "Am"},
		{name: "minor key with min suffix", value: "Emin"},
		{name: "sharp minor key", value: "C#m"},
		{name: "flat minor key", value: "Ebmin"},
		{name: "empty string is invalid", value: "", wantErr: errInvalidKey},
		{name: "lowercase root is invalid", value: "c", wantErr: errInvalidKey},
		{name: "root outside A-G is invalid", value: "H", wantErr: errInvalidKey},
		{name: "double accidental is invalid", value: "C##", wantErr: errInvalidKey},
		{name: "unsupported suffix is invalid", value: "Cmaj", wantErr: errInvalidKey},
		{name: "trailing whitespace is invalid", value: "C ", wantErr: errInvalidKey},
		{name: "leading whitespace is invalid", value: " C", wantErr: errInvalidKey},
		{name: "accidental before root is invalid", value: "#C", wantErr: errInvalidKey},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, err := MakeKey(tt.value)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, MusicalKey{}, key)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.value, key.String())
		})
	}
}

func TestMusicalKey_String(t *testing.T) {
	key, err := MakeKey("G#m")
	require.NoError(t, err)

	assert.Equal(t, "G#m", key.String())
}
