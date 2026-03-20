package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/volodymyrkovtun/confradar/internal/config"
	"github.com/volodymyrkovtun/confradar/internal/differ"
	"github.com/volodymyrkovtun/confradar/internal/health"
	"github.com/volodymyrkovtun/confradar/internal/scanner"
)

// Focus targets for Tab cycling.
type focus int

const (
	focusTree focus = iota
	focusViewer
	focusHealth
)

// App is the top-level Bubble Tea model.
type App struct {
	cfg         *config.Config
	rootPath    string
	result      *scanner.ScanResult
	issues      []health.HealthIssue

	tree        TreeView
	viewer      FileViewer
	healthPanel HealthPanel
	diffView    DiffView
	search      SearchOverlay
	help        HelpOverlay
	statusBar   StatusBar
	styles      Styles
	theme       Theme

	focus       focus
	width       int
	height      int
	themeName   string
	showPreview bool
	diffSelect  *scanner.ConfigFile // first file selected for diff

	err error
}

// scanMsg carries completed scan results.
type scanMsg struct {
	result *scanner.ScanResult
	issues []health.HealthIssue
	err    error
}

// NewApp creates the TUI application model.
func NewApp(rootPath string, cfg *config.Config) App {
	themeName := cfg.Display.Theme
	if themeName == "" || themeName == "auto" {
		themeName = "dark"
	}
	theme := GetTheme(themeName)
	styles := NewStyles(theme)

	return App{
		cfg:         cfg,
		rootPath:    rootPath,
		viewer:      NewFileViewer(),
		healthPanel: NewHealthPanel(),
		diffView:    NewDiffView(),
		statusBar:   NewStatusBar(),
		styles:      styles,
		theme:       theme,
		themeName:   themeName,
		showPreview: config.BoolVal(cfg.Display.ShowPreview, true),
		focus:       focusTree,
	}
}

// Init runs the initial scan.
func (a App) Init() tea.Cmd {
	return func() tea.Msg {
		absPath, err := filepath.Abs(a.rootPath)
		if err != nil {
			return scanMsg{err: err}
		}

		result, err := scanner.Scan(a.rootPath, a.cfg)
		if err != nil {
			return scanMsg{err: err}
		}

		issues := health.RunChecks(absPath, result, a.cfg.HealthChecks)
		autoIssues := health.RunAutoChecks(absPath, result)
		issues = append(issues, autoIssues...)

		return scanMsg{result: result, issues: issues}
	}
}

// Update handles messages.
func (a App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.updateLayout()
		return a, nil

	case scanMsg:
		if msg.err != nil {
			a.err = msg.err
			return a, nil
		}
		a.result = msg.result
		a.issues = msg.issues
		a.tree = NewTreeView(msg.result)
		a.healthPanel.SetIssues(msg.issues)
		a.statusBar.SetResult(msg.result)
		a.statusBar.SetTheme(a.themeName)
		a.updateLayout()
		return a, nil

	case tea.KeyMsg:
		return a.handleKey(msg)
	}

	return a, nil
}

