package internal

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ParseAlsaMixerVolume(t *testing.T) {
	parsed, err := parseAlsaVolume("[42.00dB]")
	require.NoError(t, err)
	require.Equal(t, 42.0, parsed)

	parsed, err = parseAlsaVolume("[-42.00dB]")
	require.NoError(t, err)
	require.Equal(t, -42.0, parsed)

	parsed, err = parseAlsaVolume("[xxx] [42.00dB]")
	require.NoError(t, err)
	require.Equal(t, 42.0, parsed)

	parsed, err = parseAlsaVolume("[8%] [42.00dB]")
	require.NoError(t, err)
	require.Equal(t, 42.0, parsed)
}
