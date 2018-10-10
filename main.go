package main

import (
	"fmt"
	"github.com/tarm/goserial"
	"os"
	"io"
)

func main() {
	if len(os.Args) < 2 {
		printf("usage: %s <port>\n", os.Args[0])
		return
	}

	port := os.Args[1]
	cport := &serial.Config{
		Name: port,
		Baud: 9600,
	}
	fport, err := serial.OpenPort(cport)
	if err != nil {
		printf("failed open port: %v\n", err)
		os.Exit(1)
		return
	}
	defer fport.Close()

	var buf [4+256]byte
	var read = func(b []byte) bool {
		_, err := io.ReadFull(fport, b)
		if err != nil {
			printf("i/o read failed: %v\n", err)
			return false
		}
		return true
	}
	for {
		if !read(buf[:1]) {
			return
		}
		if buf[0] != 0x5A {
			continue
		}

		if !read(buf[1:2]) {
			return
		}
		if buf[1] != 0x5A {
			continue
		}

		if !read(buf[2:4]) {
			return
		}
		count := int(buf[3]) + 1
		if !read(buf[4:4+count]) {
			return
		}

		decode(buf[:4+count])
	}

}

const (
	flagDataTemperature = 0x01
	flagDataHumidity    = 0x02
	flagDataPressure    = 0x04
	flagDataIAQ         = 0x08
	flagDataGas         = 0x10
	flagDataAltitude    = 0x20
)

func decode(buf []byte) {
	flag := buf[2]
	buf = buf[4:]

	if flag & flagDataTemperature != 0 && len(buf) >= 2 {
		printf("Temperature: %.2f °C\n", float64(int(buf[0]) << 8 | int(buf[1])) / 100)
		buf = buf[2:]
	}
	if flag & flagDataHumidity != 0 && len(buf) >= 2 {
		printf("Humidity:    %.2f %%\n", float64(int(buf[0]) << 8 | int(buf[1])) / 100)
		buf = buf[2:]
	}
	if flag & flagDataPressure != 0 && len(buf) >= 3 {
		printf("Pressure:    %d Pa\n", int(buf[0]) << 16 | int(buf[1]) << 8 | int(buf[2]))
		buf = buf[3:]
	}
	if flag & flagDataIAQ != 0 && len(buf) >= 2 {
		printf("IAQ:         %d\n", int(buf[0] >> 4) << 8 | int(buf[1]))
		buf = buf[2:]
	}
	if flag & flagDataGas != 0 && len(buf) >= 4 {
		printf("Gas:         %d ohm\n", int(buf[0]) << 24 | int(buf[1]) << 16 | int(buf[2]) << 8 | int(buf[0]))
		buf = buf[4:]
	}
	if flag & flagDataAltitude != 0 && len(buf) >= 2 {
		printf("Altitude:    %d m\n", int(buf[0]) << 8 | int(buf[1]))
	}
	printf("\n")
}

func printf(f string, arg ...interface{}) {
	fmt.Fprintf(os.Stderr, f, arg...)
}


