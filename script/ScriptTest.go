package script

import (
	"fmt"
	"strconv"
)

type ScriptTest struct{}

func (o ScriptTest) GetIntTest(bo map[string]interface{}, data map[string]interface{}) string {
	masterData := bo["A"].(map[string]interface{})
	attachCount, err := strconv.Atoi(fmt.Sprint(masterData["attachCount"]))
	if err != nil {
		panic(err)
	}
	return fmt.Sprint(attachCount*20 + 30)
}
