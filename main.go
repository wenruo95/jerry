/*================================================================
*   Copyright (C) 2022. All rights reserved.
*
*   file : main.go
*   coder: zemanzeng
*   date : 2022-01-16 10:51:58
*   desc : jerry: a private msg storehouse
*
================================================================*/

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wenruo95/jerry/common/log"
)

const BuildVersion string = "v0.0.1"

var (
	// build info
	BuildDate     string
	BuildCommitId string

	fi *FlagInfo
)

type FlagInfo struct {
	Version  bool
	WorkDir  string
	UserName string
}

func init() {
	fi = new(FlagInfo)
	flag.BoolVar(&fi.Version, "version", false, "-version show app version info")
	flag.StringVar(&fi.WorkDir, "dir", "./workspace", "-dir=file save dir")
	flag.StringVar(&fi.UserName, "name", "jack", "-dir=receiver username")
	flag.Parse()
}

func checkFlagInfo() {
	if fi.Version {
		fmt.Printf("BUILD_DATE:%s BUILD_COMMIT_ID:%s BUILD_VERSION:%s\n", BuildDate, BuildCommitId, BuildVersion)
		os.Exit(0)
	}

	if len(fi.WorkDir) == 0 {
		fmt.Printf("[ERROR] work dir length is zero")
		os.Exit(0)
	}

	if len(fi.UserName) == 0 {
		fmt.Printf("[ERROR] username length is zero")
		os.Exit(0)
	}

}

func main() {
	checkFlagInfo()
	log.RegisterWithFileName(log.DefaultLogName, "debug", "./log/jerry.log")
	defer log.Sync()

	log.Info("server launch...")

}
