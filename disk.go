package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func connectNBD(diskPath, nbdDevice string) error {
	cmd := exec.Command("qemu-nbd", "--connect="+nbdDevice, diskPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to connect NBD: %v, output: %s", err, output)
	}
	fmt.Printf("NBD device %s connected\n", nbdDevice)

	// Waiting time before nbd partition is mapped
	time.Sleep(3 * time.Second)

	return nil
}

func disconnectNBD(nbdDevice string) error {
	cmd := exec.Command("qemu-nbd", "--disconnect", nbdDevice)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to disconnect NBD: %v, output: %s", err, output)
	}
	fmt.Printf("NBD device %s disconnected\n", nbdDevice)
	return nil
}

func getPartitions(nbdDevice string) ([]string, error) {
	cmd := exec.Command("fdisk", "-l", nbdDevice)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get partitions: %v", err)
	}

	var partitions []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, nbdDevice) && strings.Contains(line, "Linux") {
			fields := strings.Fields(line)
			if len(fields) > 0 {
				partitions = append(partitions, fields[0])
			}
		}
	}

	return partitions, nil
}

func getFileSystem(partition string) (string, error) {
	cmd := exec.Command("file", "-s", partition)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get file system info: %v", err)
	}
	return string(output), nil
}
