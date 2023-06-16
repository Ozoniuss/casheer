package entries

import (
	"fmt"
	"strconv"

	"github.com/Ozoniuss/casheer/internal/config"
	"github.com/Ozoniuss/casheer/internal/model"
	public "github.com/Ozoniuss/casheer/pkg/casheerapi"
)

// computeRunningTotal returns the running total of a list of expenses.
func computeRunningTotal(expenses []model.Expense) int64 {
	var rt int64 = 0
	for _, exp := range expenses {
		rt += exp.Value
	}
	return rt
}

func EntryToPublic(entry model.Entry, apipath config.ApiPaths, runningTotal int64) public.EntryData {
	return public.EntryData{
		ResourceID: public.ResourceID{
			Id:   strconv.Itoa(entry.Id),
			Type: public.EntryType,
		},
		Attributes: public.EntryAtrributes{
			Year:          int(entry.Year),
			Month:         int(entry.Month),
			Category:      entry.Category,
			Subcategory:   entry.Subcategory,
			ExpectedTotal: entry.ExpectedTotal,
			Recurring:     entry.Recurring,
			Timestamps: public.Timestamps{
				CreatedAt: entry.CreatedAt,
				UpdatedAt: entry.UpdatedAt,
			},
		},
		Meta: public.EntryMeta{
			RunningTotal: runningTotal,
		},
		Links: public.EntryLinks{
			Self:       fmt.Sprintf("%s/%s", apipath.Entries, strconv.Itoa(entry.Id)),
			Collection: fmt.Sprintf("%s/", apipath.Entries),
			Expenses:   fmt.Sprintf("%s/%s/expenses/", apipath.Entries, strconv.Itoa(entry.Id)),
		},
	}
}

func EntryToPublicList(entry model.Entry, apipath config.ApiPaths, runningTotal int64) public.EntryListItemData {
	return public.EntryListItemData{
		ResourceID: public.ResourceID{
			Id:   strconv.Itoa(entry.Id),
			Type: public.EntryType,
		},
		Attributes: public.EntryAtrributes{
			Year:          int(entry.Year),
			Month:         int(entry.Month),
			Category:      entry.Category,
			Subcategory:   entry.Subcategory,
			ExpectedTotal: entry.ExpectedTotal,
			Recurring:     entry.Recurring,
			Timestamps: public.Timestamps{
				CreatedAt: entry.CreatedAt,
				UpdatedAt: entry.UpdatedAt,
			},
		},
		Meta: public.EntryMeta{
			RunningTotal: runningTotal,
		},
		Links: public.EntryListItemLinks{
			Self: fmt.Sprintf("%s/%s", apipath.Entries, strconv.Itoa(entry.Id)),
		},
	}
}
