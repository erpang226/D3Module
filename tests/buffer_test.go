package tests

import (
	"bytes"
	"encoding/binary"
	"testing"
)

type Cache struct {
	Length int64
	Split  string
	Bytes  []byte
}

func New() *Cache {
	return &Cache{
		Length: 0,
		Split:  "|",
	}
}

func (c *Cache) Put(data []byte) {
	l := len(data)
	l64 := int64(l)
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.BigEndian, l64)
	binary.Write(buffer, binary.BigEndian, byte('|'))
	binary.Write(buffer, binary.BigEndian, data)
	c.Bytes = buffer.Bytes()
	buffer.Bytes()
}

func (c *Cache) GetLenth(index int64) int64 {
	l64 := c.Bytes[index : index+8]
	var length int64
	buffer := bytes.NewBuffer(l64)
	binary.Read(buffer, binary.BigEndian, &length)
	return length
}

func (c *Cache) GetData(index int64, lenth int64) []byte {
	data := c.Bytes[index+9 : index+9+lenth]
	return data
}

func TestBuffer(t *testing.T) {
	msg := "12345678909876543210"
	c := New()
	c.Put([]byte(msg))
	l64 := c.GetLenth(0)
	t.Logf("length %v", l64)
	data := c.GetData(0, l64)
	t.Logf("data %v", data)
}
