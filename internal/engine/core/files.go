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
	mxFileSystemBlock sync.RWMutex
)

func writeBufferToDisk() bool {
	// This function is complete
	var rSlRows *[]gtypes.TRowForStore

	StWriteBuffer.Block.Lock()
	uWorkBuff := StWriteBuffer.Switch
	switch uWorkBuff {
	case 1:
		StWriteBuffer.Switch = 2
	case 2:
		StWriteBuffer.Switch = 1
	}
	StWriteBuffer.Block.Unlock()

	switch uWorkBuff {
	case 1:
		StWriteBuffer.FirstBox.BlockBuf.Lock()
		defer StWriteBuffer.FirstBox.BlockBuf.Unlock()
		rSlRows = &StWriteBuffer.FirstBox.Area
	case 2:
		StWriteBuffer.SecondBox.BlockBuf.Lock()
		defer StWriteBuffer.SecondBox.BlockBuf.Unlock()
		rSlRows = &StWriteBuffer.SecondBox.Area
	}

	mxFileSystemBlock.Lock()
	defer mxFileSystemBlock.Unlock()

	for _, stRow := range *rSlRows {
		stDBInfo, _ := GetDBInfo(stRow.DB)
		stTableInfo := stDBInfo.Tables[stRow.Table]

		sServiceCol := fmt.Sprintf("%d|%d|1|%d\n", stRow.Id, stRow.Time, stRow.Shape)

		uMaxBucket := Pow(2, stTableInfo.BucketLog)
		uHashId := stRow.Id % uMaxBucket
		if uHashId == 0 {
			uHashId = uMaxBucket
		}

		// sFileName := fmt.Sprintf("%s%s/%s/service/%s_%d", LocalCoreSettings.Storage, dbInfo.Folder, tableInfo.Folder, tableInfo.CurrentRev, hashid)
		sFileName := filepath.Join(StLocalCoreSettings.Storage, stDBInfo.Folder, stTableInfo.Folder, fmt.Sprintf("service/%s_%d", stTableInfo.CurrentRev, uHashId))
		fFileName, err := os.OpenFile(sFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			return false
		}
		if _, err := fFileName.WriteString(sServiceCol); err != nil {
			fFileName.Close()
			return false
		}
		fFileName.Close()

		sHead := fmt.Sprintf("%d|", stRow.Id)
		for _, stCol := range stRow.Row {
			sFullValue := fmt.Sprintf("%s%s\n", sHead, stCol.Value)

			stColInfo := stTableInfo.Columns[stCol.Field]
			// path := fmt.Sprintf("%s%s/%s/", LocalCoreSettings.Storage, colInfo.Parents, colInfo.Folder)
			sPath := filepath.Join(StLocalCoreSettings.Storage, stColInfo.Parents, stColInfo.Folder)

			// fileName := fmt.Sprintf("%s%s_%d", path, tableInfo.CurrentRev, hashid)
			sFileName := filepath.Join(sPath, fmt.Sprintf("%s_%d", stTableInfo.CurrentRev, uHashId))

			fRWFile, err := os.OpenFile(sFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
			if err != nil {
				return false
			}

			if _, err := fRWFile.WriteString(sFullValue); err != nil {
				fRWFile.Close()
				return false
			}
			fRWFile.Close()
		}
	}

	switch uWorkBuff {
	case 1:
		if rand.Intn(100) == 0 {
			StWriteBuffer.FirstBox.Area = nil
		} else {
			StWriteBuffer.FirstBox.Area = StWriteBuffer.FirstBox.Area[:0]
		}
	case 2:
		if rand.Intn(100) == 0 {
			StWriteBuffer.SecondBox.Area = nil
		} else {
			StWriteBuffer.SecondBox.Area = StWriteBuffer.SecondBox.Area[:0]
		}
	}

	return true
}
