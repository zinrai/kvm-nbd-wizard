package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetPrefix("VM Disk Mount: ")

	for {
		action := promptForAction()

		switch action {
		case "connect":
			vmName := promptForVMName()
			if vmName == "" {
				log.Println("No VM selected. Restarting...")
				continue
			}

			diskPath := selectDiskImage(vmName)
			if diskPath == "" {
				log.Println("No disk image selected. Restarting...")
				continue
			}

			if err := connectAndMount(diskPath); err != nil {
				log.Printf("Error during connect and mount: %v\n", err)
			}
		case "disconnect":
			if err := disconnectAndUnmount(); err != nil {
				log.Printf("Error during disconnect and unmount: %v\n", err)
			}
		case "exit":
			fmt.Println("Exiting program.")
			return
		}
	}
}
