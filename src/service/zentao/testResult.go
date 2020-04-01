package zentaoService

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"github.com/easysoft/zentaoatf/src/model"
	"github.com/easysoft/zentaoatf/src/service/client"
	configUtils "github.com/easysoft/zentaoatf/src/utils/config"
	i118Utils "github.com/easysoft/zentaoatf/src/utils/i118"
	logUtils "github.com/easysoft/zentaoatf/src/utils/log"
	"github.com/easysoft/zentaoatf/src/utils/zentao"
	"github.com/fatih/color"
	"os"
)

func CommitTestResult(report model.TestReport, testTaskId int) {
	conf := configUtils.ReadCurrConfig()
	Login(conf.Url, conf.Account, conf.Password)

	report.ZentaoData = os.Getenv("ZENTAO_DATA")
	report.BuildUrl = os.Getenv("BUILD_URL")
	report.TaskId = testTaskId

	if len(report.ZTFResults) > 0 {
		report.ProductId = report.ZTFResults[0].ProductId
	}

	url := conf.Url + zentaoUtils.GenApiUri("unittest", "commitResult", "")
	logUtils.Screen(url)

	reportJson, _ := json.Marshal(report)
	logUtils.Screen(string(reportJson))

	resp, ok := client.PostObject(url, report)

	if ok {
		json, err1 := simplejson.NewJson([]byte(resp))
		if err1 == nil {
			result, err2 := json.Get("result").String()
			if err2 != nil || result != "success" {
				ok = false
			}
		} else {
			ok = false
		}
	}

	msg := "\n"
	if ok {
		msg += color.GreenString(i118Utils.I118Prt.Sprintf("success_to_submit_unit_test_result"))
	} else {
		msg += color.RedString(i118Utils.I118Prt.Sprintf("fail_to_submit_unit_test_result", url))
		msg += "\n" + i118Utils.I118Prt.Sprintf("server_return")
		msg += "\n" + resp
	}

	logUtils.Screen(msg)

	if report.Fail > 0 || !ok {
		os.Exit(1)
	}
}
