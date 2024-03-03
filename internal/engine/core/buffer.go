package core

import (
	"fmt"
	"log/slog"
	"os"
	"slices"
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

type tCoreFile struct {
	Descriptor *os.File
	Expire     time.Time
}

type tCoreProcessing struct {
	FileDescriptors map[string]tCoreFile
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

var (
	signalWrite = make(chan struct{}, 1024)
	signalSD    = make(chan struct{}, 1)
)

var (
	CoreProcessing tCoreProcessing
	// ProccessingBlock sync.RWMutex
)

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
	var rows *[]tRowForStore
	var tempCollect = []string{}

	WriteBuffer.Block.Lock()
	workBuff := WriteBuffer.Switch
	switch workBuff {
	case 1:
		WriteBuffer.Switch = 2
	case 2:
		WriteBuffer.Switch = 1
	}
	WriteBuffer.Block.Unlock()

	tNow := time.Now()

	for nameFile, desc := range CoreProcessing.FileDescriptors {
		if desc.Expire.Compare(tNow) == -1 {
			delete(CoreProcessing.FileDescriptors, nameFile)
			slog.Debug("Delete file description", slog.String("desc", nameFile)) // FIXME: Need delete this slog
		}

	}

	switch workBuff {
	case 1:
		WriteBuffer.FirstBox.BlockBuf.Lock()
		defer WriteBuffer.FirstBox.BlockBuf.Unlock()
		rows = &WriteBuffer.FirstBox.Area
	case 2:
		WriteBuffer.SecondBox.BlockBuf.Lock()
		defer WriteBuffer.SecondBox.BlockBuf.Unlock()
		rows = &WriteBuffer.SecondBox.Area
	}

	for _, row := range *rows {
		for _, col := range row.Row {
			tc := fmt.Sprintf("%s|%s|%s", row.DB, row.Table, col.Field)
			tempCollect = append(tempCollect, tc)
		}
	}
	tempCollect = slices.Compact(tempCollect)

	// for _, tStr := range tempCollect {
	// 	tArr := strings.Split(tStr, "|")
	// 	var fileP []string
	// 	fp := GetFullPathColumn(tArr[0], tArr[1], tArr[2])
	// 	fileP := append(fileP, fp)
	// }

	return false
}
