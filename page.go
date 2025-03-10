package main

import (
	"context"
	"sort"
	"strings"

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
		filteredItems := []ListItem{}
		if text == "" {
			filteredItems = p.listItems
		} else {
			filteredItems = p.searchItems(ctx)
		}

		select {
		case <-ctx.Done():
			return
		default:
			app.tviewApp.QueueUpdateDraw(func() {
				// TODO: figure out if this can be abstracted so that we don't call the global app here
				p.updateListViewWithItems(ctx, filteredItems, app.switchToPageFunc)
			})
		}
	}()
}

// TODO: pretty sure this fucked up the Issues pages, fix it
func (p *Page) updateListViewWithItems(ctx context.Context, items []ListItem, listFn func(string) *tview.Pages) {
	select {
	case <-ctx.Done():
		return
	default:
		p.listView.Clear()
		for _, item := range items {
			var pageName string = "issues" + item.Name()
			p.listView.AddItem(item.Name(), string(item.ID()), rune(0), func() { listFn(pageName) })
		}
	}
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

	page.updateListViewWithItems(ctx, items, listFn)
}

func (p *Page) searchItems(ctx context.Context) []ListItem {
	if ctx.Err() != nil {
		return nil
	}

	searchString := p.searchField.GetText()

	name, _ := app.pages.GetFrontPage()

	// if name containst "issue" then we are on the issue page
	if name[:5] == "issue" {
	} else if name == "projects" {
		projects := app.getProjects(searchString)
		sort.Slice(projects, func(i, j int) bool {
			return strings.ToLower(projects[i].Name) < strings.ToLower(projects[j].Name)
		})
		p.listItems = []ListItem{}
		for _, project := range projects {
			p.listItems = append(p.listItems, ProjectWrapper{project})
		}
		return p.listItems
	}

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