func (a App) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	// Search input mode.
	if a.search.IsActive() {
		switch key {
		case "esc":
			a.search.Deactivate()
			a.tree.ClearFilter()
		case "enter":
			a.search.Deactivate()
		case "backspace":
			a.search.Backspace()
			a.tree.SetFilter(a.search.Query())
		default:
			if len(key) == 1 {
				a.search.InsertChar(rune(key[0]))
				a.tree.SetFilter(a.search.Query())
			}
		}
		return a, nil
	}

	// Help overlay.
	if a.help.IsVisible() {
		switch key {
		case "?", "esc", "q":
			a.help.Hide()
		}
		return a, nil
	}

	// Diff view.
	if a.diffView.IsVisible() {
		switch key {
		case "esc", "q":
			a.diffView.Hide()
		case "up", "k":
			a.diffView.ScrollUp(1)
		case "down", "j":
			a.diffView.ScrollDown(1)
		case "v":
			a.diffView.ToggleRedact()
		}
		return a, nil
	}

	// Global keys.
	switch key {
	case "q", "ctrl+c":
		return a, tea.Quit
	case "?":
		a.help.Toggle()
		return a, nil
	case "/":
		a.search.Activate()
		return a, nil
	case "tab":
		a.cycleFocus()
		return a, nil
	case "h":
		a.healthPanel.Toggle()
		a.updateLayout()
		return a, nil
	case "c":
		a.themeName = NextTheme(a.themeName)
		a.theme = GetTheme(a.themeName)
		a.styles = NewStyles(a.theme)
		a.statusBar.SetTheme(a.themeName)
		return a, nil
	case "p":
		a.showPreview = !a.showPreview
		a.updateLayout()
		return a, nil
	case "r":
		return a, a.Init()
	case "d":
		return a.handleDiff()
	case "e":
		return a.openInEditor()
	case "y":
		return a.copyPath()
	case "o":
		return a.openFile()
	}

	// Number keys for category jumping.
	if len(key) == 1 && key[0] >= '1' && key[0] <= '9' {
		n := int(key[0] - '0')
		a.tree.JumpToCategory(n)
		a.loadSelectedFile()
		return a, nil
	}

	// Focus-specific keys.
	switch a.focus {
	case focusTree:
		switch key {
		case "up", "k":
			a.tree.MoveUp()
			a.loadSelectedFile()
		case "down", "j":
			a.tree.MoveDown()
			a.loadSelectedFile()
		case "enter":
			if a.tree.SelectedCategory() != nil {
				a.tree.Toggle()
			} else if f := a.tree.Selected(); f != nil {
				a.loadSelectedFile()
				a.focus = focusViewer
			}
		case "v":
			a.viewer.ToggleRedact()
		}

	case focusViewer:
		switch key {
		case "up", "k":
			a.viewer.ScrollUp(1)
		case "down", "j":
			a.viewer.ScrollDown(1)
		case "esc":
			a.focus = focusTree
		case "v":
			a.viewer.ToggleRedact()
		}

	case focusHealth:
		switch key {
		case "up", "k":
			a.healthPanel.MoveUp()
		case "down", "j":
			a.healthPanel.MoveDown()
		case "enter":
			// Jump to the file of the selected issue.
			// (simplified: just switch to tree focus)
			a.focus = focusTree
		case "esc":
			a.focus = focusTree
		}
	}

	return a, nil
}

func (a *App) cycleFocus() {
	if a.healthPanel.IsVisible() {
		switch a.focus {
		case focusTree:
			if a.showPreview {
				a.focus = focusViewer
			} else {
				a.focus = focusHealth
			}
		case focusViewer:
			a.focus = focusHealth
		case focusHealth:
			a.focus = focusTree
		}
	} else if a.showPreview {
		if a.focus == focusTree {
			a.focus = focusViewer
		} else {
			a.focus = focusTree
		}
	}
}

func (a *App) loadSelectedFile() {
	f := a.tree.Selected()
	if f != nil && a.showPreview {
		a.viewer.LoadFile(f.Path)
	}
}

func (a App) handleDiff() (tea.Model, tea.Cmd) {
	f := a.tree.Selected()
	if f == nil {
		return a, nil
	}
	if !isEnvFile(f.Path) {
		return a, nil
	}

	if a.diffSelect == nil {
		a.diffSelect = f
		return a, nil
	}

	// Second file selected — run diff.
	result, err := differ.Diff(a.diffSelect.Path, f.Path)
	a.diffSelect = nil
	if err != nil {
		return a, nil
	}
	a.diffView.Show(result)
	return a, nil
}

