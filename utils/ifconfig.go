package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func runIfconfig() (string, error) {
	cmd := exec.Command("ifconfig", "-a")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		// 尝试无参数 ifconfig
		cmd = exec.Command("ifconfig")
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			return "", fmt.Errorf("can't run ifconfig command: %w", err)
		}
	}
	return out.String(), nil
}

func getAllInterfacesMap(output string, filter []string) map[string]map[string]string {
	interfaces := make(map[string]map[string]string)
	lines := strings.Split(output, "\n")
	var currentIface string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, ": flags=") {
			parts := strings.Split(line, ":")
			if len(parts) > 0 {
				currentIface = strings.TrimSpace(parts[0])
				interfaces[currentIface] = make(map[string]string)
			}
			continue
		}

		if currentIface == "" {
			continue
		}

		info := interfaces[currentIface]
		if len(filter) > 0 {
			for _, f := range filter {
				if strings.HasPrefix(line, f) {
					info[f] = strings.TrimSpace(strings.Split(line, f)[1])
				}
			}

		}
	}

	return interfaces
}

func GetNetworkInterfaces(filter []string) (map[string]map[string]string, error) {
	if filter == nil {
		return nil, fmt.Errorf("filter cannot be nil")
	}
	output, err := runIfconfig()
	if err != nil {
		return nil, err
	}
	interfaces := getAllInterfacesMap(output, filter)
	return interfaces, nil
}
