package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func mountPartition(partition, mountPoint string) error {
	if err := os.MkdirAll(mountPoint, 0755); err != nil {
		return fmt.Errorf("failed to create mount point: %v", err)
	}

	fsType, err := getFileSystemType(partition)
	if err != nil {
		return fmt.Errorf("failed to get file system type: %v", err)
	}

	cmd := exec.Command("mount", "-t", fsType, partition, mountPoint)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to mount partition: %v, output: %s", err, string(output))
	}

	fmt.Printf("Successfully mounted %s to %s\n", partition, mountPoint)
	return nil
}

func getFileSystemType(partition string) (string, error) {
	cmd := exec.Command("blkid", "-o", "value", "-s", "TYPE", partition)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get file system type: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func unmountPartition(mountPoint string) error {
	cmd := exec.Command("umount", mountPoint)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to unmount partition: %v", err)
	}

	return nil
}

func getMountedPartitions(nbdDevice string) ([]string, error) {
	cmd := exec.Command("mount")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get mounted partitions: %v", err)
	}

	var mountedPartitions []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, nbdDevice) {
			fields := strings.Fields(line)
			if len(fields) > 2 {
				mountedPartitions = append(mountedPartitions, fields[2])
			}
		}
	}

	return mountedPartitions, nil
}

func connectAndMount(diskPath string) error {
	nbdDevice := getNBDDevice()

	if err := connectNBD(diskPath, nbdDevice); err != nil {
		return err
	}

	partitions, err := getPartitions(nbdDevice)
	if err != nil {
		disconnectNBD(nbdDevice) // Cleanup on error
		return fmt.Errorf("failed to get partitions: %v", err)
	}

	if len(partitions) == 0 {
		disconnectNBD(nbdDevice) // Cleanup on error
		return fmt.Errorf("no partitions found")
	}

	fmt.Println("Available partitions:")
	for i, partition := range partitions {
		fsInfo, _ := getFileSystem(partition)
		fmt.Printf("%d. %s (%s)\n", i+1, partition, fsInfo)
	}

	selectedPartition := selectPartitionFromList(partitions)
	mountPoint := getMountPoint()

	if err := mountPartition(selectedPartition, mountPoint); err != nil {
		disconnectNBD(nbdDevice) // Cleanup on error
		return fmt.Errorf("failed to mount partition: %v", err)
	}

	fmt.Printf("Partition %s mounted at %s\n", selectedPartition, mountPoint)
	fmt.Printf("NBD device %s remains connected. Use 'disconnect' option to unmount and disconnect.\n", nbdDevice)
	return nil
}

func disconnectAndUnmount() error {
	nbdDevice := getNBDDevice()

	mountedPartitions, err := getMountedPartitions(nbdDevice)
	if err != nil {
		return fmt.Errorf("failed to get mounted partitions: %v", err)
	}

	for _, mountPoint := range mountedPartitions {
		if err := unmountPartition(mountPoint); err != nil {
			fmt.Printf("Error unmounting %s: %v\n", mountPoint, err)
		} else {
			fmt.Printf("Unmounted %s\n", mountPoint)
		}
	}

	if err := disconnectNBD(nbdDevice); err != nil {
		return fmt.Errorf("failed to disconnect NBD: %v", err)
	}

	return nil
}

func selectPartitionFromList(partitions []string) string {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Select partition number: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		num, err := strconv.Atoi(input)
		if err == nil && num > 0 && num <= len(partitions) {
			return partitions[num-1]
		}
		fmt.Println("Invalid selection. Please try again.")
	}
}
