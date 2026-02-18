# goschemaless

HL7 message parsing library for Go.

## Installation

```bash
go get github.com/amaster507/goschemaless
```

## Usage

```go
import (
    "fmt"
    "log"

    "github.com/amaster507/goschemaless/hl7"
)

message := "MSH|^~\\&|SendingApp|SendingFac|ReceivingApp|ReceivingFac|20240101120000||ADT^A01|123|P|2.5\rPID|1||PatientID||Doe^John"

path := hl7.ParsePath("PID-5.2")

// direct representation as hl7.ParsePath:
path = hl7.HL7Path{
    Segment:      "PID",
    SegmentIndex: 1,
    Field:        5,
    Component:    2,
}

result, err := hl7.AbstractHL7(message, path)
if err != nil {
    log.Fatal(err)
}
fmt.Println(result) // Output: John
```

## Publishing

When ready to publish:

```bash
git tag v0.1.0
git push origin v0.1.0
```

Others can then:

```bash
go get github.com/amaster507/goschemaless@v0.1.0
```
