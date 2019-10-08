package fileUtils

import (
	"fmt"
	"github.com/easysoft/zentaoatf/res"
	commonUtils "github.com/easysoft/zentaoatf/src/utils/common"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func ReadFile(filePath string) string {
	buf := ReadFileBuf(filePath)
	str := string(buf)
	str = commonUtils.RemoveBlankLine(str)
	return str
}

func ReadFileBuf(filePath string) []byte {
	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []byte(err.Error())
	}

	return buf
}

func WriteFile(filePath string, content string) {
	dir := filepath.Dir(filePath)
	MkDirIfNeeded(dir)

	var d1 = []byte(content)
	err2 := ioutil.WriteFile(filePath, d1, 0666) //写入文件(字节数组)
	check(err2)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func FileExist(path string) bool {
	var exist = true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func MkDirIfNeeded(dir string) {
	if !FileExist(dir) {
		os.MkdirAll(dir, os.ModePerm)
	}
}

func IsDir(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return fi.IsDir()
}

func AbosutePath(pth string) string {
	if !IsAbosutePath(pth) {
		pth, _ = filepath.Abs(pth)
	}

	pth = UpdateDir(pth)

	return pth
}

func IsAbosutePath(pth string) bool {
	return path.IsAbs(pth) ||
		strings.Index(pth, ":") == 1 // windows
}

func UpdateDir(pth string) string {
	sepa := string(os.PathSeparator)

	if strings.LastIndex(pth, sepa) < len(pth)-1 {
		pth += sepa
	}
	return pth
}

func GetFilesFromParams(arguments []string) []string {
	ret := make([]string, 0)

	for _, arg := range arguments {
		if strings.Index(arg, "-") != 0 {
			if arg == "." {
				arg = AbosutePath(".")
			} else if strings.Index(arg, "."+string(os.PathSeparator)) == 0 {
				arg = AbosutePath(".") + arg[2:]
			} else if !IsAbosutePath(arg) {
				arg = AbosutePath(".") + arg
			}

			ret = append(ret, arg)
		} else {
			break
		}
	}

	return ret
}

func ReadResData(path string) string {
	isRelease := commonUtils.IsRelease()

	var jsonStr string
	if isRelease {
		data, _ := res.Asset(path)
		jsonStr = string(data)
	} else {
		jsonStr = ReadFile(path)
	}

	return jsonStr
}

func GetZtfDir() string { // where ztf command in
	var dir string
	arg1 := strings.ToLower(os.Args[0])

	if strings.Index(arg1, "build") > -1 || strings.Index(arg1, "bin") > -1 || strings.Index(arg1, "temp") > -1 { // debug
		dir, _ = os.Getwd()
	} else {
		dir, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	}

	return UpdateDir(dir)
}

func GetCurrDir() string { // where you run command from
	dir, _ := os.Getwd()
	return UpdateDir(dir)
}

func GetLogDir() string {
	path := vari.ZtfDir + constant.LogDir

	dir, _ := ioutil.ReadDir(path)

	regx := `^\d\d\d$`

	numb := 0
	for _, fi := range dir {
		if fi.IsDir() {
			name := fi.Name()
			pass, _ := regexp.MatchString(regx, name)

			if pass { // 999
				name = strings.TrimLeft(name, "0")
				nm, _ := strconv.Atoi(name)

				if nm >= numb {
					numb = nm
				}
			}
		}
	}

	if numb >= 9 {
		numb = 0

		tempDir := path[:len(path)-1] + "-bak" + string(os.PathSeparator) + path[len(path):]
		childDir := path + "bak" + string(os.PathSeparator) + path[len(path):]

		os.RemoveAll(childDir)
		os.Rename(path, tempDir)

		MkDirIfNeeded(path)

		err := os.Rename(tempDir, childDir)
		_ = err
	}

	ret := getLogNumb(numb + 1)

	return UpdateDir(path + ret)
}

func getLogNumb(numb int) string {
	return fmt.Sprintf("%03s", strconv.Itoa(numb))
}

func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
