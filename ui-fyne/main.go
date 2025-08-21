// Functional UI generated with ChatGPT
package main

import (
	"fmt"
	"image/color"
	"maps"
	"math"
	"slices"
	"strconv"
	"time"

	"github.com/Ozoniuss/casheer/client/httpclient"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type appState struct {
	cl                   *httpclient.CasheerHTTPClient
	year, month          int
	categorizedEntryList []CategoryWithEntries
	debts                []DebtListItem

	mainWin           fyne.Window
	entriesAccordion  *widget.Accordion
	debtsList         *widget.List
	headerPeriodLabel *widget.Label
}

/* --------------------------------- Helpers -------------------------------- */

func moneyPow(amount int, exponent int) float32 {
	return float32(math.Pow10(exponent)) * float32(amount)
}

func divider() *canvas.Line {
	l := canvas.NewLine(color.NRGBA{60, 64, 72, 255})
	l.StrokeWidth = 1
	return l
}

func pill(text string) *canvas.Text {
	t := canvas.NewText(text, theme.ForegroundColor())
	t.TextSize = 12
	return t
}

/* ------------------------------- Data loaders ------------------------------ */

func (s *appState) loadDebts() error {
	resp, err := s.cl.ListDebts()
	if err != nil {
		return fmt.Errorf("failed to load debts: %w", err)
	}
	var items []DebtListItem
	for _, d := range resp.Data {
		id, _ := strconv.Atoi(d.Id)
		items = append(items, DebtListItem{
			Id:         id,
			Person:     d.Attributes.Person,
			TotalMoney: moneyPow(d.Attributes.Value.Amount, d.Attributes.Value.Exponent),
			Currency:   d.Attributes.Value.Currency,
			Details:    d.Attributes.Details,
		})
	}
	s.debts = items
	return nil
}

func (s *appState) loadEntries() error {
	resp, err := s.cl.ListEntriesForPeriod(s.month, s.year)
	if err != nil {
		return fmt.Errorf("failed to load entrie: %w", err)
	}
	var entries []EntryListItem
	for _, e := range resp.Data {
		id, _ := strconv.Atoi(e.Id)
		entry := EntryListItem{
			Id:           id,
			TotalMoney:   moneyPow(e.Attributes.ExpectedTotal.Amount, e.Attributes.ExpectedTotal.Exponent),
			Currency:     e.Attributes.ExpectedTotal.Currency,
			Category:     e.Attributes.Category,
			Subcategory:  e.Attributes.Subcategory,
			Recurring:    e.Attributes.Recurring,
			RunningTotal: map[string]float32{"RON": 0, "EUR": 0, "USD": 0, "GBP": 0},
		}
		// include expenses
		full, err := s.cl.GetEntry(id, true)
		if err != nil {
			return err
		}
		if full.Included != nil {
			for _, inc := range *full.Included {
				if inc.Type != "expense" {
					continue
				}
				expid, _ := strconv.Atoi(inc.Id)
				val := moneyPow(inc.Attributes.Value.Amount, inc.Attributes.Value.Exponent)
				entry.Expenses = append(entry.Expenses, ExpenseListItem{
					Id:            expid,
					TotalMoney:    val,
					Currency:      inc.Attributes.Value.Currency,
					Name:          inc.Attributes.Name,
					Description:   inc.Attributes.Description,
					PaymentMethod: inc.Attributes.PaymentMethod,
				})
				entry.RunningTotal[inc.Attributes.Value.Currency] += val
			}
		}
		entries = append(entries, entry)
	}
	// group by category
	cmap := make(map[string][]EntryListItem)
	for _, e := range entries {
		cmap[e.Category] = append(cmap[e.Category], e)
	}
	var out []CategoryWithEntries
	if _, ok := cmap["income"]; ok {
		out = append(out, CategoryWithEntries{"income", cmap["income"]})
	}
	// keys := maps.
	for _, k := range slices.Sorted(maps.Keys(cmap)) {
		if k == "income" {
			continue
		}
		out = append(out, CategoryWithEntries{k, cmap[k]})
	}
	s.categorizedEntryList = out
	return nil
}

/* ------------------------------ UI: Entries ------------------------------- */

func (s *appState) entryRow(e EntryListItem, refreshEntries func()) fyne.CanvasObject {
	// header row
	header := container.NewGridWithColumns(8,
		widget.NewLabelWithStyle(e.Subcategory, fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		widget.NewLabel(fmt.Sprintf("%d", e.Id)),
		widget.NewLabel(fmt.Sprintf("%.2f %s", e.TotalMoney, e.Currency)),
		widget.NewLabel(fmt.Sprintf("%v", e.Recurring)),
		widget.NewLabel(fmt.Sprintf("%.2f", e.RunningTotal["RON"])),
		widget.NewLabel(fmt.Sprintf("%.2f", e.RunningTotal["EUR"])),
		widget.NewLabel(fmt.Sprintf("%.2f", e.RunningTotal["USD"])),
		widget.NewLabel(fmt.Sprintf("%.2f", e.RunningTotal["GBP"])),
	)
	headerPad := container.NewPadded(header)

	// expenses list
	expList := widget.NewList(
		func() int { return len(e.Expenses) },
		func() fyne.CanvasObject {
			name := widget.NewLabel("Title")
			money := widget.NewLabel("0.00")
			del := widget.NewButtonWithIcon("", theme.DeleteIcon(), nil)
			del.Importance = widget.DangerImportance
			row := container.NewBorder(nil, divider(), nil, del,
				container.NewGridWithColumns(3, name, money, widget.NewLabel("")))
			return row
		},
		func(i widget.ListItemID, co fyne.CanvasObject) {
			exp := e.Expenses[i]
			row := co.(*fyne.Container)
			grid := row.Objects[0].(*fyne.Container) // 3-col grid
			grid.Objects[0].(*widget.Label).SetText(exp.Name)
			grid.Objects[1].(*widget.Label).SetText(fmt.Sprintf("%.2f %s", exp.TotalMoney, exp.Currency))

			del := row.Objects[2].(*widget.Button)
			del.OnTapped = func() {
				_, err := s.cl.DeleteExpenseForEntry(e.Id, exp.Id)
				if err != nil {
					dialog.ShowError(err, s.mainWin)
					return
				}
				if err := s.loadEntries(); err != nil {
					dialog.ShowError(err, s.mainWin)
					return
				}
				refreshEntries()
			}
		},
	)
	expCard := widget.NewCard("Expenses", "", container.NewMax(expList))
	expCard.Resize(fyne.NewSize(600, 220))

	// add expense
	addExpense := widget.NewButtonWithIcon("Add expense", theme.ContentAddIcon(), func() {
		s.showAddExpenseDialog(e.Category, func() {
			refreshEntries()
		})
	})
	addExpense.Importance = widget.HighImportance

	body := container.NewVBox(
		headerPad, divider(),
		container.NewHBox(layout.NewSpacer(), addExpense),
		expCard,
	)
	return body
}

func (s *appState) buildEntriesAccordion() *widget.Accordion {
	acc := widget.NewAccordion()
	// Build all panels
	for _, cat := range s.categorizedEntryList {
		c := cat // capture
		expected := 0.0
		for _, e := range c.Entries {
			expected += float64(e.TotalMoney)
		}
		title := fmt.Sprintf("%s   •   Expected: %.0f", c.Category, expected)

		makeContent := func() fyne.CanvasObject {
			rows := []fyne.CanvasObject{
				container.NewPadded(container.NewGridWithColumns(8,
					pill("Subcategory"), pill("Id"), pill("Expected"),
					pill("Recurring"), pill("RON"), pill("EUR"), pill("USD"), pill("GBP"),
				)),
				divider(),
			}
			for _, entry := range c.Entries {
				e := entry
				rows = append(rows, s.entryRow(e, func() {
					// refresh entries only
					if err := s.loadEntries(); err != nil {
						dialog.ShowError(err, s.mainWin)
						return
					}
					s.refreshEntriesOnly()
				}))
				rows = append(rows, divider())
			}
			return container.NewVBox(rows...)
		}
		acc.Append(widget.NewAccordionItem(title, makeContent()))
	}
	return acc
}

/* ------------------------------- UI: Debts -------------------------------- */

func formatDebt(d DebtListItem) string {
	return fmt.Sprintf("#%d  %s — %.2f %s  (%s)", d.Id, d.Person, d.TotalMoney, d.Currency, d.Details)
}

func (s *appState) buildDebtsList() *widget.List {
	list := widget.NewList(
		func() int { return len(s.debts) },
		func() fyne.CanvasObject {
			lbl := widget.NewLabel("debt")
			btn := widget.NewButtonWithIcon("", theme.ConfirmIcon(), nil)
			btn.Importance = widget.WarningImportance
			row := container.NewBorder(nil, divider(), nil, btn, lbl)
			return row
		},
		func(i widget.ListItemID, co fyne.CanvasObject) {
			debt := s.debts[i]
			row := co.(*fyne.Container)
			row.Objects[0].(*widget.Label).SetText(formatDebt(debt))
			btn := row.Objects[2].(*widget.Button)
			btn.SetText("Resolve")
			btn.OnTapped = func() {
				_, err := s.cl.DeleteDebt(debt.Id)
				if err != nil {
					dialog.ShowError(err, s.mainWin)
					return
				}
				if err := s.loadDebts(); err != nil {
					dialog.ShowError(err, s.mainWin)
					return
				}
				s.debtsList.Refresh()
			}
		},
	)
	return list
}

func (s *appState) debtsView() *fyne.Container {
	s.debtsList = s.buildDebtsList()
	add := widget.NewButtonWithIcon("Add debt", theme.ContentAddIcon(), func() {
		s.showAddDebtDialog(func() {
			if err := s.loadDebts(); err != nil {
				dialog.ShowError(err, s.mainWin)
				return
			}
			s.debtsList.Refresh()
		})
	})
	add.Importance = widget.HighImportance

	card := widget.NewCard("Debts", "Outstanding items", container.NewBorder(nil, add, nil, nil, s.debtsList))
	card.Resize(fyne.NewSize(400, 400))
	return container.NewMax(card)
}

/* --------------------------------- Dialogs -------------------------------- */

func (s *appState) showAddEntryDialog(onDone func()) {
	cat := widget.NewEntry()
	sub := widget.NewEntry()
	expTot := widget.NewEntry()
	curr := widget.NewEntry()
	rec := widget.NewCheck("Recurring", nil)
	curr.SetPlaceHolder("RON / EUR / USD / GBP")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Category", Widget: cat},
			{Text: "Subcategory", Widget: sub},
			{Text: "Expected Total", Widget: expTot},
			{Text: "Currency", Widget: curr},
			{Text: "", Widget: rec},
		},
		OnSubmit: func() {
			if cat.Text == "" || sub.Text == "" || expTot.Text == "" {
				dialog.ShowInformation("Missing fields", "Category, Subcategory and Expected Total are required.", s.mainWin)
				return
			}
			exp, err := strconv.Atoi(expTot.Text)
			if err != nil {
				dialog.ShowError(fmt.Errorf("invalid expected total"), s.mainWin)
				return
			}
			currency := curr.Text
			if currency == "" {
				currency = "RON"
			}
			if _, err := s.cl.CreateEntry(s.month, s.year, cat.Text, sub.Text, exp, currency, rec.Checked); err != nil {
				dialog.ShowError(err, s.mainWin)
				return
			}
			if err := s.loadEntries(); err != nil {
				dialog.ShowError(err, s.mainWin)
				return
			}
			onDone()
		},
		SubmitText: "Create entry",
	}
	d := dialog.NewCustom("New entry", "Close", form, s.mainWin)
	d.Show()
}

