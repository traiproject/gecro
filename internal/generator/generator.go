// Copyright (c) TRAI
// SPDX-License-Identifier: MIT

package generator

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"go.trai.ch/gecro/config"
)

//go:embed all:template
var templateFiles embed.FS

// ANSI escape code for green text
const green = "\033[32m"
const reset = "\033[0m"

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
	t := template.New("gecro")

	// Walk the "template" directory in the embedded FS
	err := fs.WalkDir(templateFiles, "template", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if d.IsDir() {
			return nil
		}

		// Read the file content
		content, err := templateFiles.ReadFile(path)
		if err != nil {
			return err
		}

		// Create a new template with the full path as its name and parse the content.
		// This makes the template name unique.
		_, err = t.New(path).Parse(string(content))
		return err
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &Generator{
		templateFS: templateFiles,
		templates:  t,
	}, nil
}

func printCreated(path string) {
	fmt.Printf("%sCREATED%s %s\n", green, reset, path)
}

func printDryRun(path string) {
	fmt.Printf("%sDRY RUN%s %s\n", green, reset, path)
}

func (g *Generator) GenerateConfig(force bool, cfg config.Config) error {
	outputPath := "gecro.yaml"

	if !force {
		if _, err := os.Stat(outputPath); err == nil {
			return fmt.Errorf("'%s' already exists. Use --force to overwrite", outputPath)
		}
	}

	err := g.generateFileFromTemplate("template/gecro.yaml.tmpl", outputPath, cfg)
	if err != nil {
		return fmt.Errorf("failed to generate config file: %w", err)
	}

	return nil
}

func (g *Generator) GeneratePkg(config config.Config) error {
	if config.DryRun {
		fmt.Println("--- Performing a dry run for a new package. No files will be written. ---")
	}
	fmt.Printf("Generating package '%s' at '%s/pkg/%s'\n", config.Name, config.OutputDir, config.Name)

	if config.Name == "" {
		return fmt.Errorf("package name cannot be empty")
	}

	specs := []DirectorySpec{
		{
			RelativePath: "pkg/{name}", // The generator will replace {name}
			Files: []FileSpec{
				{TemplateName: "template/pkg/go.mod.tmpl", FileName: "go.mod"},
				{TemplateName: "template/pkg/pkg.go.tmpl", FileName: config.Name + ".go"},
			},
		},
	}

	if err := g.createDirectoriesAndFiles(specs, config); err != nil {
		return fmt.Errorf("error generating shared package: %w", err)
	}

	return nil
}

func (g *Generator) GenerateService(config config.Config) error {
	if config.DryRun {
		fmt.Println("--- Performing a dry run. No files will be written. ---")
	}
	fmt.Printf("Generating service '%s' at '%s'\n", config.Name, config.OutputDir)

	// 1. Perform detailed input validation (domain-specific)
	if config.Name == "" {
		return fmt.Errorf("service name cannot be empty")
	}
	// Add more validation for ModulePath, etc.

	specs := []DirectorySpec{
		{
			RelativePath: "cmd/{name}",
			Files: []FileSpec{
				{TemplateName: "template/service/cmd/servicename/main.go.tmpl", FileName: "main.go"},
				{TemplateName: "template/service/cmd/servicename/wire.go.tmpl", FileName: "wire.go"},
				{TemplateName: "template/service/cmd/servicename/go.mod.tmpl", FileName: "go.mod"},
			},
		},
		{
			RelativePath: "proto/config/{name}/v1",
			Files: []FileSpec{
				{TemplateName: "template/service/proto/config/servicename/v1/config.proto.tmpl", FileName: "config.proto"},
			},
		},
		{
			RelativePath: "configs/{name}/v1",
			Files: []FileSpec{
				{TemplateName: "template/service/configs/servicename/v1/config.yaml.tmpl", FileName: "config.yaml"},
			},
			BuildFileContent: `filegroup(
	name = "config",
	srcs = ["config.yaml"],
	visibility = ["//visibility:public"],
)`,
		},
		{
			RelativePath: "internal/{name}/biz",
			Files: []FileSpec{
				{TemplateName: "template/service/internal/servicename/biz/biz.go.tmpl", FileName: "biz.go"},
			},
		},
		{
			RelativePath: "internal/{name}/data",
			Files: []FileSpec{
				{TemplateName: "template/service/internal/servicename/data/data.go.tmpl", FileName: "data.go"},
			},
		},
		{
			RelativePath: "internal/{name}/server",
			Files: []FileSpec{
				{TemplateName: "template/service/internal/servicename/server/grpc.go.tmpl", FileName: "grpc.go"},
				{TemplateName: "template/service/internal/servicename/server/http.go.tmpl", FileName: "http.go"},
				{TemplateName: "template/service/internal/servicename/server/server.go.tmpl", FileName: "server.go"},
			},
		},
		{
			RelativePath: "internal/{name}/service",
			Files: []FileSpec{
				{TemplateName: "template/service/internal/servicename/service/service.go.tmpl", FileName: "service.go"},
			},
		},
	}

	if err := g.createDirectoriesAndFiles(specs, config); err != nil {
		return fmt.Errorf("Error generating microservice: %w", err)
	}

	return nil
}

