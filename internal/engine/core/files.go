package core

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
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

	for _, row := range *rows {
		dbInfo, _ := GetDBInfo(row.DB)
		tableInfo := dbInfo.Tables[row.Table]

		serviceCol := fmt.Sprintf("%d|%d|1|%d\n", row.Id, row.Time, row.Shape)

		maxBucket := Pow(2, tableInfo.BucketLog)
		hashid := row.Id % maxBucket
		if hashid == 0 {
			hashid = maxBucket
		}

		// sFileName := fmt.Sprintf("%s%s/%s/service/%s_%d", LocalCoreSettings.Storage, dbInfo.Folder, tableInfo.Folder, tableInfo.CurrentRev, hashid)
		sFileName := filepath.Join(LocalCoreSettings.Storage, dbInfo.Folder, tableInfo.Folder, fmt.Sprintf("service/%s_%d", tableInfo.CurrentRev, hashid))
		srwFile, err := os.OpenFile(sFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			return false
		}
		if _, err := srwFile.WriteString(serviceCol); err != nil {
			srwFile.Close()
			return false
		}
		srwFile.Close()

		head := fmt.Sprintf("%d|", row.Id)
		for _, col := range row.Row {
			fullValue := fmt.Sprintf("%s%s\n", head, col.Value)

			colInfo := tableInfo.Columns[col.Field]
			// path := fmt.Sprintf("%s%s/%s/", LocalCoreSettings.Storage, colInfo.Parents, colInfo.Folder)
			path := filepath.Join(LocalCoreSettings.Storage, colInfo.Parents, colInfo.Folder)

			// fileName := fmt.Sprintf("%s%s_%d", path, tableInfo.CurrentRev, hashid)
			fileName := filepath.Join(path, fmt.Sprintf("%s_%d", tableInfo.CurrentRev, hashid))

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
		if rand.Intn(100) == 0 {
			WriteBuffer.FirstBox.Area = nil
		} else {
			WriteBuffer.FirstBox.Area = WriteBuffer.FirstBox.Area[:0]
		}
	case 2:
		if rand.Intn(100) == 0 {
			WriteBuffer.SecondBox.Area = nil
		} else {
			WriteBuffer.SecondBox.Area = WriteBuffer.SecondBox.Area[:0]
		}
	}

	return true
}
