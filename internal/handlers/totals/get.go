package totals

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Ozoniuss/casheer/internal/handlers/common"
	"github.com/Ozoniuss/casheer/internal/model"
	"github.com/Ozoniuss/casheer/pkg/casheerapi"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *handler) HandleGetRunningTotal(ctx *gin.Context) {

	var params casheerapi.GetTotalParams
	err := ctx.ShouldBindQuery(&params)

	if err != nil {
		common.EmitError(ctx, NewGetRunningTotalFailedError(
			http.StatusBadRequest,
			fmt.Sprintf("Could not bind query params: %s", err.Error())))
		return
	}

	var month = int(time.Now().Month())
	var year = time.Now().Year()

	if params.Month != nil {
		month = *params.Month
	}
	if params.Year != nil {
		year = *params.Year
	}

	var entries []model.Entry

	err = h.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		var dummyError error = errors.New("dummy")

		// Find all entries. (btw wtf is with that naming gorm??)
		txerr := tx.
			Preload("Expenses").
			Where("month = ?", month).
			Where("year = ?", year).
			Find(&entries).Error

		if txerr != nil {
			common.EmitError(ctx, NewGetRunningTotalFailedError(
				http.StatusInternalServerError,
				fmt.Sprintf("Could not list entries: %s", err.Error())))
			return dummyError
		}

		// If there are no entries for that month and year, return an error --
		// it likely means the month has not started yet.
		if len(entries) == 0 {
			common.EmitError(ctx, NewGetRunningTotalFailedError(
				http.StatusNotFound,
				fmt.Sprintf("No entires found for month %d and year %d. Check if period was started.", month, year)))
			return dummyError
		}
		return nil
	})

	// If there was an error, it was handled inside the transaction.
	if err != nil {
		return
	}

	// Could be done with a SUM clause, but would have to ignore salary and
	// I'd have to take a look at the syntax so fuck that.

	var expectedIncome, runningIncome, expectedTotal, runningTotal int64

	for _, e := range entries {

		// For all the incomes get their total value.
		if e.Category == "income" {
			expectedIncome += e.ExpectedTotal
			for _, exp := range e.Expenses {
				runningIncome += exp.Value
			}
		} else {
			expectedTotal += e.ExpectedTotal
			for _, exp := range e.Expenses {
				runningTotal += exp.Value
			}
		}
	}

	resp := casheerapi.GetTotalResponse{
		Data: casheerapi.TotalData{
			ResourceID: casheerapi.ResourceID{
				Id:   fmt.Sprintf("%d%d", month, year),
				Type: casheerapi.TotalType,
			},
			Month:          int(month),
			Year:           int(year),
			ExpectedIncome: expectedIncome,
			RunningIncome:  runningIncome,
			ExpectedTotal:  expectedTotal,
			RunningTotal:   runningTotal,
		},
	}
	ctx.JSON(http.StatusOK, &resp)
}
