package core

import (
	"fmt"
	"os"
	"sync"

	"github.com/Kwynto/GracefulDB/internal/engine/basicsystem/gtypes"
)

var (
	fileSystemBlock sync.RWMutex
)

func writeBufferToDisk() bool {
	// This function is complete
	var rows *[]gtypes.TRowForStore

	WriteBuffer.Block.Lock()
	workBuff := WriteBuffer.Switch
	switch workBuff {
	case 1:
		WriteBuffer.Switch = 2
	case 2:
		WriteBuffer.Switch = 1
	}
	WriteBuffer.Block.Unlock()

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

	fileSystemBlock.Lock()
	defer fileSystemBlock.Unlock()

	// надо оптимизировать этот цикл
	for _, row := range *rows {
		dbInfo, _ := GetDBInfo(row.DB)
		tableInfo := dbInfo.Tables[row.Table]
		serviceCol := fmt.Sprintf("%d|%d|1|%d\n", row.Id, row.Time, row.Shape)

		sMaxBucket := Pow(2, tableInfo.BucketLog)
		shashid := row.Id % sMaxBucket
		if shashid == 0 {
			shashid = sMaxBucket
		}

		sFileName := fmt.Sprintf("%s%s/%s/service/%s_%d", LocalCoreSettings.Storage, dbInfo.Folder, tableInfo.Folder, tableInfo.CurrentRev, shashid)
		srwFile, err := os.OpenFile(sFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			fmt.Println("Point 1")
			return false
		}

		if _, err := srwFile.WriteString(serviceCol); err != nil {
			srwFile.Close()
			fmt.Println("Point 2")
			return false
		}
		srwFile.Close()

		// head := fmt.Sprintf("%d|%d|1|%d|", row.Id, row.Time, row.Shape)
		head := fmt.Sprintf("%d|", row.Id)
		for _, col := range row.Row {
			fullValue := fmt.Sprintf("%s%s\n", head, col.Value)

			// -
			dc := GetDescriptionColumn(row.DB, row.Table, col.Field)

			// -
			maxBucket := Pow(2, dc.BucketLog)
			hashid := row.Id % maxBucket
			if hashid == 0 {
				hashid = maxBucket
			}

			// -
			fileName := fmt.Sprintf("%s%s_%d", dc.Path, dc.CurrentRev, hashid)

			rwFile, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
			if err != nil {
				return false
			}

			if _, err := rwFile.WriteString(fullValue); err != nil {
				rwFile.Close()
				return false
			}
			rwFile.Close()
		}
	}

	switch workBuff {
	case 1:
		clear(WriteBuffer.FirstBox.Area)
	case 2:
		clear(WriteBuffer.SecondBox.Area)
	}

	return true
}
