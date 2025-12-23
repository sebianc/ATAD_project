package services

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"atad_project/models"
)

type filterMode int

// filter modes
const (
	modeNone filterMode = iota
	modeYearMonth
	modeDescription
	modeCategory
)

type screen int

// screens
const (
	screenMenu screen = iota
	screenFilter
)

// model holds the state of the TUI
type model struct {
	all      []*models.Transaction
	filtered []*models.Transaction

	screen screen
	mode   filterMode

	input textinput.Model
	tbl   table.Model

	status string
	width  int
	height int
}

// styles
var (
	titleStyle  = lipgloss.NewStyle().Bold(true)
	helpStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
)

// RunTransactionTUI runs the TUI for browsing transactions
func RunTransactionTUI(transactions []*models.Transaction) error {
	m := newModel(transactions)
	_, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	return err
}

// newModel creates a new TUI model
func newModel(tx []*models.Transaction) model {
	cols := []table.Column{
		{Title: "YYYY-MM", Width: 9},
		{Title: "Amount", Width: 10},
		{Title: "Category", Width: 14},
		{Title: "Description", Width: 60},
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithHeight(16),
	)

	in := textinput.New()
	in.Prompt = "Filter: "
	in.Placeholder = ""
	in.CharLimit = 80
	in.Width = 40
	in.Blur()

	m := model{
		all:      append([]*models.Transaction{}, tx...),
		filtered: append([]*models.Transaction{}, tx...),
		screen:   screenMenu,
		mode:     modeNone,
		input:    in,
		tbl:      t,
		status:   "Choose filter: 1=YYYY-MM, 2=Description, 3=Category. q quits.",
	}

	m.refreshTable()
	return m
}

// empty Init function
func (m model) Init() tea.Cmd { return nil }

// Update handles incoming messages and updates the model accordingly
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch v := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = v.Width, v.Height

		// resize description column + table height
		descW := maxInt(20, m.width-(11+10+14)-10)
		cols := m.tbl.Columns()
		cols[3].Width = descW
		m.tbl.SetColumns(cols)
		m.tbl.SetHeight(maxInt(8, m.height-8))
		return m, nil

	case tea.KeyMsg:
		switch v.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

		if m.screen == screenMenu {
			return m.updateMenu(v)
		}
		return m.updateFilter(v)
	}

	// non-key: only input updates matter in filter screen
	if m.screen == screenFilter {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		m.applyFilter()
		return m, cmd
	}

	return m, nil
}

// updateMenu handles key events in the menu screen
func (m model) updateMenu(k tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch k.String() {
	case "1":
		m.mode = modeYearMonth
		m.switchToFilter("Year-Month", "YYYY-MM (e.g. 2025-01)")
	case "2":
		m.mode = modeDescription
		m.switchToFilter("Description", "type text (case-insensitive)")
	case "3":
		m.mode = modeCategory
		m.switchToFilter("Category", "type text (case-insensitive)")
	}
	return m, nil
}

// updateFilter handles key events in the filter screen
func (m model) updateFilter(k tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch k.String() {
	case "esc":
		// back to menu
		m.screen = screenMenu
		m.mode = modeNone
		m.input.Blur()
		m.input.SetValue("")
		m.filtered = append([]*models.Transaction{}, m.all...)
		m.refreshTable()
		m.status = "Choose filter: 1=YYYY-MM, 2=Description, 3=Category. q quits."
		return m, nil

	case "down", "up", "pgdown", "pgup", "home", "end":
		var cmd tea.Cmd
		m.tbl, cmd = m.tbl.Update(k)
		return m, cmd
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(k)
	m.applyFilter()
	return m, cmd
}

// switchToFilter switches the model to filter screen with given label and placeholder
func (m *model) switchToFilter(label, placeholder string) {
	m.screen = screenFilter
	m.input.SetValue("")
	m.input.Placeholder = placeholder
	m.input.Prompt = fmt.Sprintf("%s: ", label)
	m.input.Focus()
	m.status = "Type to filter. Esc=back, q=quit."
	m.filtered = append([]*models.Transaction{}, m.all...)
	m.refreshTable()
}

// applyFilter filters the transactions based on the current input and mode
func (m *model) applyFilter() {
	q := strings.TrimSpace(m.input.Value())
	if q == "" {
		m.filtered = append([]*models.Transaction{}, m.all...)
		m.status = statusStyle.Render(fmt.Sprintf("%d/%d shown", len(m.filtered), len(m.all)))
		m.refreshTable()
		return
	}

	switch m.mode {
	case modeYearMonth:
		out := make([]*models.Transaction, 0)
		for _, t := range m.all {
			if strings.HasPrefix(t.DATE, q) {
				out = append(out, t)
			}
		}
		m.filtered = out

	case modeDescription:
		qq := strings.ToLower(q)
		out := make([]*models.Transaction, 0)
		for _, t := range m.all {
			if strings.Contains(strings.ToLower(t.DESCRIPTION), qq) {
				out = append(out, t)
			}
		}
		m.filtered = out

	case modeCategory:
		qq := strings.ToLower(q)
		out := make([]*models.Transaction, 0)
		for _, t := range m.all {
			if strings.Contains(strings.ToLower(t.CATEGORY), qq) {
				out = append(out, t)
			}
		}
		m.filtered = out
	}

	m.status = statusStyle.Render(fmt.Sprintf("%d/%d shown", len(m.filtered), len(m.all)))
	m.refreshTable()
}

// refreshTable updates the table rows based on the filtered transactions
func (m *model) refreshTable() {
	m.tbl.SetRows(toRows(m.filtered))
}

// toRows converts transactions to table rows
func toRows(tx []*models.Transaction) []table.Row {
	rows := make([]table.Row, 0, len(tx))
	for _, t := range tx {
		rows = append(rows, table.Row{
			t.DATE, // already YYYY-MM
			fmt.Sprintf("%.2f", t.AMOUNT),
			t.CATEGORY,
			t.DESCRIPTION,
		})
	}
	return rows
}

// View renders the TUI based on the current screen
func (m model) View() string {
	switch m.screen {
	case screenMenu:
		return lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render("Browse Transactions"),
			"",
			"Choose filter mode:",
			"  1) Year-Month (YYYY-MM)",
			"  2) Description",
			"  3) Category",
			"",
			helpStyle.Render("Press 1/2/3. q to quit."),
			"",
		)

	default:
		header := titleStyle.Render("Browse Transactions") + "  " + helpStyle.Render("(Esc back • q quit • arrows to navigate)")
		return lipgloss.JoinVertical(lipgloss.Left,
			header,
			"",
			m.input.View(),
			"",
			m.tbl.View(),
			"",
			m.status,
		)
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
