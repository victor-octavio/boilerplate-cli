package templates

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func CreateProjectStructure(projectName, projectType, framework, db string) {
	baseDirs := []string{}

	switch projectType {
	case "API (REST/gRPC/GraphQL)":
		baseDirs = []string{
			"cmd",
			"cmd/app",
			"internal/config",
			"pkg",
			"pkg/db",
			"api",
		}
	case "CLI":
		baseDirs = []string{
			"cmd",
			"internal/config",
			"pkg",
		}
	}

	// Criar diretórios
	for _, dir := range baseDirs {
		path := filepath.Join(projectName, dir)
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", path, err)
			log.Fatal(err)
			return
		}
	}

	// Conteúdo do main.go
	mainContent := FrameworkInit(framework)
	mainContent = strings.ReplaceAll(mainContent, "{{.ProjectName}}", projectName)

	// go.mod
	goModContent := fmt.Sprintf("module %s\n\ngo 1.20", projectName)

	switch framework {
	case "Gin Gonic":
		goModContent += "\nrequire github.com/gin-gonic/gin latest"
	case "Echo":
		goModContent += "\nrequire github.com/labstack/echo/v4 latest"
	case "Fiber":
		goModContent += "\nrequire github.com/gofiber/fiber/v2 latest"
	}

	switch db {
	case "PostgreSQL":
		goModContent += "\nrequire github.com/lib/pq latest"
	case "MySQL":
		goModContent += "\nrequire github.com/go-sql-driver/mysql latest"
	case "MongoDB":
		goModContent += "\nrequire go.mongodb.org/mongo-driver latest"
	}

	readmeContent := fmt.Sprintf(`# %s
Awesome project generated with boilerplate-cli! 🚀
`, projectName)

	// Criação de arquivos
	files := map[string]string{
		"main.go":   mainContent,
		"go.mod":    goModContent,
		"README.md": readmeContent,
	}

	for fileName, content := range files {
		var path string
		if fileName == "main.go" {
			path = filepath.Join(projectName, "cmd", "app", fileName)
		} else {
			path = filepath.Join(projectName, fileName)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			fmt.Printf("Error writing file %s: %v\n", path, err)
			return
		}
	}

	if db != "none" {
		createDBFiles(projectName, db)
	}
	return
}
