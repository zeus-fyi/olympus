package code_driver

import (
	"bytes"
	"fmt"

	"github.com/zeus-fyi/jennifer/jen"
	"github.com/zeus-fyi/olympus/pkg/utils/printer/v0/structs"
)

func AutoGenCodeFromTemplate(codeGenPathInfo structs.Path) {
	f := jen.NewFile(codeGenPathInfo.PackageName)
	f.Func().Id("main").Params().Block()
	buf := &bytes.Buffer{}
	err := f.Render(buf)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(buf.String())
	}

	f.Save("m.go")
}
