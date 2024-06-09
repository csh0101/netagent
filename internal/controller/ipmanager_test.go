package controller_test

import (
	"fmt"
	"testing"

	"github.com/csh0101/netagent.git/internal/controller"
)

func TestIPmanager(t *testing.T) {
	manager, err := controller.NewIPManager("192.168.24.0/24")
	if err != nil {
		fmt.Printf("Error creating IP manager: %v\n", err)
		return
	}

	for i := 0; i < 5; i++ {
		ip, err := manager.GenerateUniqueIP()
		if err != nil {
			fmt.Printf("Error generating IP: %v\n", err)
		} else {
			fmt.Printf("Generated IP: %s\n", ip.String())
		}
	}
}
