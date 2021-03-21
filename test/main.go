package main

import(
	"fmt"
	"time"
	"log"

	"kmsgparser"
)

func main(){
	parser, err := kmsgparser.NewParser()
	if err != nil{
		log.Fatalf("unable to create parser: %v", err)
	}

	defer parser.Close()

	err := parser.SeekEnd()
	if err != nil{
		log.Fatalf("ERROR >> could not tail: %v", err)
	}

	kmsg := parser.Parse()

	for msg := range kmsg{
		fmt.Printf("(%d) - %s: %s", msg.SequenceNumber, msg.Timestamp.Format(time.RFC3339Nano), msg.Message)
	}
}