func (s *appState) showAddExpenseDialog(category string, onDone func()) {
	sub := widget.NewEntry()
	name := widget.NewEntry()
	amt := widget.NewEntry()
	curr := widget.NewEntry()
	details := widget.NewEntry()
	pay := widget.NewEntry()
	curr.SetPlaceHolder("RON / EUR / USD / GBP")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Subcategory", Widget: sub},
			{Text: "Name", Widget: name},
			{Text: "Amount", Widget: amt},
			{Text: "Currency", Widget: curr},
			{Text: "Details", Widget: details},
			{Text: "Payment method", Widget: pay},
		},
		OnSubmit: func() {
			if sub.Text == "" || amt.Text == "" || curr.Text == "" {
				dialog.ShowInformation("Missing fields", "Subcategory, Amount and Currency are required.", s.mainWin)
				return
			}
			value, err := strconv.Atoi(amt.Text)
			if err != nil {
				dialog.ShowError(fmt.Errorf("invalid amount"), s.mainWin)
				return
			}
			// find entry id; create if missing
			var targetEntryID int = -1
			for _, c := range s.categorizedEntryList {
				if c.Category != category {
					continue
				}
				for _, e := range c.Entries {
					if e.Subcategory == sub.Text {
						targetEntryID = e.Id
						break
					}
				}
			}
			if targetEntryID == -1 {
				created, err := s.cl.CreateEntry(s.month, s.year, category, sub.Text, value, curr.Text, false)
				if err != nil {
					dialog.ShowError(err, s.mainWin)
					return
				}
				targetEntryID, _ = strconv.Atoi(created.Data.Id)
			}
			if _, err := s.cl.CreateBasicExpense(targetEntryID, name.Text, details.Text, pay.Text, value, curr.Text); err != nil {
				dialog.ShowError(err, s.mainWin)
				return
			}
			if err := s.loadEntries(); err != nil {
				dialog.ShowError(err, s.mainWin)
				return
			}
			onDone()
		},
		SubmitText: "Create expense",
	}
	d := dialog.NewCustom(fmt.Sprintf("New expense in %s", category), "Close", form, s.mainWin)
	d.Show()
}

