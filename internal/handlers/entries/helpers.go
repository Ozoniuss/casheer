package entries

import (
	"github.com/Ozoniuss/casheer/internal/model"
	public "github.com/Ozoniuss/casheer/pkg/casheerapi"
)

func EntryToPublic(entry model.Entry) public.EntryData {
	return public.EntryData{
		ResourceID: public.ResourceID{
			Id:   entry.Id.String(),
			Type: public.EntryType,
		},
		Year:          int(entry.Year),
		Month:         int(entry.Month),
		Category:      entry.Category,
		Subcategory:   entry.Subcategory,
		ExpectedTotal: entry.ExpectedTotal,
		RunningTotal:  entry.RunningTotal,
		Recurring:     entry.Recurring,
		Timestamps: public.Timestamps{
			CreatedAt: entry.CreatedAt,
			UpdatedAt: entry.UpdatedAt,
		},
	}
}
