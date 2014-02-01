package aqm0802

import (
	"bitbucket.org/gmcbay/i2c"
	"errors"
	"fmt"
	"time"
)

type AQM0802 struct {
	bus  *i2c.I2CBus
	addr byte
}

func NewAQM0802(busNumber, addr byte) (lcd *AQM0802, err error) {
	lcd.bus, err = i2c.Bus(busNumber)
	if err != nil {
		return
	}
	lcd.addr = addr
	return
}

func (lcd *AQM0802) Init() (err error) {
	//WriteByteBlock(addr,reg,[]value)
	err = lcd.bus.WriteByteBlock(lcd.addr, 0x40, 0x00,
		[]byte{0x38, 0x39, 0x14, 0x70, 0x56, 0x6c})
	if err != nil {
		return
	}
	time.Sleep(time.Millisecond * 200)
	err = lcd.bus.WriteByteBlock(lcd.addr, 0x40, 0x00,
		[]byte{0x38, 0x39, 0x14, 0x70, 0x56, 0x6c})
	if err != nil {
		return
	}
}

func (lcd *AQM0802) writeDisplay(lineNum int, s string) (err error) {
	switch lineNum {
	case 1:
		err = lcd.moveFirstLine()
		if err != nil {
			return
		}
	case 2:
		err = lcd.moveSecondLine()
		if err != nil {
			return
		}
	default:
		err = errors.New("lineNum is 1 or 2")
		return
	}
	err = lcd.writeLine(s)
	if err != nil {
		return
	}

}
func (lcd *AQM0802) clearDisplay() (err error) {
	err = lcd.bus.WriteByte(
		lcd.addr,
		0x00,
		0x01)
}

func (lcd *AQM0802) writeLine(s string) (err error) {
	code, err := stringToLcdcode(s)
	if err != nil {
		return
	}
	err = lcd.bus.WriteByteBlock(
		lcd.addr,
		0x40,
		code)
	if err != nil {
		return
	}
}

func (lcd *AQM0802) moveFirstLine() (err error) {
	err = lcd.bus.WriteByte(
		lcd.addr,
		0x00,
		0x80)
	if err != nil {
		return
	}
}

func (lcd *AQM0802) moveSecondLine() (err error) {
	err = lcd.bus.WriteByte(
		lcd.addr,
		0x00,
		0xc0)
	if err != nil {
		return
	}
}
func stringToLcdcode(s string) (code [8]byte, err error) {
	//initialize with white space code
	for i := 0; i < len(code); i++ {
		code[i] = 0xa0
	}
	runeArray := []rune(s)
	if len(runeArray) > 8 {
		err = errors.New("Exceeded number of charactor per one line")
		return
	}
	for i, r := range runeArray {
		code[i], err = charToLcdcode(r)
	}
	return
}

func charToLcdcode(r rune) (code byte, err error) {
	switch string(r) {
	case "A":
		code = 0x41
	case "B":
		code = 0x42
	case "C":
		code = 0x11
	default:
		err = errors.New("Unsupported string: " + string(r))
	}
	return
}
