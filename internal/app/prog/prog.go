package prog

import (
	"fmt"
	"os"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

var (
	cBSel0 rpio.Pin
	cBSel1 rpio.Pin
	cReset rpio.Pin
	cDone  rpio.Pin
	spiSS  rpio.Pin
)

//Init sets up the rPI pins prior to use
func Init() {
	if err := rpio.Open(); err != nil {
		panic(err)
	}

	fmt.Println("Setting up SPI...")
	// SPI clock must be between 1MHz and 25MHz
	if err := rpio.SpiBegin(rpio.Spi0); err != nil {
		panic(err)
	}
	rpio.SpiMode(1, 1) // clk rests high, clk rise in middle of bit
	rpio.SpiChipSelect(0)
	rpio.SpiSpeed(1000000) // 1MHz

	fmt.Println("Setting up GPIO...")
	// These are the BCM GPIO numbers, not the rPI pin #
	cBSel0 = rpio.Pin(23)
	cBSel0.Output()

	cBSel1 = rpio.Pin(17)
	cBSel1.Output()

	cReset = rpio.Pin(22)
	cReset.Output()
	cReset.Low() // active-low reset

	cDone = rpio.Pin(27)
	cDone.Input()

	spiSS = rpio.Pin(25)
	spiSS.Output()
	spiSS.Low()
}

func readBinFile(binFileName string) []byte {
	file, err := os.Open(binFileName)
	if err != nil {
		panic(err)
	}

	fileInfo, _ := file.Stat()
	fileData := make([]byte, fileInfo.Size())
	_, err = file.Read(fileData)
	if err != nil {
		panic(err)
	}

	file.Close()

	return fileData
}

// Program takes the given bin file and sends it to the FPGA
func Program(binFileName string) {
	fmt.Println("Programming with " + binFileName + "...")
	fileData := readBinFile(binFileName)

	//Wait in reset for 200us
	time.Sleep(200 * time.Microsecond)

	//Come out of reset
	fmt.Println("\tBringing FPGA out of reset...")
	cReset.High()

	//Wait 1200us
	time.Sleep(1200 * time.Microsecond)

	//Drive SpiSS high to force slave mode
	fmt.Println("\tForcing FPGA into slave mode...")
	spiSS.High()

	//Send 8 clocks on SPI clock
	rpio.SpiTransmit(0x00)

	//Send bin file (MSB first)
	fmt.Println("\tSending binfile...")
	spiSS.Low()
	rpio.SpiTransmit(fileData...)
	spiSS.High()

	//Wait 100 SPI clocks
	fmt.Println("\tWaiting for CDone...")
	dummyData := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0}
	rpio.SpiTransmit(dummyData...)
	rpio.SpiTransmit(dummyData...)

	//Sample CDone as high (else error)
	res := cDone.Read()
	if res != rpio.High {
		panic("CDone did not assert when expected")
	}

	//Send 49 more SPI clocks to allow I/O to come up
	rpio.SpiTransmit(dummyData...)

	//Done
	fmt.Println("Done!")
	rpio.SpiEnd(rpio.Spi0)
	rpio.Close()
}
