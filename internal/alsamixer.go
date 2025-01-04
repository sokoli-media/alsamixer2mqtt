package internal

import (
	"errors"
	"log"
	"os/exec"
	"regexp"
	"strconv"
)

func getAlsaVolume(device, control string) (int, error) {
	out, err := exec.Command("amixer", "-D", device, "get", control).Output()
	if err != nil {
		log.Printf("Failed to get volume: %v", err)
		return 0, errors.New("couldn't get alsa volume")
	}
	return parseAlsaVolume(string(out))
}

func parseAlsaVolume(output string) (int, error) {
	// example output: "[42%]" -> extract "42"
	r, _ := regexp.Compile("\\[(\\d+)%\\]")
	matched := r.FindStringSubmatch(output)
	if len(matched) < 2 {
		log.Printf("couldn't find output volume in: " + output)
		return 0, errors.New("couldn't get alsa volume")
	}

	num, err := strconv.Atoi(matched[1])
	if err != nil {
		log.Printf("coudln't convert volume to int in: " + output)
		return 0, errors.New("couldn't get alsa volume")
	}

	return num, nil
}

func setAlsaVolume(device, control string, volume int) error {
	cmd := exec.Command("amixer", "-D", device, "set", control, strconv.Itoa(volume))
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
