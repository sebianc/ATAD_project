package services

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"

	"atad_project/models"

	bubbletea "github.com/charmbracelet/bubbletea"
)

type transactionItem struct {
	T *models.Transaction
}

func (t transactionItem) Title() string {
	return fmt.Sprintf("%s | %.2f | %s", t.T.DATE, t.T.AMOUNT, t.T.CATEGORY)
}

func (t transactionItem) Description() string {
	return t.T.DESCRIPTION
}

func (t transactionItem) FilterValue() string {
	return t.T.DATE + " " + fmt.Sprintf("%.2f", t.T.AMOUNT) + " " + t.T.DESCRIPTION + " " + t.T.CATEGORY
}

type model struct {
	list      list.Model
	ti        textinput.Model
	filtering bool
	filterBy  string
	all       []*models.Transaction
}

func NewTransactionTUI(transactions []*models.Transaction) model {
	items := make([]list.Item, len(transactions))
	for i, tx := range transactions {
		items[i] = transactionItem{T: tx}
	}
	l := list.New(items, list.NewDefaultDelegate(), 100, 30)
	l.Title = "Transactions Browser"

	ti := textinput.New()
	ti.Placeholder = "filter value"
	ti.CharLimit = 128
	ti.Width = 40

	return model{list: l, ti: ti, filtering: false, filterBy: "description", all: transactions}
}

func (m model) Init() bubbletea.Cmd {
	return nil
}

func (m model) Update(msg bubbletea.Msg) (bubbletea.Model, bubbletea.Cmd) {
	var cmd bubbletea.Cmd

	switch msg := msg.(type) {
	case bubbletea.KeyMsg:
		if m.filtering {
			switch msg.String() {
			case "enter":
				val := strings.TrimSpace(m.ti.Value())
				filtered := FilterTransactions(m.all, m.filterBy, val)
				items := make([]list.Item, len(filtered))
				for i, tx := range filtered {
					items[i] = transactionItem{T: tx}
				}
				m.list.SetItems(items)
				m.filtering = false
				m.ti.Blur()
				m.ti.SetValue("")
				return m, nil
			case "esc":
				m.filtering = false
				m.ti.Blur()
				m.ti.SetValue("")
				return m, nil
			}
			var cmd2 bubbletea.Cmd
			m.ti, cmd2 = m.ti.Update(msg)
			return m, cmd2
		}

		switch msg.String() {
		case "/":
			m.filtering = true
			m.ti.Focus()
			return m, nil
		case "1":
			m.filterBy = "date"
			return m, nil
		case "2":
			m.filterBy = "amount"
			return m, nil
		case "3":
			m.filterBy = "category"
			return m, nil
		case "4":
			m.filterBy = "description"
			return m, nil
		case "q", "ctrl+c":
			return m, bubbletea.Quit
		}
	}

	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	header := fmt.Sprintf("[ / ] Filter  [1] Date  [2] Amount  [3] Category  [4] Description  [q] Quit  | Current filter: %s\n", m.filterBy)
	if m.filtering {
		return header + "Filter value: " + m.ti.View() + "\n\n" + m.list.View()
	}
	return header + m.list.View()
}

func RunTransactionTUI(transactions []*models.Transaction) error {
	p := bubbletea.NewProgram(NewTransactionTUI(transactions))
	return p.Start()
}

func FilterTransactions(transactions []*models.Transaction, by, val string) []*models.Transaction {
	by = strings.ToLower(strings.TrimSpace(by))
	val = strings.TrimSpace(val)
	if val == "" {
		return transactions
	}

	switch by {
	case "amount":
		return filterByAmount(transactions, val)
	case "category":
		return filterByTextField(transactions, val, func(t *models.Transaction) string { return t.CATEGORY })
	case "description":
		return filterByTextField(transactions, val, func(t *models.Transaction) string { return t.DESCRIPTION })
	case "date":
		return filterByDate(transactions, val)
	default:
		return transactions
	}
}

func filterByTextField(transactions []*models.Transaction, val string, getter func(*models.Transaction) string) []*models.Transaction {
	var out []*models.Transaction
	q := strings.ToLower(val)
	for _, t := range transactions {
		if strings.Contains(strings.ToLower(getter(t)), q) {
			out = append(out, t)
		}
	}
	return out
}

func filterByDate(transactions []*models.Transaction, val string) []*models.Transaction {
	var out []*models.Transaction
	// range using ':' as separator start:end
	if strings.Contains(val, ":") {
		parts := strings.SplitN(val, ":", 2)
		start := strings.TrimSpace(parts[0])
		end := strings.TrimSpace(parts[1])
		for _, t := range transactions {
			if t.DATE >= start && t.DATE <= end {
				out = append(out, t)
			}
		}
		return out
	}

	// otherwise substring match
	q := strings.ToLower(val)
	for _, t := range transactions {
		if strings.Contains(strings.ToLower(t.DATE), q) {
			out = append(out, t)
		}
	}
	return out
}

func filterByAmount(transactions []*models.Transaction, val string) []*models.Transaction {
	var out []*models.Transaction
	// range with '-'
	if strings.Contains(val, "-") {
		parts := strings.SplitN(val, "-", 2)
		lo, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		hi, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err1 == nil && err2 == nil {
			if lo > hi {
				lo, hi = hi, lo
			}
			for _, t := range transactions {
				if t.AMOUNT >= lo && t.AMOUNT <= hi {
					out = append(out, t)
				}
			}
			return out
		}
	}

	if strings.HasPrefix(val, ">=") || strings.HasPrefix(val, "<=") {
		op := val[:2]
		numStr := strings.TrimSpace(val[2:])
		num, err := strconv.ParseFloat(numStr, 64)
		if err == nil {
			for _, t := range transactions {
				switch op {
				case ">=":
					if t.AMOUNT >= num {
						out = append(out, t)
					}
				case "<=":
					if t.AMOUNT <= num {
						out = append(out, t)
					}
				}
			}
			return out
		}
	}

	if strings.HasPrefix(val, ">") || strings.HasPrefix(val, "<") {
		op := val[:1]
		numStr := strings.TrimSpace(val[1:])
		num, err := strconv.ParseFloat(numStr, 64)
		if err == nil {
			for _, t := range transactions {
				switch op {
				case ">":
					if t.AMOUNT > num {
						out = append(out, t)
					}
				case "<":
					if t.AMOUNT < num {
						out = append(out, t)
					}
				}
			}
			return out
		}
	}

	// exact match
	if num, err := strconv.ParseFloat(val, 64); err == nil {
		for _, t := range transactions {
			if t.AMOUNT == num {
				out = append(out, t)
			}
		}
		return out
	}

	return out
}
