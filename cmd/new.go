package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Generate a new project using Package Oriented Design for Go",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		interactiveProjectSetup(projectName)
	},
}

func interactiveProjectSetup(projectName string) {
	generateFinalMessage()
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(
		"\n1. API (REST, gRPC, GraphQL) 🍃\n" +
			"\n2. CLI (Cobra) 🐍\n\n" +
			"\nWhich kind of project are you creating?  ")

	projectType, _ := reader.ReadString('\n')
	projectType = strings.TrimSpace(projectType)

	switch projectType {
	case "1":
		fmt.Print("\nWould you like to use some framework? (gin/echo/fiber/none): ")
		framework, _ := reader.ReadString('\n')
		framework = strings.TrimSpace(framework)
		fmt.Print("\nDatabase (postgres/mysql/mongo/none): ")
		db, _ := reader.ReadString('\n')
		db = strings.TrimSpace(db)
		createProjectStructure(projectName, projectType, framework, db)
		break
	case "2":
		createProjectStructure(projectName, projectType, "", "")
		break
	default:
		interactiveProjectSetup(projectName)
		break
	}

}

func createProjectStructure(projectName, projectType, framework, db string) {
	baseDirs := []string{}

	switch projectType {
	case "1":
		baseDirs = []string{
			"cmd",
			"cmd/app",
			"internal/config",
			"pkg",
			"pkg/db",
			"api",
		}
	case "2":
		baseDirs = []string{
			"cmd",
			"internal/config",
			"pkg",
		}
	}

	for _, dir := range baseDirs {
		path := filepath.Join(projectName, dir)
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			fmt.Printf("Error creating project packages. 😔 %s: %v\n", path, err)
			return
		}
	}

	mainContent := FrameworkInit(framework)
	mainContent = strings.ReplaceAll(mainContent, "{{.ProjectName}}", projectName)

	goModContent := fmt.Sprintf("module %s\n\ngo 1.20", projectName)

	if framework == "gin" {
		goModContent += "\nrequire github.com/gin-gonic/gin latest"
	} else if framework == "echo" {
		goModContent += "\nrequire github.com/labstack/echo/v4 latest"
	} else if framework == "fiber" {
		goModContent += "\nrequire github.com/gofiber/fiber/v2 latest"
	}

	switch db {
	case "postgres":
		goModContent += "\nrequire github.com/lib/pq latest"
	case "mysql":
		goModContent += "\nrequire github.com/go-sql-driver/mysql latest"
	case "mongo":
		goModContent += "\nrequire go.mongodb.org/mongo-driver latest"
	}

	readmeContent := fmt.Sprintf(`
# %s

Awesome project generated with boilerplate-cli! 🚀🚀🚀
	
`, projectName)

	files := map[string]string{
		"main.go":   mainContent,
		"go.mod":    goModContent,
		"README.md": readmeContent,
	}

	for file, content := range files {
		if file == "main.go" {
			path := filepath.Join(projectName, "cmd", "app", file)
			if err := os.WriteFile(path, []byte(content), 0644); err != nil {
				fmt.Printf("Error writing main.go file: %v\n", err)
			}
			continue
		}
		path := filepath.Join(projectName, file)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			fmt.Printf("Error creating file %s: %v\n", path, err)
		}
	}

	if db != "none" {
		createDBFiles(projectName, db)
	}

	fmt.Println("\nProject created successfully! 😮‍💨")
}

func init() {
	rootCmd.AddCommand(newCmd)
}
