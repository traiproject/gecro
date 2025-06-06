// Copyright (c) TRAI
// SPDX-License-Identifier: MIT

package generator

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed all:template
var templateFiles embed.FS

// ANSI escape code for green text
const green = "\033[32m"
const reset = "\033[0m"

// ServiceParams encapsulates all parameters needed to generate a service.
type ServiceParams struct {
	ServiceName           string
	MonorepoPrefix        string
	OutputDir             string
	KratosVersion         string
	WireVersion           string
	GoVersion             string
	DefaultConfigPath     string
	ServiceNamePascalCase string
	ServiceNameLowerCase  string
}

type FileSpec struct {
	TemplateName string
	FileName     string
}

type DirectorySpec struct {
	RelativePath     string
	Files            []FileSpec
	BuildFileContent string
}

type Generator struct {
	templateFS embed.FS
	templates  *template.Template
}

func NewGenerator() (*Generator, error) {
	tmpl, err := template.ParseFS(templateFiles,
		"template/*/*/*.tmpl",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &Generator{
		templateFS: templateFiles,
		templates:  tmpl,
	}, nil
}

func (g *Generator) GenerateService(params ServiceParams) error {
	fmt.Printf("Generating service '%s' at '%s'\n", params.ServiceName, params.OutputDir)

	// 1. Perform detailed input validation (domain-specific)
	if params.ServiceName == "" {
		return fmt.Errorf("service name cannot be empty")
	}
	// Add more validation for ModulePath, etc.

	specs := []DirectorySpec{
		{
			RelativePath: "cmd/{service}",
			Files: []FileSpec{
				{TemplateName: "main.go.tmpl", FileName: "main.go"},
				{TemplateName: "wire.go.tmpl", FileName: "wire.go"},
				{TemplateName: "go.mod.tmpl", FileName: "go.mod"},
			},
		},
		// {
		// 	RelativePath: "internal/{service}/biz",
		// 	Files: []FileSpec{
		// 		{TemplateName: "biz.go.tmpl", FileName: "biz.go"},
		// 	},
		// },
		// {
		// 	RelativePath: "internal/{service}/data",
		// 	Files: []FileSpec{
		// 		{TemplateName: "data.go.tmpl", FileName: "data.go"},
		// 	},
		// },
		// {
		// 	RelativePath: "internal/{service}/server",
		// 	Files: []FileSpec{
		// 		{TemplateName: "grpc.go.tmpl", FileName: "grpc.go"},
		// 		{TemplateName: "http.go.tmpl", FileName: "http.go"},
		// 		{TemplateName: "server.go.tmpl", FileName: "server.go"},
		// 	},
		// },
		// {
		// 	RelativePath: "internal/{service}/service",
		// 	Files: []FileSpec{
		// 		{TemplateName: "service.go.tmpl", FileName: "service.go"},
		// 	},
		// },
		// {
		// 	RelativePath: "proto/config/{service}/v1",
		// 	Files: []FileSpec{
		// 		{TemplateName: "config.proto.tmpl", FileName: "config.proto"},
		// 	},
		// },
		{
			RelativePath: "configs/{service}/v1",
			Files: []FileSpec{
				{TemplateName: "config.yaml.tmpl", FileName: "config.yaml"},
			},
			BuildFileContent: `filegroup(
			name = "config",
			srcs = ["config.yaml"],
			visibility = ["//visibility:public"],
			)`,
		},
	}

	if err := g.createDirectoriesAndFiles(specs, params); err != nil {
		return fmt.Errorf("Error generating microservice: %w", err)
	}

	return nil
}

func printCreated(path string) {
	fmt.Printf("%sCREATED%s %s\n", green, reset, path)
}

func (g *Generator) createDirectoriesAndFiles(specs []DirectorySpec, params ServiceParams) error {
	for _, spec := range specs {
		// Replace the {service} placeholder with the actual service name
		replacedPath := strings.ReplaceAll(spec.RelativePath, "{service}", params.ServiceName)

		// Construct the full directory path
		fullDirPath := filepath.Join(params.OutputDir, replacedPath)

		// Create the directory structure
		if err := os.MkdirAll(fullDirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory '%s': %w", fullDirPath, err)
		}

		// Generate each file within the directory
		for _, file := range spec.Files {
			outputPath := filepath.Join(fullDirPath, file.FileName)
			if err := g.generateFileFromTemplate(file.TemplateName, outputPath, params); err != nil {
				return err
			}
		}

		// Create an empty BUILD.bazel file in the directory
		buildFilePath := filepath.Join(fullDirPath, "BUILD.bazel")
		if spec.BuildFileContent != "" {
			if err := g.createBuildFileWithContent(buildFilePath, spec.BuildFileContent); err != nil {
				return err
			}
		} else {
			if err := g.createEmptyBuildFile(buildFilePath); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *Generator) generateFileFromTemplate(templateName, outputPath string, params ServiceParams) error {
	if _, err := os.Stat(outputPath); err == nil {
		// File exists, no need to create
		return nil
	} else if !os.IsNotExist(err) {
		// An error other than non-existence occurred
		return fmt.Errorf("failed to check existence of '%s': %w", outputPath, err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file '%s': %w", outputPath, err)
	}
	defer file.Close()

	if err := g.templates.ExecuteTemplate(file, templateName, params); err != nil {
		return fmt.Errorf("failed to execute template '%s' for '%s': %w", templateName, outputPath, err)
	}
	printCreated(outputPath)
	return nil
}

func (g *Generator) createBuildFileWithContent(outputPath, content string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create BUILD.bazel at '%s': %w", outputPath, err)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return fmt.Errorf("failed to write to BUILD.bazel at '%s': %w", outputPath, err)
	}

	printCreated(outputPath)
	return nil
}

func (g *Generator) createEmptyBuildFile(outputPath string) error {
	if _, err := os.Stat(outputPath); err == nil {
		// File exists, no need to create
		return nil
	} else if !os.IsNotExist(err) {
		// An error other than non-existence occurred
		return fmt.Errorf("failed to check existence of '%s': %w", outputPath, err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create BUILD.bazel file at '%s': %w", outputPath, err)
	}
	defer file.Close()

	printCreated(outputPath)
	return nil
}
