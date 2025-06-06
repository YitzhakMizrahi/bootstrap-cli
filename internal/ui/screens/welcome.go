package screens

import (
	"fmt"
	"strings"

	"github.com/YitzhakMizrahi/bootstrap-cli/internal/system"
	"github.com/YitzhakMizrahi/bootstrap-cli/internal/ui/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// WelcomeScreen handles the welcome screen display & system info
type WelcomeScreen struct {
	quitting bool
	done     bool
	width    int
	height   int
	sysInfo  *system.Info
}

// NewWelcomeScreen creates a new welcome screen
func NewWelcomeScreen() *WelcomeScreen {
	return &WelcomeScreen{sysInfo: nil}
}

// SetInfo allows the parent model to inject system info when ready
func (w *WelcomeScreen) SetInfo(info *system.Info) {
	w.sysInfo = info
}

// Init implements tea.Model
func (w *WelcomeScreen) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (w *WelcomeScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			w.quitting = true
			return w, tea.Quit
		case "enter":
			w.done = true
			return w, nil
		}
	case tea.WindowSizeMsg:
		w.width = msg.Width
		w.height = msg.Height
	}
	return w, nil
}

// View implements tea.Model
func (w *WelcomeScreen) View() string {
	if w.quitting {
		return ""
	}

	// --- ASCII Art Title ---
	igniteArt := `
      (  .      )
   )           (              )
         .  '   .   '  .  '  .
  (    , )       (.   )  (   ',    )
   .' ) ( . )    ,  ( ,     )   ( .
). , ( .   (  ) ( , ')  .' (  ,    )
(_,) . ), ) _) _,')  (, ) '. )  ,. (' )

██╗ ██████╗ ███╗   ██╗██╗████████╗███████╗
██║██╔════╝ ████╗  ██║██║╚══██╔══╝██╔════╝
██║██║  ███╗██╔██╗ ██║██║   ██║   █████╗  
██║██║   ██║██║╚██╗██║██║   ██║   ██╔══╝  
██║╚██████╔╝██║ ╚████║██║   ██║   ███████╗
╚═╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝   ╚═╝   ╚══════╝
`
	styledArt := styles.TitleStyle.Copy(). 
		Foreground(styles.ColorAccent).
		Align(lipgloss.Center).
		Render(igniteArt)
	// --- End ASCII Art ---

	var content strings.Builder

	content.WriteString(styledArt)
	content.WriteString("\n\n") 

	// --- System Info ---
	sysInfoStr := ""
	if w.sysInfo != nil {
		labelStyle := styles.NormalTextStyle.Copy().Width(15)
		valueStyle := styles.NormalTextStyle.Copy().Foreground(styles.ColorDimText)
		infoLines := []string{
			lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("OS:"), valueStyle.Render(fmt.Sprintf("%s %s", w.sysInfo.Distro, w.sysInfo.Version))),
			lipgloss.JoinHorizontal(lipgloss.Left, labelStyle.Render("Architecture:"), valueStyle.Render(w.sysInfo.Arch)),
		}
		sysInfoStr = lipgloss.JoinVertical(lipgloss.Left, infoLines...)
	} else {
		sysInfoStr = styles.HelpStyle.Render("Detecting system info...")
	}
	content.WriteString(sysInfoStr)
	content.WriteString("\n\n")

	// Subtitle/Description
	subtitle := styles.SubtitleStyle.Render("Setup your development environment with ease.")
	content.WriteString(subtitle)
	content.WriteString("\n\n\n")

	// Help text
	helpText := styles.HelpStyle.Render("Press Enter to continue.")
	content.WriteString(helpText)

	// Center the entire block using lipgloss.Place
	centeredContent := lipgloss.Place(
		w.width,
		w.height,
		lipgloss.Center, // Horizontal alignment
		lipgloss.Center, // Vertical alignment
		content.String(),
	)

	return centeredContent
}

// Finished returns true if the screen exited normally (not by quitting)
func (w *WelcomeScreen) Finished() bool {
	return w.done && !w.quitting
} 