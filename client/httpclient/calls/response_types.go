package calls

import "github.com/Ozoniuss/casheer/pkg/casheerapi"

type ResponseType interface {
	casheerapi.CreateDebtResponse |
		casheerapi.ListDebtResponse |
		casheerapi.UpdateDebtResponse |
		casheerapi.DeleteDebtResponse |
		casheerapi.GetDebtResponse |
		casheerapi.CreateEntryRequest |
		casheerapi.ListEntryResponse |
		casheerapi.CreateEntryResponse |
		casheerapi.CreateExpenseResponse |
		casheerapi.ListExpenseResponse |
		casheerapi.GetExpenseResponse |
		casheerapi.UpdateExpenseResponse |
		casheerapi.DeleteExpenseResponse |
		casheerapi.GetEntryResponse
}
