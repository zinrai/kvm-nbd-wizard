package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func getShutOffVMs() ([]string, error) {
	cmd := exec.Command("virsh", "list", "--state-shutoff", "--name")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error getting shut off VMs: %v", err)
	}

	vms := strings.Split(strings.TrimSpace(string(output)), "\n")
	var nonEmptyVMs []string
	for _, vm := range vms {
		if vm != "" {
			nonEmptyVMs = append(nonEmptyVMs, vm)
		}
	}

	return nonEmptyVMs, nil
}

func getDiskPaths(vmName string) ([]string, error) {
	cmd := exec.Command("virsh", "dumpxml", vmName)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error getting VM XML: %v", err)
	}

	var diskPaths []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "source file=") {
			path := strings.Split(line, "'")[1]
			diskPaths = append(diskPaths, path)
		}
	}

	return diskPaths, nil
}
