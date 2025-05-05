// Package main contains the main function for the ports cli example.
package main

import (
	"fmt"
	"log/slog"
	"time"

	"go.bug.st/serial/enumerator"
)

func main() {
	for {
		doStuff()
	}
}

func doStuff() {
	targetID := "1a86:7523" // the esp boards

	var found *enumerator.PortDetails
	for {
		ports, err := enumerator.GetDetailedPortsList()
		if err != nil {
			slog.Error(err.Error())
		}
		if len(ports) == 0 {
			fmt.Println("No serial ports found!")
			time.Sleep(1 * time.Second)
			continue
		}
		fmt.Printf("Found %d serial ports\n", len(ports))
		for _, port := range ports {
			fmt.Printf("Found port: %s\n", port.Name)
			if port.IsUSB {
				fmt.Printf("   USB ID     %s:%s\n", port.VID, port.PID)
				fmt.Printf("   USB serial %s\n", port.SerialNumber)
				if port.VID+":"+port.PID == targetID {
					if found != nil {
						slog.Error("Multiple ports found with expected ID", "id", targetID)
						time.Sleep(1 * time.Second)
						continue
					}
					found = port
				}
			}
			if port.Product != "" {
				fmt.Printf("   Product    %s\n", port.Product)
			}
		}

		if found != nil {
			break
		}

		time.Sleep(1 * time.Second)
	}
}
