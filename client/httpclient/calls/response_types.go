package calls

import "github.com/Ozoniuss/casheer/pkg/casheerapi"

type ResponseType interface {
	casheerapi.CreateDebtResponse |
		casheerapi.ListDebtResponse |
		casheerapi.UpdateDebtResponse |
		casheerapi.DeleteDebtResponse |
		casheerapi.GetDebtResponse |
		casheerapi.CreateEntryRequest |
		casheerapi.CreateEntryResponse
}
