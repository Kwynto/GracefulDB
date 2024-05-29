package core

import (
	"time"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
)

var StWriteBuffer = gtypes.TCollectBuffers{
	Switch: 1,
}

var (
	chSignalWrite    = make(chan struct{}, 1024)
	chSignalShutdown = make(chan struct{}, 1)
)

func InsertIntoBuffer(stRowsForStore []gtypes.TRowForStore) {
	// This function is complete
	StWriteBuffer.Block.Lock()
	defer StWriteBuffer.Block.Unlock()

	switch StWriteBuffer.Switch {
	case 1:
		StWriteBuffer.FirstBox.Area = append(StWriteBuffer.FirstBox.Area, stRowsForStore...)
	case 2:
		StWriteBuffer.SecondBox.Area = append(StWriteBuffer.SecondBox.Area, stRowsForStore...)
	}
	chSignalWrite <- struct{}{}
}

func WriteBufferService() {
	// This function is complete
labelLoop:
	select {
	case <-chSignalWrite:
		if !writeBufferToDisk() {
			time.Sleep(1 * time.Second)
			chSignalWrite <- struct{}{}
		}
		goto labelLoop
	case <-chSignalShutdown:
		StWriteBuffer.Block.Lock()
		fLen := len(StWriteBuffer.FirstBox.Area) != 0
		sLen := len(StWriteBuffer.SecondBox.Area) != 0
		StWriteBuffer.Block.Unlock()
		if fLen || sLen {
			writeBufferToDisk()
			chSignalShutdown <- struct{}{}
			goto labelLoop
		}
	}
}
