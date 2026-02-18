package main

import (
	"fmt"

	"github.com/amaster507/goschemaless/hl7"
)

func main() {
	// Example usage
	message := "MSH|^~\\&|SendingApp|SendingFac|ReceivingApp|ReceivingFac|20240101120000||ADT^A01|123|P|2.5\rPID|1||PatientID|||Doe^John"
	path := hl7.HL7Path{
		Segment:      "PID",
		SegmentIndex: 1,
		Field:        5,
		Component:    2,
	}

	result, err := hl7.AbstractHL7(message, path)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Result: %s\n", result)
}
