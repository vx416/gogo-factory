package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/vx416/gogo-factory/codegen"
)

var (
	input   = flag.String("i", ".", "Input file or directory. (required)")
	structs = flag.String("s", "", "The needed struct name, these names should  be concat by comma. (required)")
	output  = flag.String("o", "", "Output directory for generated factory code. (optional)")
	p       = flag.Bool("p", false, "Print the result. (optional)")
)

func main() {
	flag.Parse()

	inputInfo, err := os.Stat(*input)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("input(%s) not exists\n", *input)
			os.Exit(1)
		}
		fmt.Printf("input(%s) invalid\n", *input)
		os.Exit(1)
	}
	if inputInfo.IsDir() {
		filepath.Walk(*input, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				parseFile(path)
			}
			return nil
		})

		return
	}
	parseFile(*input)
}

func parseFile(filePath string) {
	if filepath.Ext(filePath) != ".go" {
		fmt.Printf("input(%s) should be go file\n", filePath)
		os.Exit(1)
	}

	node, err := codegen.ParseFile(filePath)
	if err != nil {
		fmt.Printf("parse input content failed, err:%+v", err)
		os.Exit(1)
	}
	if *structs == "" {
		fmt.Println("struct name cannot be empty")
		os.Exit(1)
	}
	modelsName := strings.Split(*structs, ",")
	fileMeta := codegen.ParseFileMeta(node, modelsName...)

	t, err := codegen.GetTempalte(false, fileMeta)
	if err != nil {
		fmt.Printf("get factory content failed, err:%+v", err)
		os.Exit(1)
	}
	if *p {
		fmt.Println(filePath)
		fmt.Println(t)
	}
	if *output != "" {
		path := filepath.Base(filePath)
		fiileName := strings.Split(path, ".")[0]
		filepath.Join(*output, fiileName)
		outputFile, err := os.Create(filepath.Join(*output, fiileName+"_factory.go"))
		if err != nil {
			fmt.Printf("create output(%s) file failed, err:%+v", *output, err)
			os.Exit(1)
		}
		outputFile.WriteString(t)
	}
}

var (
	usage = `Usage: gofactory [OPTIONS]
	-i <input_path>
	-m <models>
	-o <output_path>
	-p <only_print>
	`
)
