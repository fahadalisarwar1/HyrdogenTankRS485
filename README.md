**Modbus RS485 Data Extraction for Hydrogen**
============================================

**Introduction**
----------------


The program helps to receive the data from the hydrogen sensor

**Requirements**
-----------------
The program requires the following libraries

    go get www.github.com/goburrow/modbus.git
    

**Usage**
-----------

In order to be able to use it.


`go run Pressure.go -delaytime 5 -filename HydrogenTank.csv`

USAGE:

	//  delaytime is the amount of time you want between each reading
	//  filename is name of the file you want to store
	//  COMPORT is the serial port to be used for communicaiton


MODBUS RTU characteristics
--------------------------

	handler.BaudRate = 9600
	handler.DataBits = 8
	handler.Parity = "N"
	handler.StopBits = 1
	handler.SlaveId = 1
	handler.Timeout = 5 * time.Second
	
In order to read register for current pressure

    PressureReg, err := client.ReadHoldingRegisters(0, 2)
    
As mentioned in the documentation of the Pressure sensor, The values of 
register 0 and 1 represents the low and high bits respectively.
To convert them in Big endian order, we can use following statement.

    pCurr := []byte{PressureReg[2], PressureReg[3], PressureReg[0], PressureReg[1]}
   
The 32 bit number is stored in float32 bit form, to convert it into deciaml form,
we need to make use of following lines

    bits := binary.BigEndian.Uint32(bytes)
	float := math.Float32frombits(bits)

These lines of code convert the 32 bit representation to the decimal form, the pressure values are in Pascals.
