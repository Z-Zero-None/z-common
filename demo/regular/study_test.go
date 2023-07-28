package regular

import (
	"regexp"
	"testing"
)

func TestRegx(t *testing.T) {
	content := "For the basic usage introduction we will be installing pendulum, a datetime library. If you have not yet installed Poetry, refer to the Introduction chapter."
	// 匹配符合[...]的字符 result:"y. I"
	condition := "[^A-Z]{2}[\\s]+[IPC]"
	compile := regexp.MustCompile(condition)
	find := compile.FindString(content)
	t.Logf("condition:%s,find:%s\n", condition, find)
}
