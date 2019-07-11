package main

import (
	"biz"
	"flag"
	"model"
	"os"
	"utils"
)

func main() {
	langType := flag.String("l", "", "Script Language like python, php etc.")
	workDir := flag.String("p", "", "Folder that contains the scripts")

	flag.Parse()

	if *langType == "" || *workDir == "" {
		flag.Usage()
		os.Exit(1)
	}

	files, _ := utils.GetAllFiles(*workDir, *langType)

	var report = model.TestReport{Path: *workDir, Env: utils.GetOs(),
		Pass: 0, Fail: 0, Total: 0, Cases: make([]model.CaseLog, 0)}

	biz.RunScripts(files, *workDir, *langType, &report)

	biz.CheckResults(*workDir, *langType, &report)
	biz.Print(report, *workDir)
}
