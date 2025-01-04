package internal

import (
	"errors"
	"log"
	"os/exec"
	"regexp"
	"strconv"
)

func getAlsaVolume(device, control string) (float64, error) {
	out, err := exec.Command("amixer", "-D", device, "get", control).Output()
	if err != nil {
		log.Printf("Failed to get volume: %v", err)
		return 0, errors.New("couldn't get alsa volume")
	}
	return parseAlsaVolume(string(out))
}

func parseAlsaVolume(output string) (float64, error) {
	// example output: "[42%]" -> extract "42"
	r, _ := regexp.Compile("\\[(-?[\\d.]+)dB\\]")
	matched := r.FindStringSubmatch(output)
	if len(matched) < 2 {
		log.Printf("couldn't find output volume in: " + output)
		return 0, errors.New("couldn't get alsa volume")
	}

	num, err := strconv.ParseFloat(matched[1], 64)
	if err != nil {
		log.Printf("coudln't convert volume to int in: " + output)
		return 0, errors.New("couldn't get alsa volume")
	}

	return num, nil
}

func setAlsaVolume(device, control string, volume float64) error {
	currentVolume, err := getAlsaVolume(device, control)
	if err != nil {
		return err
	}

	if volume > currentVolume {
		return decreaseVolume(device, control, volume)
	} else {
		return increaseVolume(device, control, volume)
	}
}

func decreaseVolume(device, control string, expectedVolume float64) error {
	i := 0
	for {
		i += 1

		cmd := exec.Command("amixer", "-D", device, "set", control, "1%-")
		if err := cmd.Run(); err != nil {
			return err
		}

		currentVolume, err := getAlsaVolume(device, control)
		if err != nil {
			return err
		}

		if currentVolume <= expectedVolume || currentVolume == -36.0 {
			return nil
		}

		if i > 100 {
			return errors.New("couldn't find volume to be set")
		}
	}
}

func increaseVolume(device, control string, expectedVolume float64) error {
	i := 0
	for {
		i += 1

		cmd := exec.Command("amixer", "-D", device, "set", control, "1%+")
		if err := cmd.Run(); err != nil {
			return err
		}

		currentVolume, err := getAlsaVolume(device, control)
		if err != nil {
			return err
		}

		if currentVolume >= expectedVolume || currentVolume == 8.0 {
			return nil
		}

		if i > 100 {
			return errors.New("couldn't find volume to be set")
		}
	}
}
