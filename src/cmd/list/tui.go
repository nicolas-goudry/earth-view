package list

import (
  "strconv"
  "strings"
  "time"

  "github.com/charmbracelet/bubbles/progress"
  tea "github.com/charmbracelet/bubbletea"
  "github.com/charmbracelet/lipgloss"
)

const padding = 2
var (
  red = lipgloss.AdaptiveColor{
    Light: "#d20f39",
    Dark: "#f38ba8",
  }
  green = lipgloss.AdaptiveColor{
    Light: "#40a02b",
    Dark: "#a6e3a1",
  }
  orange = lipgloss.AdaptiveColor{
    Light: "#fe640b",
    Dark: "#fab387",
  }
)

// Defer program end to allow progress bar to go to 100% before UI is cleared
func finalPause() tea.Cmd {
  return tea.Tick(time.Second, func(_ time.Time) tea.Msg {
    return fetchCompleteMsg(true)
  })
}

// Custom messages
type progressMsg      fetchProgress
type fetchCompleteMsg bool

// UI state
type model struct {
  progress progress.Model
  success  int
  skipped  int
  errored  int
  abort    bool
  clear    bool
}

func (m model) Init() tea.Cmd {
  return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  var cmds []tea.Cmd

  switch msg := msg.(type) {
  case tea.KeyMsg:
    key := msg.String()
    if key == "ctrl+c" {
      m.abort = true
      return m, tea.Quit
    }

  case tea.WindowSizeMsg:
    m.progress.Width = msg.Width - padding * 2 - 4
    return m, nil

  case progress.FrameMsg:
    progressModel, cmd := m.progress.Update(msg)
    m.progress = progressModel.(progress.Model)
    return m, cmd

  case progressMsg:
    for _, result := range msg.results {
      if result.error != nil && strings.Contains(result.error.Error(), "not found") {
        m.skipped++
      } else if result.error != nil {
        m.errored++
      } else {
        m.success++
      }
    }

    var cmds []tea.Cmd

    if msg.percent >= 1.0 {
      cmds = append(cmds, tea.Sequence(finalPause(), tea.Quit))
    }

    cmds = append(cmds, m.progress.SetPercent(float64(msg.percent)))
    return m, tea.Batch(cmds...)

  case fetchCompleteMsg:
    m.clear = bool(msg)
    return m, nil
  }

  return m, tea.Batch(cmds...)
}

func (m model) View() string {
  if quiet {
    return ""
  }

  if m.clear {
    return ""
  }

  if m.abort {
    return lipgloss.NewStyle().
      PaddingBottom(1).
      Foreground(red).
      Render("Operation aborted before end")
  }

  separator := lipgloss.NewStyle().Padding(0, padding).Render("•")

  return lipgloss.JoinVertical(
    lipgloss.Left,
    lipgloss.NewStyle().Padding(1, 0, 1, padding).Bold(true).Render("Searching images..."),
    lipgloss.NewStyle().Padding(0, 0, 1, padding).Render(m.progress.View()),
    lipgloss.JoinHorizontal(
      lipgloss.Top,
      lipgloss.NewStyle().
        PaddingLeft(padding).
        PaddingBottom(1).
        Foreground(green).
        Render("Found: " + strconv.Itoa(m.success)),
      separator,
      lipgloss.NewStyle().
        Foreground(orange).
        Render("Skipped: " + strconv.Itoa(m.skipped)),
      separator,
      lipgloss.NewStyle().
        Foreground(red).
        Render("Errors: " + strconv.Itoa(m.errored)),
    ),
  )
}
