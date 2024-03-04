package core

import (
	"fmt"
	"log/slog"
	"os"
	"slices"
	"sync"
	"time"
)

type tCoreFile struct {
	Descriptor *os.File
	Expire     time.Time
}

type tCoreProcessing struct {
	FileDescriptors map[string]tCoreFile
}

var (
	CoreProcessing  tCoreProcessing
	fileSystemBlock sync.RWMutex
)

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

	fileSystemBlock.Lock()
	defer fileSystemBlock.Unlock()

	// for _, tStr := range tempCollect {
	// 	tArr := strings.Split(tStr, "|")
	// 	var fileP []string
	// 	fp := GetFullPathColumn(tArr[0], tArr[1], tArr[2])
	// 	fileP := append(fileP, fp)
	// }

	return false
}
