package actions

import (
	"budget_tracker/models"
	"net/http"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	u := c.Value("current_user")
	if u != nil {
		user := u.(*models.User)
		tx := c.Value("tx").(*pop.Connection)
		now := time.Now()

		// Helper to get sum
		getSum := func(startDate, endDate time.Time, tType string) float64 {
			var res struct {
				Amount *float64 `db:"total"`
			}
			query := tx.RawQuery("SELECT SUM(amount) as total FROM transactions WHERE user_id = ? AND type = ? AND transaction_date >= ? AND transaction_date <= ?", user.ID, tType, startDate, endDate)
			if err := query.First(&res); err != nil {
				return 0
			}
			if res.Amount == nil {
				return 0
			}
			return *res.Amount
		}

		// Month Dates
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		endOfMonth := startOfMonth.AddDate(0, 1, -1)

		// Year Dates
		startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
		endOfYear := time.Date(now.Year(), 12, 31, 23, 59, 59, 0, now.Location())

		monthIncome := getSum(startOfMonth, endOfMonth, "Income")
		monthExpense := getSum(startOfMonth, endOfMonth, "Expense")
		yearIncome := getSum(startOfYear, endOfYear, "Income")
		yearExpense := getSum(startOfYear, endOfYear, "Expense")

		c.Set("monthIncome", monthIncome)
		c.Set("monthExpense", monthExpense)
		c.Set("monthMargin", monthIncome-monthExpense)

		c.Set("yearIncome", yearIncome)
		c.Set("yearExpense", yearExpense)
		c.Set("yearMargin", yearIncome-yearExpense)

		// Recent Transactions
		recent := &models.Transactions{}
		if err := tx.Where("user_id = ?", user.ID).Order("transaction_date desc").Limit(5).All(recent); err == nil {
			c.Set("recentTransactions", recent)
		}
	}

	return c.Render(http.StatusOK, r.HTML("home/index.plush.html"))
}
