package internal

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ParseAlsaMixerVolume(t *testing.T) {
	require.Equal(t, 42, parseAlsaVolume("[42%]"))
	require.Equal(t, 42, parseAlsaVolume("[xxx] [42%]"))
	require.Equal(t, 42, parseAlsaVolume("[42%] [-8dB]"))
}
