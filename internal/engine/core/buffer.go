package core

import (
	"sync"
	"time"
)

type tWriteBuffer struct {
	Area     []tRowForStore
	BlockBuf sync.RWMutex
}

type tCollectBuffers struct {
	FirstBox  tWriteBuffer
	SecondBox tWriteBuffer
	Block     sync.RWMutex
	Switch    uint8
}

var WriteBuffer = tCollectBuffers{
	// FirstBox: tWriteBuffer{
	// 	Area:     []tRowForStore{},
	// 	BlockBuf: sync.RWMutex{},
	// },
	// SecondBox: tWriteBuffer{
	// 	Area:     []tRowForStore{},
	// 	BlockBuf: sync.RWMutex{},
	// },
	// Block:  sync.RWMutex{},
	Switch: 1,
}

var signalWrite = make(chan struct{}, 1024)
var signalSD = make(chan struct{}, 1)

func InsertIntoBuffer(rowsForStore []tRowForStore) {
	// -
	WriteBuffer.Block.Lock()
	defer WriteBuffer.Block.Unlock()

	switch WriteBuffer.Switch {
	case 1:
		WriteBuffer.FirstBox.Area = append(WriteBuffer.FirstBox.Area, rowsForStore...)
	case 2:
		WriteBuffer.SecondBox.Area = append(WriteBuffer.SecondBox.Area, rowsForStore...)
	}
	signalWrite <- struct{}{}
}

func WriteBufferService() {
loop:
	select {
	case <-signalWrite:
		if !writeBufferToDisk() {
			time.Sleep(1 * time.Second)
			signalWrite <- struct{}{}
		}
		goto loop
	case <-signalSD:
		WriteBuffer.Block.Lock()
		fLen := len(WriteBuffer.FirstBox.Area) != 0
		sLen := len(WriteBuffer.SecondBox.Area) != 0
		WriteBuffer.Block.Unlock()
		if fLen || sLen {
			writeBufferToDisk()
			signalSD <- struct{}{}
			goto loop
		}
	}
}

func writeBufferToDisk() bool {

	return false
}
