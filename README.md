**Modbus RS485 Data Extraction for Hydrogen**
============================================

**Introduction**
----------------


The program helps to receive the data from the hydrogen sensor


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
    
