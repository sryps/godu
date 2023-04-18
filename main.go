package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"time"

	. "github.com/blacked/go-zabbix"
	"github.com/spf13/cast"
)

var (
	defaultHost = os.Getenv("ZABBIX_HOST_NAME")
	defaultIP   = os.Getenv("ZABBIX_SERVER")
	defaultPort = 10051
	defaultDir  = os.Getenv("HOME")
	defaultKey  = os.Getenv("ZABBIX_TRAPPER_KEY")
)

func DirSize(path string) string {
	var dirSize int64 = 0

	readSize := func(path string, file os.FileInfo, err error) error {
		if !file.IsDir() {
			dirSize += file.Size()
		}
		return nil
	}

	filepath.Walk(path, readSize)

	sizeMB := float64(dirSize) / 1024 / 1024 / 1024
	sizeGB := int(math.Floor(sizeMB*100) / 100)

	size := strconv.Itoa(sizeGB)
	return size
}

func zabbixSend(dirSize string, dir string) {
	var metrics []*Metric
	metrics = append(metrics, NewMetric(defaultHost, defaultKey, dir))
	metrics = append(metrics, NewMetric(defaultHost, defaultKey, dirSize))

	packet := NewPacket(metrics)

	z := NewSender(defaultIP, defaultPort)
	resp, err := z.Send(packet)

	if err != nil {
		fmt.Printf("Zabbix send failed: %v", err)
	}

	fmt.Println(cast.ToString(resp))
}

func main() {
	var arg string
	flag.StringVar(&arg, "f", "/path/to/your/directory", "Enter your arg")
	flag.Parse()

	i := 0
	for i <= 2 {
		run := DirSize(arg)
		fmt.Printf("%s%s", run, " GB\n")
		fmt.Println(arg)
		zabbixSend(run, arg)
		time.Sleep(1 * time.Hour)
	}
}
