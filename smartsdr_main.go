/* SPDX-License-Identifier: GPL-3.0
 *
 * Copyright (C) 2018 Brady OBrien. All Rights Reserved.
 */

package main

import (
	"fmt"
	"os"
	"time"
)

func topError(err error) {
	fmt.Printf("Error in main: %v\n", err)
	os.Exit(1)
}

func main() {
	/* Create and start the discovery client */

	disClient, err := CreateDiscoveryClient("0.0.0.0:4992")
	if err != nil {
		topError(err)
	}

	go disClient.doDiscoveryListen()
	select {
	case radio := <-disClient.radios:
		fmt.Println("Found Radio:", radio)
		disClient.Close()
	case err = <-disClient.errors:
		topError(err)
	case <-time.After(time.Second * 30):
		fmt.Println("Failed to discover any radios in 30 seconds")
		disClient.Close()
	}
	os.Exit(0)
}
