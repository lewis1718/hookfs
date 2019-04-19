package main

import (
	"flag"
	"fmt"
	"hookfs"
	"math/rand"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func probab(percentage int) bool {
	return rand.Intn(99) < percentage
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s [OPTIONS] MOUNTPOINT ORIGINAL...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options\n")
		flag.PrintDefaults()
	}

	mountpoint := flag.String("mountpoint", "/mnt/faultfs", "(required) mount point , must set this option")
	originalDir := flag.String("original", "/mnt/fs", "(required) original dir, must set this option")

	faultType := flag.Int("type", 0, "(required)type value, type list: \n" +
		"   0   OpenFileEIO \n" +
		"   1   OpenFileEPERM \n" +
		"   2   ReadFileDelay \n" +
		"   3   ReadFileErr \n" +
		"   4   WriteFileENOSPC \n" +
		"   5   WriteFileDelay \n" +
		"   6   MkDirEACCES \n" +
		"   7   MkDirEPERM \n" +
		"   8   RmDirEACCESS \n" +
		"   9   RmDirEPERM \n" +
		"   10  FsycnDelay \n" +
		"   11  FsycnEIO \n" +
		"   12  OpenDirEACCESS \n" +
		"   13  OpenDirEPERM \n")


	percent := flag.Int("percent", 0, "fault percentage (0, 99]")
	delay := flag.Duration("delay", 0, "delay time (ms, s, min, h ...) ,if use type (ReadFileDelay , " +
		"WriteFileDelay, " +
		"FsycnDelay), " +
		"must set this option")
	logLevel := flag.Int("log-level", 0, fmt.Sprintf("log level (%d..%d)", hookfs.LogLevelMin, hookfs.LogLevelMax))

	flag.Parse()
	log.Infof("percent %s", *percent)
	log.Infof("delay %s", *delay)
	log.Infof("type %s", *faultType)

	hookfs.SetLogLevel(*logLevel)

	serve(*originalDir, *mountpoint, *faultType, *percent, *delay)
}

func serve(original string, mountpoint string, faulttype int, percent int, delay time.Duration) {
	fs, err := hookfs.NewHookFs(original, mountpoint, &MyHook{faulttype, percent, delay})
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Serving %s", fs)
	log.Infof("Please run `fusermount -u %s` after using this, manually", mountpoint)
	if err = fs.Serve(); err != nil {
		log.Fatal(err)
	}
}
