package database

import (
	"database/sql"
	"testing"

	"clean-architecture/internal/entity"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
)

type OrderRepositoryTestSuite struct {
	suite.Suite
	Db *sql.DB
}

func (suite *OrderRepositoryTestSuite) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	suite.NoError(err)
	db.Exec("CREATE TABLE orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))")
	suite.Db = db
}

func (suite *OrderRepositoryTestSuite) TearDownTest() {
	suite.Db.Close()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(OrderRepositoryTestSuite))
}

func (suite *OrderRepositoryTestSuite) TestGivenOrders_WhenList_ThenShouldReturnAllOrders() {
	order1, err := entity.NewOrder("123", 10.0, 2.0)
	suite.NoError(err)
	suite.NoError(order1.CalculateFinalPrice())
	order2, err := entity.NewOrder("456", 20.0, 4.0)
	suite.NoError(err)
	suite.NoError(order2.CalculateFinalPrice())

	repo := NewOrderRepository(suite.Db)
	err = repo.Save(order1)
	suite.NoError(err)
	err = repo.Save(order2)
	suite.NoError(err)

	orders, err := repo.List()
	suite.NoError(err)
	suite.Len(orders, 2)

	suite.Equal(order1.ID, orders[1].ID)
	suite.Equal(order1.Price, orders[1].Price)
	suite.Equal(order1.Tax, orders[1].Tax)
	suite.Equal(order1.FinalPrice, orders[1].FinalPrice)

	suite.Equal(order2.ID, orders[0].ID)
	suite.Equal(order2.Price, orders[0].Price)
	suite.Equal(order2.Tax, orders[0].Tax)
	suite.Equal(order2.FinalPrice, orders[0].FinalPrice)
}

func (suite *OrderRepositoryTestSuite) TestGivenAnOrder_WhenSave_ThenShouldSaveOrder() {
	order, err := entity.NewOrder("123", 10.0, 2.0)
	suite.NoError(err)
	suite.NoError(order.CalculateFinalPrice())
	repo := NewOrderRepository(suite.Db)
	err = repo.Save(order)
	suite.NoError(err)

	var orderResult entity.Order
	err = suite.Db.QueryRow("SELECT id, price, tax, final_price FROM orders WHERE id = ?", order.ID).
		Scan(&orderResult.ID, &orderResult.Price, &orderResult.Tax, &orderResult.FinalPrice)

	suite.NoError(err)
	suite.Equal(order.ID, orderResult.ID)
	suite.Equal(order.Price, orderResult.Price)
	suite.Equal(order.Tax, orderResult.Tax)
	suite.Equal(order.FinalPrice, orderResult.FinalPrice)
}
