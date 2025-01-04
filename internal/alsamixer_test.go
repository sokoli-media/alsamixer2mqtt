package internal

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ParseAlsaMixerVolume(t *testing.T) {
	parsed, err := parseAlsaVolume("[42%]")
	require.NoError(t, err)
	require.Equal(t, 42, parsed)

	parsed, err = parseAlsaVolume("[42%] [12.00dB]")
	require.NoError(t, err)
	require.Equal(t, 42, parsed)

	parsed, err = parseAlsaVolume("[12] [42%]")
	require.NoError(t, err)
	require.Equal(t, 42, parsed)
}
