package xiaoniu

import (
	"os"
	"testing"
)

func TestXiaoNiuTranslate(t *testing.T) {
	for _, unit := range []struct {
		text, from, to string
	}{
		{`Hello world。`, "auto", "zh-CN"},
	} {
		result, err := (&XiaoNiu{
			ApiKey: os.Getenv("XIAONIU_API_KEY"),
		}).Translate(unit.text, unit.from, unit.to)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(result)
	}
}
