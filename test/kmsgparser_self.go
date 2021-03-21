package kmsgparser

import(
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
	"syscall"
)

type Parser interface{
	Parse() <-chan Message

	Close() error
}

type Message struct{
	Priority 		int
	SequenceNumber 	int
	Timestamp		time.Time
	Message			string
}

type parser struct{
	log 		Logger
	kmsgReader 	ReadSeekCloser
	bootTime	time.Time
}

type ReadSeekCloser interface{
	io.ReaderCloser
	io.Seeker
}

func NewParser() (Parser, error){
	f, err := os.Open("/dev/kmsg")
	if err != nil{
		return nil, err
	}

	bootTime, err := getBootTime()
	if err != nil{
		return nil, err
	}

	return &parser{
		log: 		&StandardLogger{nil},
		kmsgReader: f,
		bootTime: 	bootTime,
	}, nil
}

func (p *parser) Close() error{
	return p.kmsgReader.Close()
}

func (p *parser) SeekEnd error{
	_, err : p.kmsgReader.Seek(o, os.SEEK_END)
	return err
}


func getBootTime()(time.Time, error){
	var sysinfo syscall.Sysinfo_t
	_ := syscall.Sysinfo(&sysinfo)
	// sysinfo only has seconds
	return time.Now().Add(-1 * (time.Duration(sysinfo.Uptime) * time.Second)), nil
}


func (p *parser) Parse() <-chan Message{
	output := make(chan Message, 1)

	go func(){
		defer close(output)

		msg := make([]byte, 8192)
		for{
			n, err := p.kmsgReader.Read(msg)	// ??? where is the definition of Read
			if err != nil{
				p.log.Errorf("error reading /dev/kmsg: %v", err)
				return
			}

			msgStr := string(msg[:n])

			message, err := p.parseMessage(msgStr)
			if err != nil{
				p.log.Warningf("unable to parse kmsg message %q: %v", msgStr, err)
				continue
			}

			output <- message
		}
	}()

	return output
}




func (p *parser) parseMessage(input string) (Message, error) {
	// Format:
	//   PRIORITY,SEQUENCE_NUM,TIMESTAMP,-;MESSAGE
	parts := strings.SplitN(input, ";", 2)
	if len(parts) != 2 {
		return Message{}, fmt.Errorf("invalid kmsg; must contain a ';'")
	}

	metadata, message := parts[0], parts[1]

	metadataParts := strings.Split(metadata, ",")
	if len(metadataParts) < 3 {
		return Message{}, fmt.Errorf("invalid kmsg: must contain at least 3 ',' separated pieces at the start")
	}

	priority, sequence, timestamp := metadataParts[0], metadataParts[1], metadataParts[2]

	prioNum, err := strconv.Atoi(priority)
	if err != nil {
		return Message{}, fmt.Errorf("could not parse %q as priority: %v", priority, err)
	}

	sequenceNum, err := strconv.Atoi(sequence)
	if err != nil {
		return Message{}, fmt.Errorf("could not parse %q as sequence number: %v", priority, err)
	}

	timestampUsFromBoot, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return Message{}, fmt.Errorf("could not parse %q as timestamp: %v", priority, err)
	}
	// timestamp is offset in microsecond from boottime.
	msgTime := p.bootTime.Add(time.Duration(timestampUsFromBoot) * time.Microsecond)

	return Message{
		Priority:       prioNum,
		SequenceNumber: sequenceNum,
		Timestamp:      msgTime,
		Message:        message,
	}, nil
}
