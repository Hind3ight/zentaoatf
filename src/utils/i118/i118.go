package i118Utils

import (
	"encoding/json"
	"github.com/easysoft/zentaoatf/res"
	"github.com/easysoft/zentaoatf/src/utils/common"
	constant "github.com/easysoft/zentaoatf/src/utils/const"
	fileUtils "github.com/easysoft/zentaoatf/src/utils/file"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"io/ioutil"
)

var I118Prt *message.Printer

func Sprintf(tmpl string, params ...interface{}) (ret string) {
	ret = I118Prt.Sprintf(tmpl, params...)
	return
}

func InitI118(lang string) {
	//var once sync.Once
	//once.Do(func() {
	isRelease := commonUtils.IsRelease()

	var langRes string
	if lang == constant.LanguageEN {
		langRes = constant.EnRes
	} else {
		langRes = constant.ZhRes
	}

	if isRelease {
		data, _ := res.Asset(langRes)
		InitResFromAsset(data)
	} else {
		InitRes(langRes)
	}

	if lang == "zh" {
		I118Prt = message.NewPrinter(language.SimplifiedChinese)
	} else {
		I118Prt = message.NewPrinter(language.AmericanEnglish)
	}
	//})
}

type I18n struct {
	Language string    `json:"language"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Id          string `json:"id"`
	Message     string `json:"message,omitempty"`
	Translation string `json:"translation,omitempty"`
}

func Check(e error) {
	if e != nil {
		panic(e)
	}
}
func ReadI18nJson(file string) string {
	b, err := ioutil.ReadFile(file)
	Check(err)
	str := string(b)
	return str

}

func InitRes(jsonPath string) {
	if !fileUtils.FileExist(jsonPath) { // for debug with ide only, run unit test in another project dir
		if commonUtils.IsWin() {
			jsonPath = "C:\\dev\\project\\go\\ztf\\" + jsonPath
		} else {
			jsonPath = "/Users/aaron/rd/project/zentao/go/ztf/" + jsonPath
		}
	}

	var i18n I18n
	str := ReadI18nJson(jsonPath)
	json.Unmarshal([]byte(str), &i18n)

	msgArr := i18n.Messages
	tag := language.MustParse(i18n.Language)

	for _, e := range msgArr {
		message.SetString(tag, e.Id, e.Translation)
	}
}

func InitResFromAsset(bytes []byte) {
	var i18n I18n
	json.Unmarshal(bytes, &i18n)

	msgArr := i18n.Messages
	tag := language.MustParse(i18n.Language)

	for _, e := range msgArr {
		message.SetString(tag, e.Id, e.Translation)
	}
}
