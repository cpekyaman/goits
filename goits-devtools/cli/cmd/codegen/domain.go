package codegen

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var types *[]string

func createDomainCommand() *cobra.Command {
	domain := &cobra.Command{
		Use:   "domain",
		Short: "generate domain entities and repositories",
		Long:  "generate domain entities and repositories",
		Run: func(cmd *cobra.Command, args []string) {
			var tplData = prepareDomainCmdParams()

			outDir, err := ensureDir(*moduleName)
			if err != nil {
				os.Exit(1)
			}

			generate(outDir, "domain", "domain.go", tplData)
			generate(outDir, "repository", "repository.go", tplData)
		},
	}

	moduleName = domain.Flags().String("module", "", "name of module (required)")
	domain.MarkFlagRequired("module")

	types = domain.Flags().StringSlice("types", []string{}, "comma separated list of types in type:base form (required)")
	domain.MarkFlagRequired("types")

	return domain
}

func prepareDomainCmdParams() map[string]interface{} {
	var tplData = make(map[string]interface{})
	tplData["Module"] = *moduleName

	tDataList := make([]map[string]string, len(*types))
	for i, v := range *types {
		typeAndBase := strings.Split(v, ":")
		tData := map[string]string{
			"Name":  typeAndBase[0],
			"Base":  typeAndBase[1],
			"LName": strings.ToLower(typeAndBase[0][:1]) + typeAndBase[0][1:],
		}
		tDataList[i] = tData
	}
	tplData["Types"] = tDataList
	return tplData
}
