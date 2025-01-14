package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gosnmp/gosnmp"
)

func main() {

	g := &gosnmp.GoSNMP{
		Target:    "127.0.0.1",
		Port:      1321,
		Community: "public",
		Timeout:   time.Second * 10,
		// Logger:    gosnmp.NewLogger(log.New(os.Stdout, "", 0)),
	}

	// Default is a pointer to a GoSNMP struct that contains sensible defaults
	// eg port 161, community public, etc
	err := g.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer g.Conn.Close()

	oids := []string{"1.3.6.1.2.1.1.5"}
	result, err2 := g.Get(oids) // Get() accepts up to g.MAX_OIDS
	if err2 != nil {
		log.Fatalf("Get() err: %v", err2)
	}

	for i, variable := range result.Variables {
		fmt.Printf("%d: oid: %s ", i, variable.Name)

		// the Value of each variable returned by Get() implements
		// interface{}. You could do a type switch...
		switch variable.Type {
		case gosnmp.OctetString:
			bytes := variable.Value.([]byte)
			fmt.Printf("string: %s\n", string(bytes))
		default:
			// ... or often you're just interested in numeric values.
			// ToBigInt() will return the Value as a BigInt, for plugging
			// into your calculations.
			fmt.Printf("number: %d\n", gosnmp.ToBigInt(variable.Value))
		}
	}
}
