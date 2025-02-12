package xiaoniu

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/metatube-community/metatube-sdk-go/common/fetch"
	"github.com/metatube-community/metatube-sdk-go/translate"
)

var _ translate.Translator = (*XiaoNiu)(nil)

const scTranslateAPI = "https://api.niutrans.com/NiuTransServer/translation"

type XiaoNiu struct {
	ApiKey string `json:"xiaoniu-api-key"`
}

func (xn *XiaoNiu) Translate(text, from, to string) (result string, err error) {
	var (
		resp   *http.Response
		apikey = xn.ApiKey
	)

	if resp, err = fetch.Post(
		scTranslateAPI,
		fetch.WithURLEncodedBody(map[string]string{
			"from":     parseToSupportedLanguage(from),
			"to":       parseToSupportedLanguage(to),
			"apikey":   apikey,
			"src_text": text,
		}),
		fetch.WithHeader("Content-Type", "application/x-www-form-urlencoded"),
	); err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", err
	}

	data := struct {
		From      string `json:"from"`
		To        string `json:"to"`
		Text      string `json:"tgt_text"`
		ErrorCode string `json:"error_code"`
		ErrorMsg  string `json:"error_msg"`
	}{}
	if err = json.NewDecoder(resp.Body).Decode(&data); err == nil {
		if len(data.Text) > 0 {
			result = strings.TrimSpace(data.Text)
		} else {
			err = fmt.Errorf("%s: %s", data.ErrorCode, data.ErrorMsg)
		}
	}
	return
}

func parseToSupportedLanguage(lang string) string {
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
		return "ko"
	case "vie", "vi":
		return "vi"
	case "spa", "es":
		return "es"
	case "fra", "fr":
		return "fr"
	case "ara", "ar":
		return "ar"
	case "bul", "bg":
		return "bg"
	case "est", "et":
		return "et"
	case "dan", "da":
		return "da"
	case "fin", "fi":
		return "fi"
	case "rom", "ro":
		return "ro"
	case "slo", "sl":
		return "sl"
	case "swe", "sv":
		return "sv"
	}
	return lang
}

func init() {
	translate.Register(&XiaoNiu{})
}