func (g *Generator) createDirectoriesAndFiles(specs []DirectorySpec, config config.Config) error {
	for _, spec := range specs {
		// Replace the {name} placeholder with the actual service name
		replacedPath := strings.ReplaceAll(spec.RelativePath, "{name}", config.Name)

		// Construct the full directory path
		fullDirPath := filepath.Join(config.OutputDir, replacedPath)

		if !config.DryRun {
			// Create the directory structure
			if err := os.MkdirAll(fullDirPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory '%s': %w", fullDirPath, err)
			}
		}

		// Generate each file within the directory
		for _, file := range spec.Files {
			outputPath := filepath.Join(fullDirPath, file.FileName)
			if err := g.generateFileFromTemplate(file.TemplateName, outputPath, config); err != nil {
				return err
			}
		}

		// Create an empty BUILD.bazel file in the directory
		buildFilePath := filepath.Join(fullDirPath, "BUILD.bazel")
		if spec.BuildFileContent != "" {
			if err := g.createBuildFileWithContent(buildFilePath, spec.BuildFileContent, config); err != nil {
				return err
			}
		} else {
			if err := g.createEmptyBuildFile(buildFilePath, config); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *Generator) generateFileFromTemplate(templateName, outputPath string, config config.Config) error {
	if _, err := os.Stat(outputPath); err == nil {
		// File exists, no need to create
		return nil
	} else if !os.IsNotExist(err) {
		// An error other than non-existence occurred
		return fmt.Errorf("failed to check existence of '%s': %w", outputPath, err)
	}

	if config.DryRun {
		// Execute the template into a discard buffer to test it
		if err := g.templates.ExecuteTemplate(io.Discard, templateName, config); err != nil {
			return fmt.Errorf("failed to execute template '%s' for '%s' during dry run: %w", templateName, outputPath, err)
		}
		printDryRun(outputPath)
		return nil
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file '%s': %w", outputPath, err)
	}
	defer file.Close()

	if err := g.templates.ExecuteTemplate(file, templateName, config); err != nil {
		return fmt.Errorf("failed to execute template '%s' for '%s': %w", templateName, outputPath, err)
	}
	printCreated(outputPath)
	return nil
}

func (g *Generator) createBuildFileWithContent(outputPath, content string, config config.Config) error {
	if _, err := os.Stat(outputPath); err == nil {
		// File exists, no need to create
		return nil
	} else if !os.IsNotExist(err) {
		// An error other than non-existence occurred
		return fmt.Errorf("failed to check existence of '%s': %w", outputPath, err)
	}

	if config.DryRun {
		printDryRun(outputPath)
		return nil
	}

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

func (g *Generator) createEmptyBuildFile(outputPath string, config config.Config) error {
	if _, err := os.Stat(outputPath); err == nil {
		// File exists, no need to create
		return nil
	} else if !os.IsNotExist(err) {
		// An error other than non-existence occurred
		return fmt.Errorf("failed to check existence of '%s': %w", outputPath, err)
	}

	if config.DryRun {
		printDryRun(outputPath)
		return nil
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create BUILD.bazel file at '%s': %w", outputPath, err)
	}
	defer file.Close()

	printCreated(outputPath)
	return nil
}