func (a App) openInEditor() (tea.Model, tea.Cmd) {
	f := a.tree.Selected()
	if f == nil {
		return a, nil
	}
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	cmd := exec.Command(editor, f.Path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return a, tea.ExecProcess(cmd, func(err error) tea.Msg { return nil })
}

func (a App) copyPath() (tea.Model, tea.Cmd) {
	f := a.tree.Selected()
	if f == nil {
		return a, nil
	}
	// Use pbcopy on macOS, xclip on Linux.
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	case "linux":
		cmd = exec.Command("xclip", "-selection", "clipboard")
	default:
		return a, nil
	}
	cmd.Stdin = strings.NewReader(f.Path)
	cmd.Run()
	return a, nil
}

func (a App) openFile() (tea.Model, tea.Cmd) {
	f := a.tree.Selected()
	if f == nil {
		return a, nil
	}
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", f.Path)
	case "linux":
		cmd = exec.Command("xdg-open", f.Path)
	default:
		return a, nil
	}
	cmd.Run()
	return a, nil
}

func (a *App) updateLayout() {
	treeHeight := a.height - 3 // status bar + borders
	if a.healthPanel.IsVisible() {
		healthH := min(8, len(a.issues)+2)
		a.healthPanel.SetHeight(healthH)
		treeHeight -= healthH
	}
	if a.search.IsActive() {
		treeHeight--
	}
	a.tree.SetHeight(treeHeight)
	a.viewer.SetHeight(treeHeight)
	a.diffView.SetHeight(treeHeight)
}

// View renders the full TUI.
func (a App) View() string {
	if a.err != nil {
		return fmt.Sprintf("\n  Error: %v\n\n  Press q to quit.\n", a.err)
	}

	if a.result == nil {
		return "\n  Scanning...\n"
	}

	var b strings.Builder

	// Help overlay takes over the screen.
	if a.help.IsVisible() {
		return a.help.View(a.styles, a.width, a.height)
	}

	// Diff view takes over the main area.
	if a.diffView.IsVisible() {
		b.WriteString(a.diffView.View(a.styles, a.width))
		b.WriteByte('\n')
		b.WriteString(a.statusBar.View(a.styles, a.healthPanel.Summary(a.styles), a.width))
		return b.String()
	}

	// Main layout: tree (left) + viewer (right).
	treeWidth := a.width
	if a.showPreview {
		treeWidth = a.width * 2 / 5
		if treeWidth < 30 {
			treeWidth = 30
		}
	}
	viewerWidth := a.width - treeWidth - 1

	treeContent := a.tree.View(a.styles, treeWidth)

	var mainContent string
	if a.showPreview && viewerWidth > 10 {
		viewerContent := a.viewer.View(a.styles, viewerWidth)

		treeBorder := lipgloss.NewStyle().
			Width(treeWidth).
			BorderRight(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(a.theme.Border)

		mainContent = lipgloss.JoinHorizontal(lipgloss.Top,
			treeBorder.Render(treeContent),
			viewerContent,
		)
	} else {
		mainContent = treeContent
	}

	b.WriteString(mainContent)

	// Health panel.
	if a.healthPanel.IsVisible() {
		separator := strings.Repeat("─", a.width)
		b.WriteString(a.styles.Dim.Render(separator))
		b.WriteByte('\n')
		b.WriteString(a.healthPanel.View(a.styles, a.width))
	}

	// Search bar.
	if a.search.IsActive() {
		b.WriteString(a.search.View(a.styles, a.width))
		b.WriteByte('\n')
	}

	// Diff select indicator.
	if a.diffSelect != nil {
		b.WriteString(a.styles.Warning.Render(fmt.Sprintf(
			" d: Select second .env file to diff (first: %s, press Esc to cancel)",
			filepath.Base(a.diffSelect.Path))))
		b.WriteByte('\n')
	}

	// Status bar.
	b.WriteString(a.statusBar.View(a.styles, a.healthPanel.Summary(a.styles), a.width))

	return b.String()
}

