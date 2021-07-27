package codegen

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var svcType *string

func createServiceCommand() *cobra.Command {
	service := &cobra.Command{
		Use:   "service",
		Short: "generate service and service tests",
		Long:  "generate service and service tests",
		Run: func(cmd *cobra.Command, args []string) {
			tplData := prepareServiceCmdParams()

			outDir, err := ensureDir(*moduleName)
			if err != nil {
				os.Exit(1)
			}

			generate(outDir, "service", fmt.Sprintf("%vService.go", tplData["LName"]), tplData)
			generate(outDir, "service_test", fmt.Sprintf("%vService_test.go", tplData["LName"]), tplData)
		},
	}

	moduleName = service.Flags().String("module", "", "name of module (required)")
	service.MarkFlagRequired("module")

	svcType = service.Flags().String("typeName", "", "name of the type for which a service is generated (required)")
	service.MarkFlagRequired("typeName")

	return service
}

func prepareServiceCmdParams() map[string]interface{} {
	var tplData = make(map[string]interface{})

	tplData["Module"] = *moduleName
	tplData["Name"] = svcType
	tplData["LName"] = strings.ToLower((*svcType)[:1]) + (*svcType)[1:]

	return tplData
}
