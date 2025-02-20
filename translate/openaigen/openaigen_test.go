package openaigen

import (
	"os"
	"testing"
)

func TestOpenaiGenTranslate(t *testing.T) {
	for _, unit := range []struct {
		text, from, to string
	}{
		{`Oh yeah! I'm a translator!`, "", "zh-CN"},
	} {
		result, err := (&OpenAIGen{
			Url:   os.Getenv("OPENAIGEN_URL"),
			Model: os.Getenv("OPENAIGEN_MODEL"),
			Auth:  os.Getenv("OPENAIGEN_AUTH"),
		}).Translate(unit.text, unit.from, unit.to)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(result)
	}
}
