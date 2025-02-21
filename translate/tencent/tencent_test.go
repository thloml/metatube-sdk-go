package tencent

import (
	"os"
	"strconv"
	"testing"
)

func TestTencentTranslate(t *testing.T) {
	for _, unit := range []struct {
		text, from, to string
	}{
		{`hello world`, "auto", "zh-CN"},
	} {
		num, err := strconv.ParseInt(os.Getenv("TENCENT_PROJECT_ID"), 10, 64)
		if err != nil {
			t.Fatal(err)
		}
		result, err := (&Tencent{
			SecretId:  os.Getenv("TENCENT_SECRET_ID"),
			SecretKey: os.Getenv("TENCENT_SECRET_KEY"),
			ProjectId: num,
		}).Translate(unit.text, unit.from, unit.to)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(result)
	}
}
