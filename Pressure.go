package main

import (
	"encoding/binary"
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/goburrow/modbus"
	"log"
	"math"
	"os"
	"os/exec"
	"runtime"
	"time"
)

func main(){


	var delaytime int
	var filename string
	var COMPORT string
	//	USAGE:
	//  delaytime is the amount of time you want between each reading
	//  filename is name of the file you want to store
	//  COMPORT is the serial port to be used for communicaiton

	flag.IntVar(&delaytime, "delaytime", 600, "Time delay in seconds to be used for logging")
	flag.StringVar(&filename, "filename", "HyrdogenTank.csv", "File name to be used for storing data")
	flag.StringVar(&COMPORT, "port", "COM5", "serial port to be used for communicaiton")

	flag.Parse()


	// Settings for modbus RTU
	handler := modbus.NewRTUClientHandler(COMPORT)
	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.Timeout = 5 * time.Second


	defer handler.Close()

	client := modbus.NewClient(handler)
	for {
		CallClear()
		PressureReg, err := client.ReadHoldingRegisters(0, 2)
		// The values are stored in little endian in 4 bytes of 32 bits. This numbering is done to correct their order
		pCurr := []byte{PressureReg[2], PressureReg[3], PressureReg[0], PressureReg[1]}

		checkError("Curr Pressure reading error", err)
		PressurePa := Float32frombytes(pCurr)



		MaxPressureReg, err := client.ReadHoldingRegisters(2, 2)
		checkError("Max Pressure Reading Error", err)
		pMax := []byte{MaxPressureReg[2], MaxPressureReg[3], MaxPressureReg[0], MaxPressureReg[1]}
		MaxPressure := Float32frombytes(pMax)

		MinPressureReg, err := client.ReadHoldingRegisters(4, 2)
		checkError("Min Pressure Reading Error", err)
		pMin := []byte{MinPressureReg[2], MinPressureReg[3], MinPressureReg[0], MinPressureReg[1]}
		MinPressure := Float32frombytes(pMin)

		TempReg, err := client.ReadHoldingRegisters(6, 1)
		tempC := float32(TempReg[1])/10.0 // The temperature has to be dividied by 10 as per specifications in the documentation

		fmt.Println("Curr Pressure:\t", PressurePa)
		fmt.Println("CMax Pressure:\t", MaxPressure)
		fmt.Println("Curr Pressure:\t", MinPressure)
		fmt.Println("Temperature:\t",tempC)

		d := time.Duration(delaytime)
		WriteDataToCSV(PressurePa, MaxPressure, MinPressure,tempC, filename)
		time.Sleep(d*time.Second)

	}
}


var clear map[string]func() //create a map for storing clear funcs

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok { //if we defined a clear func for that platform:
		value()  //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func Float32frombytes(bytes []byte) float32 {
	bits := binary.BigEndian.Uint32(bytes)
	float := math.Float32frombits(bits)

	return float
}


func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}



func WriteDataToCSV(cp float32, maxp float32, minp float32, temp float32, path string){
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
	if err != nil{
		log.Fatal(err)
	}
	defer f.Close()

	//var data [][]string
	cp_str := fmt.Sprintf("%f", cp)
	p_bars := cp/100000
	p_bars_str := fmt.Sprintf("%f", p_bars) // store values as strings to write them to csv file
	cmax_str := fmt.Sprintf("%f", maxp)
	cmin_str := fmt.Sprintf("%f", minp)
	temp_str := fmt.Sprintf("%.5f", temp)
	curr_time := time.Now()

	w := csv.NewWriter(f)
	err = w.Write([]string{curr_time.String(), cp_str,p_bars_str, cmax_str, cmin_str, temp_str})

	if err != nil{
		log.Fatal(err)
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Fatal(err)
	}
}