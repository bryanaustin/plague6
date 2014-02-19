
package main

import (
	"fmt"
)

func Message(message string, args ... interface{}) {
	if !ako.AppConfig.Quiet {
		fmt.Printf(message + "\n", args...)
	}
}