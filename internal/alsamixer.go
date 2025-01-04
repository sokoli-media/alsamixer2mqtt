package internal

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
)

func getAlsaVolume(device, control string) int {
	out, err := exec.Command("amixer", "-D", device, "get", control).Output()
	if err != nil {
		log.Printf("Failed to get volume: %v", err)
		return -1
	}
	return parseAlsaVolume(string(out))
}

func parseAlsaVolume(output string) int {
	// example output: "[42%]" -> extract "42"
	r, _ := regexp.Compile("\\[(\\d+)%\\]")
	matched := r.FindStringSubmatch(output)
	if len(matched) < 2 {
		log.Printf("couldn't find output volume in: " + output)
		return -1
	}

	num, err := strconv.Atoi(matched[1])
	if err != nil {
		log.Printf("coudln't convert volume to int in: " + output)
		return -1
	}

	return num
}

func setAlsaVolume(device, control string, volume int) error {
	cmd := exec.Command("amixer", "-D", device, "set", control, fmt.Sprintf("%d%%", volume))
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to set volume: %v", err)
		return err
	}
	return nil
}
