package purchases

import (
	"net/http"
	"strings"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/labstack/echo/v4"
	"github.com/mustafabalila/golang-api/db"
	"github.com/mustafabalila/golang-api/utils/common"
	"github.com/mustafabalila/golang-api/utils/logger"
)

// GetUnPaidPurchases queries the database for the list of purchases where
// status is "created" it accepts users (comma separated user ids) and
// created_at (cursor)
func GetUnPaidPurchases(c echo.Context) (e error) {
	logger := logger.GetLoggerInstance()
	var err error
	users := c.QueryParams().Get("users")
	createdAt := c.QueryParams().Get("createdAt")
	subscriptions := &[]db.PurchaseSubscription{}
	userId := c.Get("userId")

	query := db.Database.Model(subscriptions).
		Where("purchase_subscription.user_id = ?", userId).
		Where("status = ?", common.Statuses[common.CREATED]).
		Relation("Purchase").
		Relation("Purchase.User.full_name")

	if createdAt != "" {
		query.Where("purchase_subscription.created_at >= ?", createdAt)
	}

	if users != "" {
		userIds := strings.Split(users, ",")
		query.WhereGroup(func(q *orm.Query) (*orm.Query, error) {
			q = q.
				Where("purchase.user_id in (?) ", pg.In(userIds)).
				Where("purchase_subscription.user_id = ?", c.Get("userId"))
			return q, nil
		})

	}
	err = query.
		Order("created_at desc").
		Limit(10). // it's paginated using the created_at field (cursor)
		Select()
	if err != nil {
		logger.Error(err.Error())
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	response := map[string]interface{}{
		"subscriptions": subscriptions,
	}
	return c.JSON(http.StatusOK, response)
}
