package app

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func processLock(pidDir string) (*os.File, string, error) {
	pidDir = strings.TrimSuffix(pidDir, "/") + "/"
	lockfileName := filepath.Base(os.Args[0])
	lockfileFullName := pidDir + lockfileName + ".pid"

	file, err := os.OpenFile(lockfileFullName, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return file, lockfileFullName, err
	}

	prePIDByte, err := ioutil.ReadAll(file)
	if err != nil {
		return file, lockfileFullName, err
	}

	if len(prePIDByte) != 0 {
		prePID, err1 := strconv.Atoi(string(prePIDByte))
		if err1 != nil {
			return file, lockfileFullName, err1
		}

		if err2 := syscall.Kill(prePID, 0); err2 == nil {
			return file, lockfileFullName, fmt.Errorf(
				"existing lock %s: another copy is running as pid %d",
				lockfileFullName,
				prePID,
			)
		}
	}

	file, err = os.OpenFile(lockfileFullName, os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		return file, lockfileFullName, err
	}

	pid := fmt.Sprint(os.Getpid())
	_, err = file.WriteString(pid)

	return file, lockfileFullName, err
}
