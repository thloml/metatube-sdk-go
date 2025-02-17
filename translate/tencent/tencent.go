package tencent

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/metatube-community/metatube-sdk-go/translate"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tmt "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tmt/v20180321"
)

var _ translate.Translator = (*Tencent)(nil)

const api = "https://tmt.tencentcloudapi.com"

type Tencent struct {
	SecretId  string `json:"tencent-secret-id"`
	SecretKey string `json:"tencent-secret-key"`
	ProjectId int64  `json:"tencent-project-id"`
}
type TencentTranslationResponse struct {
	Response struct {
		Source     string `json:"Source"`
		Target     string `json:"Target"`
		TargetText string `json:"TargetText"`
		RequestId  string `json:"RequestId"`
	} `json:"Response"`
}

func (tx *Tencent) Translate(text, from, to string) (result string, err error) {
	credential := common.NewCredential(
		tx.SecretId,
		tx.SecretKey,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "tmt.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := tmt.NewClient(credential, "ap-chengdu", cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := tmt.NewTextTranslateRequest()

	request.SourceText = common.StringPtr(text)
	request.Source = common.StringPtr(parseToBaiduSupportedLanguage(from))
	request.Target = common.StringPtr(parseToBaiduSupportedLanguage(to))
	request.ProjectId = common.Int64Ptr(tx.ProjectId)

	// 返回的resp是一个TextTranslateResponse的实例，与请求对象对应
	response, err := client.TextTranslate(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	if err != nil {
		panic(err)
	}

	var data TencentTranslationResponse
	if err = json.Unmarshal([]byte(response.ToJsonString()), &data); err != nil {
		return
	}
	result = data.Response.TargetText
	return
}

func parseToBaiduSupportedLanguage(lang string) string {
	if lang = strings.ToLower(lang); lang == "" || lang == "auto" /* auto detect */ {
		return "auto"
	}
	switch lang {
	case "zh", "chs", "zh-cn", "zh_cn", "zh-hans":
		return "zh"
	case "cht", "zh-tw", "zh_tw", "zh-hk", "zh_hk", "zh-hant":
		return "cht"
	case "jp", "ja":
		return "ja"
	case "kor", "ko":
		return "kor"
	case "vie", "vi":
		return "vie"
	case "spa", "es":
		return "spa"
	case "fra", "fr":
		return "fra"
	case "ara", "ar":
		return "ara"
	case "bul", "bg":
		return "bul"
	case "est", "et":
		return "est"
	case "dan", "da":
		return "dan"
	case "fin", "fi":
		return "fin"
	case "rom", "ro":
		return "rom"
	case "slo", "sl":
		return "slo"
	case "swe", "sv":
		return "swe"
	}
	return lang
}

func init() {
	translate.Register(&Tencent{})
}