func (s *appState) showAddDebtDialog(onDone func()) {
	person := widget.NewEntry()
	amt := widget.NewEntry()
	curr := widget.NewEntry()
	dets := widget.NewEntry()
	curr.SetPlaceHolder("RON / EUR / USD / GBP")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Person", Widget: person},
			{Text: "Amount", Widget: amt},
			{Text: "Currency", Widget: curr},
			{Text: "Details", Widget: dets},
		},
		OnSubmit: func() {
			if person.Text == "" || amt.Text == "" || curr.Text == "" {
				dialog.ShowInformation("Missing fields", "Person, Amount, Currency are required.", s.mainWin)
				return
			}
			value, err := strconv.Atoi(amt.Text)
			if err != nil {
				dialog.ShowError(fmt.Errorf("invalid amount"), s.mainWin)
				return
			}
			if _, err := s.cl.CreateDebt(person.Text, dets.Text, value, curr.Text, -2); err != nil {
				dialog.ShowError(err, s.mainWin)
				return
			}
			if err := s.loadDebts(); err != nil {
				dialog.ShowError(err, s.mainWin)
				return
			}
			onDone()
		},
		SubmitText: "Create debt",
	}
	d := dialog.NewCustom("New debt", "Close", form, s.mainWin)
	d.Show()
}

