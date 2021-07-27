package codegen

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
)

func ensureDir(moduleName string) (string, error) {
	outDir := filepath.Join(os.Getenv("GOITS_HOME"), "codegen", "out", moduleName)
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		err := os.Mkdir(outDir, os.ModeDir)
		if err != nil {
			log.Println("could not create output directory", err)
			return "", err
		}
	}

	return outDir, nil
}

func generate(outDir string, template string, outFile string, data map[string]interface{}) {
	f, err := os.Create(filepath.Join(outDir, outFile))
	if err != nil {
		log.Fatal("could not create output file", err)
	}

	defer f.Close()
	w := bufio.NewWriter(f)
	defer w.Flush()

	log.Printf("generating %s under %s", outFile, outDir)
	t.ExecuteTemplate(w, template+".tmpl", data)
}
