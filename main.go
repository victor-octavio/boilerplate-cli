package main

import (
	"boilerplate-cli/cmd/templates"
	"boilerplate-cli/cmd/ui"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"io"
	"os"
	"strings"
	"time"
)

type step int

const (
	askProjectName step = iota
	askProjectType
	askFramework
	askDatabase
	creatingProject
	finish
	done
)

type item string

func (i item) Title() string       { return string(i) }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return string(i) }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i := listItem.(item)
	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render(">> " + strings.Join(s, " "))
		}
	}
	fmt.Fprint(w, fn(fmt.Sprintf("%d. %s", index+1, i)))
}

type tickMsg time.Time
type doneMsg struct{}

type Model struct {
	step        step
	projectName string
	projectType string
	framework   string
	db          string
	input       string
	list        list.Model
	err         error
	progress    progress.Model
}

var (
	titleStyle        = lipgloss.NewStyle().PaddingLeft(2).PaddingRight(2).Foreground(lipgloss.Color("")).Background(lipgloss.Color("#049488")).Bold(true)
	itemStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#049488"))
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(2)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(2).PaddingBottom(1)
)

func (m Model) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.step {
	case askProjectName:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				m.list = createList([]string{"API (REST/gRPC/GraphQL)", "CLI"}, "Selecione o tipo de projeto")
				m.projectName = strings.TrimSpace(m.input)
				m.input = ""
				m.step = askProjectType
				return m, nil
			case tea.KeyCtrlC, tea.KeyEsc:
				return m, tea.Quit
			case tea.KeyBackspace, tea.KeyDelete:
				if len(m.input) > 0 {
					m.input = m.input[:len(m.input)-1]
				}
			case tea.KeyRunes:
				m.input += msg.String()
			}
		}

	case askProjectType, askFramework, askDatabase:
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)

		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.Type == tea.KeyEnter {
			selected := m.list.SelectedItem().(item)
			switch m.step {
			case askProjectType:
				m.projectType = string(selected)
				if m.projectType == "CLI" {
					m.step = finish
					m.progress = progress.New(progress.WithDefaultGradient(), progress.WithWidth(50))
					return m, tea.Batch(tickCmd(), doProjectCreation(m.projectName, m.projectType, m.framework, m.db))
				} else {
					m.step = askFramework
					m.list = createList([]string{"Gin Gonic", "Echo", "Fiber", "Nenhum"}, "Selecione um framework")
				}
			case askFramework:
				m.framework = string(selected)
				m.step = askDatabase
				m.list = createList([]string{"PostgreSQL", "MySQL", "MongoDB", "Nenhum"}, "Selecione o Banco de Dados desejado")
			case askDatabase:
				m.db = string(selected)
				m.step = finish
				m.progress = progress.New(progress.WithDefaultGradient(), progress.WithWidth(50))
				return m, tea.Batch(tickCmd(), doProjectCreation(m.projectName, m.projectType, m.framework, m.db))
			}
			return m, nil
		}
		return m, cmd

	case finish:
		m.step = creatingProject
		m.progress = progress.New(progress.WithDefaultGradient(), progress.WithWidth(50))
		return m, tea.Batch(tickCmd(), doProjectCreation(m.projectName, m.projectType, m.framework, m.db))

	case creatingProject:
		switch msg := msg.(type) {
		case tickMsg:
			if m.progress.Percent() >= 1.0 {
				m.step = done
				return m, tea.Quit
			}
			cmd := m.progress.IncrPercent(0.25)
			return m, tea.Batch(cmd, tickCmd())

		case progress.FrameMsg:
			pm, cmd := m.progress.Update(msg)
			m.progress = pm.(progress.Model)
			return m, cmd

		case doneMsg:
			return m, nil
		}
	}
	return m, nil
}

func (m Model) View() string {
	switch m.step {
	case askProjectName:
		logo := ui.LogoMessage()
		header := titleStyle.Render("Choose your project name 🚀")
		input := fmt.Sprintf("%s\n> %s", header, m.input)
		return fmt.Sprintf("%s\n\n%s", logo, itemStyle.Render(input))

	case askProjectType, askFramework, askDatabase:
		return m.list.View()

	case creatingProject:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Render("Criando projeto...\n\n") + m.progress.View()

	case done:
		return "\n🎉 Projeto criado com sucesso!\n\nPressione Ctrl+C para sair."

	}
	return ""
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func doProjectCreation(name, typ, fw, db string) tea.Cmd {
	return func() tea.Msg {
		templates.CreateProjectStructure(name, typ, fw, db)
		return doneMsg{}
	}
}

func createList(items []string, title string) list.Model {
	var listItems []list.Item
	for _, s := range items {
		listItems = append(listItems, item(s))
	}

	const height = 10
	const width = 40

	l := list.New(listItems, itemDelegate{}, width, height)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return l
}

func main() {
	if err := tea.NewProgram(Model{}).Start(); err != nil {
		fmt.Println("Erro:", err)
		os.Exit(1)
	}
}
