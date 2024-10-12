package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func promptForVMName() string {
	vms, err := getShutOffVMs()
	if err != nil {
		log.Printf("Error getting shut off VMs: %v\n", err)
		return ""
	}

	if len(vms) == 0 {
		fmt.Println("No shut off VMs found.")
		return ""
	}

	fmt.Println("Available shut off VMs:")
	for i, vm := range vms {
		fmt.Printf("%d. %s\n", i+1, vm)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Select VM number: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		num, err := strconv.Atoi(input)
		if err == nil && num > 0 && num <= len(vms) {
			return vms[num-1]
		}
		fmt.Println("Invalid selection. Please try again.")
	}
}

func selectDiskImage(vmName string) string {
	diskPaths, err := getDiskPaths(vmName)
	if err != nil {
		log.Printf("Error getting disk paths: %v\n", err)
		return ""
	}

	if len(diskPaths) == 0 {
		log.Println("No disk images found for the VM.")
		return ""
	}

	fmt.Println("Available disk images:")
	for i, path := range diskPaths {
		fmt.Printf("%d. %s\n", i+1, path)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Select disk image number (default is 1): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		fmt.Println("Using default selection: 1")
		return diskPaths[0]
	}

	num, err := strconv.Atoi(input)
	if err != nil || num < 1 || num > len(diskPaths) {
		fmt.Println("Invalid selection. Using default: 1")
		return diskPaths[0]
	}

	return diskPaths[num-1]
}

func promptForAction() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Select action:")
	fmt.Println("1. Connect")
	fmt.Println("2. Disconnect")
	fmt.Println("3. Exit")
	fmt.Print("Enter choice (1/2/3): ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	switch choice {
	case "1":
		return "connect"
	case "2":
		return "disconnect"
	case "3":
		return "exit"
	default:
		fmt.Println("Invalid choice. Please try again.")
		return promptForAction()
	}
}

func getNBDDevice() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter NBD device (default is /dev/nbd0): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return "/dev/nbd0"
	}
	return input
}

func getMountPoint() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter mount point (default is /mnt/vm_partition): ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return "/mnt/vm_partition"
	}
	return input
}
