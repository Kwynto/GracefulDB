package core

import "sync"

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
}

// func WriteBufferToDisk(writeSignal <-chan bool) {

// }
