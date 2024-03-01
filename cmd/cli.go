package cmd

import (
	"fmt"
	"log"
	"os"
	"shedock/internal/core"
	"shedock/internal/instance"
	"shedock/pkg/docker/file"
	"shedock/pkg/parsers/shellscript"
	"shedock/pkg/shell"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mitchellh/go-wordwrap"
	"github.com/spf13/cobra"
)

var (
	spinners = []spinner.Spinner{
		// spinner.MiniDot
		spinner.Points,
		{
			Frames: []string{"🧡", "💛", "💚", "💙", "💜", "🖤", "🤍"},
			FPS:    time.Second / 19,
		},
	}

	textStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Bold(true).
			Render
	bodyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("251")).
			Render
	textStyleInsights = lipgloss.NewStyle().
				Foreground(lipgloss.Color("45")).
				Bold(true).PaddingLeft(2).
				PaddingRight(2).
				Border(lipgloss.RoundedBorder()).
				Render
	commandStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Background(lipgloss.Color("12")).PaddingLeft(1).PaddingRight(1).Render
	shellNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("45")).Render
	spinnerStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
)

type model struct {
	step          int
	spinner       spinner.Model
	commandOutput string
	builder       *core.ImageBuilder
}
type analyzeShellScriptMsg string
type getTransitiveDependenciesMsg string
type generateDockerfileMsg string
type buildImageMsg string
type insightsMsg string

func (m *model) analyzeShellScriptCmd() tea.Cmd {
	return func() tea.Msg {
		// we pull a sneaky on ya
		instance.Init()
		err := m.builder.LoadSystemBuiltins()
		if err != nil {
			return err
		}

		shell, _ := m.builder.Script.GetShell()
		err = m.builder.LoadScriptDeps()
		if err != nil {
			return err
		}

		scriptDeps := m.builder.GetScriptDeps()
		if err != nil {
			return err
		}

		err = m.builder.UsedSystemBuiltins()
		if err != nil {
			return err
		}

		totalDependencies := len(scriptDeps)

		return analyzeShellScriptMsg(
			bodyStyle(
				fmt.Sprintf("\n✅ %s\n├── Shell recognized: %s\n└── Found %d root dependencies\n", textStyle("Analyzing shell script"), shellNameStyle(shell), totalDependencies),
			),
		)
	}
}

func (m *model) getTransitiveDependenciesCmd() tea.Cmd {
	return func() tea.Msg {
		err := m.builder.LoadShellBuiltins()
		if err != nil {
			return err
		}
		// remove not-supported commands from script deps
		// remove shell-builtins and system-builtins from script deps and find what we can get from package manager
		filteredDeps := m.builder.FilterCmdsToInstall()
		err = m.builder.DependenciesAvailableOnPackageHost(filteredDeps)
		if err != nil {
			return err
		}
		err = m.builder.LoadAllSharedLibs()
		if err != nil {
			return err
		}

		cmdOnApk := m.builder.GetCmdOnApk()
		// commands not available on apk come under not found
		var notFound []string
		for _, dep := range filteredDeps {
			var found bool
			for _, cmd := range cmdOnApk {
				if dep == cmd.Name {
					found = true
					break
				}
			}
			if !found {
				notFound = append(notFound, dep)
			}
		}
		m.builder.UpdateCmdsNotFound(notFound)

		libs := m.builder.GetSharedLibs()

		transitiveDependencies := len(libs)
		return getTransitiveDependenciesMsg(
			bodyStyle(
				fmt.Sprintf("\n✅ %s\n└── Found %d transitive dependencies\n", textStyle("Getting transitive dependencies"), transitiveDependencies),
			),
		)
	}
}

func (m *model) generateDockerfileCmd() tea.Cmd {
	return func() tea.Msg {
		var deps file.Dependencies
		var bins []file.Dependency
		var libs []file.Dependency

		shell, _ := m.builder.Script.GetShell()
		systemBuiltins := m.builder.GetUsedSystemBuiltins()
		externalCommands := m.builder.GetCmdOnApk()
		sharedLibs := m.builder.GetSharedLibs()

		for _, cmd := range systemBuiltins {
			bins = append(bins, file.Dependency{
				FromPath: cmd.Path,
				ToPath:   cmd.Path,
			})
		}

		for _, cmd := range externalCommands {
			bins = append(bins, file.Dependency{
				FromPath: cmd.Path,
				ToPath:   cmd.Path,
			})
		}

		for _, lib := range sharedLibs {
			libs = append(libs, file.Dependency{
				FromPath:   lib.FullPath,
				ToPath:     lib.FullPath,
				Requiredby: lib.DependencyOf,
			})
		}

		deps = file.Dependencies{
			Bin: bins,
			Lib: libs,
		}

		file := &file.Dockerfile{
			DependenciesToInstall: externalCommands,
			Script:                m.builder.Script.ScriptPath,
			ShellPath:             shell,
			Dependencies:          deps,
		}

		dockerFilePath, err := file.Generate()
		if err != nil {
			return err
		}

		return generateDockerfileMsg(
			bodyStyle(
				fmt.Sprintf("\n✅ %s\n└── Dockerfile generated at %s\n", textStyle("Generating Dockerfile"), dockerFilePath),
			),
		)
	}
}

