package main

import (
	"context"
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/rivo/tview"
)

type Page struct {
	searchField  *tview.InputField
	searchCancel context.CancelFunc
	gridView     *tview.Flex
	columnView   *tview.Flex
	listView     *tview.List
	textView     *tview.TextView
	listItems    []ListItem
	currentItem  ListItem
}

func (p *Page) InputFieldChangedFunc(text string) {
	if p.searchCancel != nil {
		p.searchCancel()
		p.searchCancel = nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	p.searchCancel = cancel

	go func() {
		// TODO: figure out if this can be abstracted so that we don't call the global app here
		app.tviewApp.QueueUpdateDraw(func() {
			p.PopulateListView(ctx, app.switchToPageFunc)
		})
	}()
}

func (p *Page) CreateSearchField(SetFocus func(tview.Primitive), SwitchToPage func(string) *tview.Pages) {
	p.searchField = tview.NewInputField().
		SetLabel("Search: ")
	p.searchField.SetChangedFunc(p.InputFieldChangedFunc)

	p.searchField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab && p.searchField.HasFocus() {
			SetFocus(p.listView)
		} else if event.Key() == tcell.KeyEnter {
			SwitchToPage("issues" + p.currentItem.Name())
		}
		return event
	})
}

func (p *Page) CreatePageGrid() {
	p.listView.SetBorder(true)
	p.textView.SetBorder(true)

	p.columnView = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(p.searchField, 1, 0, true).
		AddItem(p.listView, 0, 1, false)

	p.gridView = tview.NewFlex().
		AddItem(p.columnView, 45, 1, true).
		AddItem(p.textView, 0, 3, false)

	p.gridView.SetInputCapture(viewInputCapture)
}

func (page *Page) PopulateListView(ctx context.Context, listFn func(string) *tview.Pages) {
	var items []ListItem
	if page.searchField.GetText() == "" {
		items = page.listItems
	} else {
		items = page.searchItems(ctx)
	}
	select {
	case <-ctx.Done():
		return
	default:
		page.listView.Clear()
		for _, item := range items {
			var pageName string = "issues" + item.Name()
			page.listView.AddItem(item.Name(), string(item.ID()), rune(0), func() { listFn(pageName) })
		}
	}
}

func (p *Page) searchItems(ctx context.Context) []ListItem {
	if ctx.Err() != nil {
		return nil
	}

	searchString := p.searchField.GetText()

	itemNamesMap := make(map[string]ListItem, len(p.listItems))
	itemNames := make([]string, 0, len(p.listItems))

	for _, item := range p.listItems {
		select {
		case <-ctx.Done():
			return nil
		default:
			itemNames = append(itemNames, item.Name())
			itemNamesMap[item.Name()] = item
		}
	}

	names := fuzzy.RankFindFold(searchString, itemNames)
	sort.Sort(names)
	projects := addToArray([]ListItem{}, itemNamesMap, names)
	return projects
}
