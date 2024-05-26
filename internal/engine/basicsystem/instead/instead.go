package instead

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Kwynto/GracefulDB/pkg/lib/closer"
	"github.com/Kwynto/GracefulDB/pkg/lib/ecowriter"
)

// NOTE: This package does not require data synchronization, it is written with a lack of significant competitiveness.

type TRow []string
type TFile []TRow

type tInstead map[string]TFile

var mCache = make(tInstead)

var chTick = make(chan struct{}, 1)

// Private

func (mi tInstead) set(sKey string, slFile TFile) error {
	// This method is complete
	mi[sKey] = slFile

	_, isOk := mi[sKey]
	if isOk {
		return nil
	}

	return errors.New("the cache value cannot be set")
}

func (mi tInstead) get(sKey string) (TFile, error) {
	// This method is complete
	mValue, isOk := mi[sKey]
	if isOk {
		return mValue, nil
	}

	return TFile{}, errors.New("it is impossible to get the cache value")
}

func (mi tInstead) loadFile(sPath string, iStamp int64) (TFile, bool) {
	// This method is complete
	var slFile = make(TFile, 0, 4)

	sKey := fmt.Sprintf("%s:%d", sPath, iStamp)
	slValue, isOkStart := mi[sKey]
	if isOkStart {
		fmt.Println("Отдали структуру") // FIXME: удалить
		return slValue, true
	}

	sFileText, err := ecowriter.FileRead(sPath)
	if err != nil {
		// return TFile{}, errors.New("file reading error")
		return TFile{}, false
	}

	slSFileData := strings.Split(sFileText, "\n")

	for _, sLine := range slSFileData {
		slLineData := strings.Split(sLine, "|")
		if len(slLineData) < 1 {
			continue
		}
		if slLineData[0] == "" {
			continue
		}
		slFile = append(slFile, slLineData)
	}

	slFile = slices.Clip(slFile)

	_, isOkFinish := mi[sKey]
	if !isOkFinish {
		mi[sKey] = slFile
		fmt.Println("Записали структуру") // FIXME: удалить
	}

	fmt.Println("... и отдали.") // FIXME: удалить
	return slFile, true
}

func (mi tInstead) remove(sKey string) bool {
	// This method is complete
	delete(mi, sKey)
	_, isOk := mi[sKey]

	return !isOk
}

func (mi tInstead) destroyStamps(iStamp int64) {
	// This method is complete
	var slSBuffer = make([]string, 0, 4)
	for sKey := range mi {
		slSKey := strings.Split(sKey, ":")
		if len(slSKey) < 2 {
			continue
		}

		sCacheStump := slSKey[1]
		iCacheStump, err := strconv.ParseInt(sCacheStump, 10, 64)
		if err != nil {
			continue
		}

		if iCacheStump == iStamp {
			slSBuffer = append(slSBuffer, sKey)
		}
	}

	for _, sKey := range slSBuffer {
		delete(mi, sKey)
	}
}

func (mi tInstead) ifInstalled(sKey string) bool {
	// This method is complete
	_, isOk := mi[sKey]
	return isOk
}

func (mi tInstead) cleaning() {
	// This method is complete
	var slSBuffer = make([]string, 0, 4)
	for {
		<-chTick
		// cleaning
		dtNow := time.Now().Unix()
		for sKey := range mi {
			slSKey := strings.Split(sKey, ":")
			if len(slSKey) < 2 {
				continue
			}

			sCacheStump := slSKey[1]
			iCacheStump, err := strconv.ParseInt(sCacheStump, 10, 64)
			if err != nil {
				continue
			}

			iCacheStump = iCacheStump + 20
			if iCacheStump < dtNow {
				slSBuffer = append(slSBuffer, sKey)
			}
		}
		for _, sKey := range slSBuffer {
			delete(mi, sKey)
			fmt.Println("Почистили", sKey) // FIXME: удалить
		}
		if rand.Intn(100) == 0 {
			slSBuffer = nil
		} else {
			slSBuffer = slSBuffer[:0]
		}
	}
}

func (mi tInstead) serv() {
	// This method is complete
	for {
		time.Sleep(5 * time.Second)
		chTick <- struct{}{}
	}
}

// Public

func Set(sPrefix string, iStamp int64, slFile TFile) error {
	// This function is complete
	sKey := fmt.Sprintf("%s:%d", sPrefix, iStamp)
	return mCache.set(sKey, slFile)
}

func Get(sPrefix string, iStamp int64) (TFile, error) {
	// This function is complete
	sKey := fmt.Sprintf("%s:%d", sPrefix, iStamp)
	return mCache.get(sKey)
}

func LoadFile(sPath string, iStamp int64) (TFile, bool) {
	// This function is complete
	return mCache.loadFile(sPath, iStamp)
}

func Remove(sPrefix string, iStamp int64) bool {
	// This function is complete
	sKey := fmt.Sprintf("%s:%d", sPrefix, iStamp)
	return mCache.remove(sKey)
}

func DestroyStamps(iStamp int64) {
	// This function is complete
	mCache.destroyStamps(iStamp)
}

func IfInstalled(sPrefix string, iStamp int64) bool {
	// This function is complete
	sKey := fmt.Sprintf("%s:%d", sPrefix, iStamp)
	return mCache.ifInstalled(sKey)
}

// Main

func Start() {
	// This function is complete
	go mCache.cleaning()
	go mCache.serv()
	slog.Info("The cache of the DBMS was started.")
}

func Shutdown(ctx context.Context, c *closer.TCloser) {
	// This function is complete
	slog.Info("The caching system is stopped.")
	c.Done()
}
