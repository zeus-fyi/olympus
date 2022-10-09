package drivers

import (
	code_driver "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/drivers/code"
	template_driver "github.com/zeus-fyi/olympus/pkg/hera/lib/v0/drivers/template"
)

type DriverLib struct {
	code_driver.CodeDriverLib
	template_driver.TemplateDriverLib
}
