package core

import (
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"
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

type tDescColumn struct {
	DB         string
	Table      string
	Column     string
	Path       string
	Spec       TColumnSpecification
	CurrentRev string
	BucketSize int64
	BucketLog  uint8
}

var (
	CoreProcessing  tCoreProcessing
	fileSystemBlock sync.RWMutex
)

func writeBufferToDisk() bool {
	var rows *[]tRowForStore
	var tempCollect = []string{}
	var deskColumns = make(map[string]tDescColumn)

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

	for name, desc := range CoreProcessing.FileDescriptors {
		if desc.Expire.Compare(tNow) == -1 {
			CoreProcessing.FileDescriptors[name].Descriptor.Close()
			delete(CoreProcessing.FileDescriptors, name)
			slog.Debug("Delete file description", slog.String("desc", name)) // FIXME: Need delete this slog
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

	// возможно под удаление - начало
	for _, row := range *rows {
		for _, col := range row.Row {
			tc := fmt.Sprintf("%s|%s|%s", row.DB, row.Table, col.Field)
			tempCollect = append(tempCollect, tc)
		}
	}
	tempCollect = slices.Compact(tempCollect)

	for _, tStr := range tempCollect {
		tArr := strings.Split(tStr, "|")
		dc := GetDescriptionColumn(tArr[0], tArr[1], tArr[2])
		// deskColumns = append(deskColumns, dc)
		deskColumns[tStr] = dc
	}
	// возможно под удаление - конец

	fileSystemBlock.Lock()
	defer fileSystemBlock.Unlock()

	for _, row := range *rows {
		head := fmt.Sprintf("%d|%d|%d|%d|", row.Id, row.Time, row.Status, row.Shape)
		for _, col := range row.Row {
			fullValue := fmt.Sprintf("%s%s\n", head, col.Value)
			key := fmt.Sprintf("%s|%s|%s", row.DB, row.Table, col.Field)
			dc := deskColumns[key]

			maxBucket := Pow(2, dc.BucketLog)
			hashid := row.Id % maxBucket
			if hashid == 0 {
				hashid = Pow(2, dc.BucketLog)
			}

			fileName := fmt.Sprintf("%s_%d", dc.CurrentRev, hashid)

			_, ok := CoreProcessing.FileDescriptors[key]
			if !ok {
				rwFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
				if err != nil {
					return false
				}
				CoreProcessing.FileDescriptors[key] = tCoreFile{
					Descriptor: rwFile,
					Expire:     tNow.Add(time.Minute),
				}
				// defer rwFile.Close()
			}

			if _, err := CoreProcessing.FileDescriptors[key].Descriptor.WriteString(fullValue); err != nil {
				return false
			}
		}
	}

	return true
}