func (m *model) buildImageCmd() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(1 * time.Second)
		imageName := ""
		return buildImageMsg(
			bodyStyle(
				fmt.Sprintf("\n✅ %s\n├── Image built successfully!\n└── Run the following command to start a container for your image:\n    > docker run -it --rm %s:latest", textStyle("Building docker image"), imageName),
			),
		)
	}
}

func (m *model) getInsightsCmd() tea.Cmd {
	return func() tea.Msg {

		not_found := m.builder.GetCmdNotOnApk()
		not_supported := m.builder.GetCmdNotSupported()
		if len(not_found) > 0 || len(not_supported) > 0 {
			var insights string

			title := fmt.Sprintf("\n\n%s\n", textStyleInsights("We are not perfect!"))
			if len(not_supported) > 0 {
				insights += "- We have recognized some dependencies that cannot work in a containerized environment. Consider removing them from your script or adding workarounds for them:\n"
				for _, cmd := range not_supported {
					insights += fmt.Sprintf("  - %s\n", commandStyle(cmd))
				}
			}
			if len(not_found) > 0 {
				insights += "- We couldn't find the following dependencies. Consider installing them manually:\n"
				for _, cmd := range not_found {
					insights += fmt.Sprintf("  - %s\n", commandStyle(cmd))
				}
			}

			insights += fmt.Sprintf("\n\n%s\n\n", "Report any issues at link")
			wrappedText := wordwrap.WrapString(insights, 80) // Wrap text after 80 characters
			return insightsMsg(title + wrappedText)
		}
		return insightsMsg("")
	}
}

func (m *model) Init() tea.Cmd {
	return tea.Sequence(
		m.spinner.Tick,
		m.analyzeShellScriptCmd(),
		m.getTransitiveDependenciesCmd(),
		m.generateDockerfileCmd(),
		// m.buildImageCmd(),
		m.getInsightsCmd(),
	)
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case analyzeShellScriptMsg:
		m.commandOutput += string(msg)
		m.step++
	case getTransitiveDependenciesMsg:
		m.commandOutput += string(msg)
		m.step++
	case generateDockerfileMsg:
		m.commandOutput += string(msg)
		m.step++
	// case buildImageMsg:
	// 	m.commandOutput += string(msg)
	// 	m.step++
	case insightsMsg:
		m.commandOutput += string(msg)
		m.step++
	default:
		return m, nil
	}
	if m.step > 3 {
		return m, tea.Quit
	}
	return m, nil
}

func (m *model) initSpinner() {
	m.spinner = spinner.New()
	m.spinner.Style = spinnerStyle
	m.spinner.Spinner = spinners[0]
}

func (m *model) View() string {
	var s string
	s += m.commandOutput
	switch m.step {
	case 0:
		s += fmt.Sprintf("\n%s %s", m.spinner.View(), textStyle("Analyzing shell script"))
	case 1:
		s += fmt.Sprintf("\n\n%s %s", m.spinner.View(), textStyle("Getting transitive dependencies"))
	case 2:
		s += fmt.Sprintf("\n\n%s %s", m.spinner.View(), textStyle("Generating Dockerfile"))
		// case 3:
		// 	s += fmt.Sprintf("\n\n%s %s", m.spinner.View(), textStyle("Building image"))
	}
	return s
}

func DriverCli(cmd *cobra.Command, args []string) {
	// check if filepath is a shell script
	script := shellscript.Script{ScriptPath: args[0]}
	isScript, err := script.Validate()
	if err != nil || !isScript {
		log.Fatalf("Failed to validate script: %v", err)
	}

	shellType, err := script.GetShell()
	if err != nil {
		log.Fatalf("Failed to get shell type: %v", err)
	}
	shell, err := shell.NewShell(shellType)
	if err != nil {
		log.Fatalf("Failed to get shell: %v", err)
	}

	imageBuilder := core.NewImageBuilder(
		&script,
		shell,
	)

	m := &model{
		builder: imageBuilder,
	}
	m.initSpinner()

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}

	defer instance.Destroy()
}
