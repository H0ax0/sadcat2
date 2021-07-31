package utils

import (
	"io/ioutil"
	"os"
	"time"
)

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func GetPathFiles(path string) []string {
	files, _ := ioutil.ReadDir(path)
	var t []string
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			t = append(t, file.Name())
		}
	}
	return t
}

func GetRecentLogs(path string, n int) []string {
	var paths []string
	if !PathExists(path) {
		return paths
	}

	if path[len(path)-1:] != "/" {
		path += "/"
	}
	data := time.Now()
	d, _ := time.ParseDuration("-24h")

	nt := Min(n, len(GetPathFiles(path)))
	for i := 1; i <= nt; {
		if PathExists(path + data.Format("2006-01-02") + ".log") {
			paths = append(paths, data.Format("2006-01-02")+".log")
			i++
		}
		data = data.Add(d)
	}
	return paths
}
