package aqm0802

import (
	"bitbucket.org/gmcbay/i2c"
	"errors"
	"time"
)

type AQM0802 struct {
	bus  *i2c.I2CBus
	addr byte
}

func NewAQM0802(busNumber, addr byte) (lcd *AQM0802, err error) {
	lcd = new(AQM0802)
	lcd.bus, err = i2c.Bus(busNumber)
	if err != nil {
		return
	}
	lcd.addr = addr
	return
}

func (lcd *AQM0802) Init() (err error) {
	//WriteByteBlock(addr,reg,[]value)
	err = lcd.bus.WriteByteBlock(lcd.addr, 0x00,
		[]byte{0x38, 0x39, 0x14, 0x70, 0x56, 0x6c})
	if err != nil {
		return
	}
	time.Sleep(time.Millisecond * 200)
	err = lcd.bus.WriteByteBlock(lcd.addr, 0x00,
		[]byte{0x38, 0x0d, 0x01})
	if err != nil {
		return
	}
	return
}

func (lcd *AQM0802) WriteDisplay(lineNum int, s string) (err error) {
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
	return
}

func (lcd *AQM0802) ClearDisplay() (err error) {
	err = lcd.bus.WriteByte(lcd.addr, 0x00, 0x01)
	return
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
	return
}

func (lcd *AQM0802) moveFirstLine() (err error) {
	err = lcd.bus.WriteByte(
		lcd.addr,
		0x00,
		0x80)
	if err != nil {
		return
	}
	return
}

func (lcd *AQM0802) moveSecondLine() (err error) {
	err = lcd.bus.WriteByte(
		lcd.addr,
		0x00,
		0xc0)
	if err != nil {
		return
	}
	return
}
