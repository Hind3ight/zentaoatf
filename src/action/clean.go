package action

import (
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/vari"
	"os"
)

func Clean() {
	path := vari.ExeDir + constant.LogDir
	bak := path[:len(path)-1] + "-bak" + string(os.PathSeparator) + path[len(path):]

	os.RemoveAll(path)
	os.RemoveAll(bak)

	logUtils.PrintTo(i118Utils.Sprintf("success_to_clean_logs"))
}
