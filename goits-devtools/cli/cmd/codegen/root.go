package codegen

import (
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

var CodegenCommand *cobra.Command

var moduleName *string

var t *template.Template

func init() {
	initTemplates()

	CodegenCommand = &cobra.Command{
		Use:     "codegen",
		Aliases: []string{"gen"},
		Short:   "generate code from templates",
		Long:    "generate boiler plate code for application components",
	}

	CodegenCommand.AddCommand(createDomainCommand())
	CodegenCommand.AddCommand(createServiceCommand())
}

func initTemplates() {
	tplRootDir := filepath.Join(os.Getenv("GOITS_HOME"), "codegen", "template")
	t = template.Must(template.New("codegen").ParseGlob(tplRootDir + "/*"))
}