/* ---------------------------- Refresh / Rebuild ---------------------------- */

func (s *appState) refreshEntriesOnly() {
	s.entriesAccordion.Items = s.buildEntriesAccordion().Items
	s.entriesAccordion.Refresh()
	s.headerPeriodLabel.SetText(fmt.Sprintf("Period: %d-%02d", s.year, s.month))
}

func (s *appState) refreshAllData() {
	if err := s.loadEntries(); err != nil {
		dialog.ShowError(err, s.mainWin)
		return
	}
	if err := s.loadDebts(); err != nil {
		dialog.ShowError(err, s.mainWin)
		return
	}
	// rebuild entries accordion
	s.entriesAccordion.Items = s.buildEntriesAccordion().Items
	s.entriesAccordion.Refresh()
	// debts
	s.debtsList.Refresh()
	s.headerPeriodLabel.SetText(fmt.Sprintf("Period: %d-%02d", s.year, s.month))
}

/* --------------------------------- Header -------------------------------- */

func (s *appState) headerBar() fyne.CanvasObject {
	years := []string{"2022", "2023", "2024", "2025"}
	months := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"}

	yearSel := widget.NewSelect(years, nil)
	monthSel := widget.NewSelect(months, nil)

	yearSel.Selected = fmt.Sprintf("%d", s.year)
	monthSel.Selected = fmt.Sprintf("%d", s.month)

	yearSel.OnChanged = func(v string) {
		if y, err := strconv.Atoi(v); err == nil {
			s.year = y
			s.refreshAllData()
		}
	}
	monthSel.OnChanged = func(v string) {
		if m, err := strconv.Atoi(v); err == nil {
			s.month = m
			s.refreshAllData()
		}
	}

	newEntry := widget.NewButtonWithIcon("New entry", theme.ContentAddIcon(), func() {
		s.showAddEntryDialog(func() { s.refreshAllData() })
	})
	newEntry.Importance = widget.HighImportance

	left := container.NewHBox(
		canvas.NewText("Casheer", theme.PrimaryColor()),
		layout.NewSpacer(),
		widget.NewLabel("Year"),
		yearSel,
		widget.NewLabel("Month"),
		monthSel,
	)
	right := container.NewHBox(newEntry, widget.NewButtonWithIcon("Refresh", theme.ViewRefreshIcon(), func() { s.refreshAllData() }))
	return container.NewBorder(nil, divider(), left, right, s.headerPeriodLabel)
}

/* ----------------------------------- main --------------------------------- */

func main() {
	a := app.New()
	a.Settings().SetTheme(&sleekDark{})

	w := a.NewWindow("Casheer — Planner")
	w.Resize(fyne.NewSize(1200, 720))

	cl, err := httpclient.NewCasheerHTTPClient()
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	now := time.Now()
	state := &appState{
		cl:      cl,
		year:    now.Year(),
		month:   int(now.Month()),
		mainWin: w,
	}

	state.headerPeriodLabel = widget.NewLabel("")
	header := state.headerBar()

	state.entriesAccordion = state.buildEntriesAccordion()
	entriesCard := widget.NewCard("Entries", "Predictions vs actuals", container.NewMax(state.entriesAccordion))
	entriesCard.Resize(fyne.NewSize(800, 600))

	state.debtsList = state.buildDebtsList()
	debts := state.debtsView()

	content := container.NewBorder(header, nil, nil, debts, container.NewVScroll(entriesCard))
	w.SetContent(content)

	state.refreshAllData()

	w.ShowAndRun()
}
