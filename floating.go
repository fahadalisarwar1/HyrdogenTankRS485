package main

import (
	"encoding/binary"
	"fmt"
	"math"
)

func main(){
	//b1 := []byte{141, 138}
	//b2 := []byte{73, 156}
	data := []byte{73, 156, 141, 138}
	bits_data := binary.BigEndian.Uint32(data)
	fmt.Printf("%x",bits_data)
	number := math.Float32frombits(bits_data)
	fmt.Println(number)
}
