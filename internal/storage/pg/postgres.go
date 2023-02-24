package pg

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Xrefullx/yandexDiplom2/internal/api/consta"
	"github.com/Xrefullx/yandexDiplom2/internal/models"
)

type PgStorage struct {
	connect *sql.DB
}

func New(uri string) (*PgStorage, error) {
	connect, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, err
	}
	return &PgStorage{connect: connect}, nil
}

func (PG *PgStorage) Ping() error {
	if err := PG.connect.Ping(); err != nil {
		return err
	}
	err := createTables(PG.connect)
	if err != nil {
		return err
	}
	return nil
}

func (PG *PgStorage) Close() error {
	if err := PG.connect.Close(); err != nil {
		return err
	}
	return nil
}

func createTables(connect *sql.DB) error {
	_, err := connect.Exec(`
	create table if not exists public.user(
		login text primary key,
		password text,
		createUser timestamp default now()
	);
	
	create table if not exists public.orders(
		 numberOrder text primary key,
		 login text,
		 statusOrder varchar(50),
		 accrualOrder double precision,
		 uploadedOrder timestamp default now(),
		 createdOrder timestamp default now(),
		 foreign key (login) references public.user (login)
	);
	
	create table if not exists public.withdraws(
		 login text,
		 numberOrder text,
		 sum double precision,
		 uploaded timestamp default now()
	);
	`)
	if err != nil {
		return err
	}
	return nil
}

func (PG *PgStorage) Adduser(ctx context.Context, user models.User) error {
	result, err := PG.connect.ExecContext(ctx,
		`insert into public.users (login, password) 
		values ($1, $2) on conflict do nothing`,
		user.Login, user.Password)
	if err != nil {
		return err
	}
	row, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if row == 0 || row == ' ' {
		return consta.ErrorNoUNIQUE
	}
	return nil
}

func (PG *PgStorage) Authentication(ctx context.Context, user models.User) (bool, error) {
	var done int
	err := PG.connect.QueryRowContext(ctx, `select count(1) from public.users where login=$1 and password=$2`,
		user.Login, user.Password).Scan(&done)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	}
	if done == 0 || done == ' ' {
		return false, nil
	}
	return true, nil
}

func (PG *PgStorage) GetOrder(ctx context.Context, numberOrder string) (models.Order, error) {
	var order models.Order
	err := PG.connect.QueryRowContext(ctx, `select numberOrder, login, statusOrder, accrualOrder, uploadedOrder
	from public.orders where number_order=$1 order by created_order desc`,
		numberOrder).Scan(&order.NumberOrder, &order.UserLogin, &order.Status,
		&order.Accrual, &order.Uploaded)
	if err != nil {
		return order, err
	}
	return order, nil
}

func (PG *PgStorage) GetOrders(ctx context.Context, userLogin string) ([]models.Order, error) {
	var orders []models.Order
	rows, err := PG.connect.QueryContext(ctx, `select numberOrder, login, statusOrder, accrualOrder, uploadedOrder
	from public.orders where login_user=$1`, userLogin)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.NumberOrder, &order.UserLogin, &order.Status, &order.Accrual, &order.Uploaded)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

func (PG *PgStorage) AddOrder(ctx context.Context, numberOrder string, order models.Order) error {
	result, err := PG.connect.ExecContext(ctx, `insert into public.orders 
    (numberOrder, login, statusOrder, uploadedOrder, accrualOrder)  values ($1, $2, $3, $4, $5) on conflict do nothing`,
		numberOrder, order.UserLogin, order.Status, order.Uploaded, order.Accrual)
	if err != nil {
		return err
	}
	row, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if row == 0 {
		return consta.ErrorNoUNIQUE
	}
	return nil
}

func (PG *PgStorage) UpdateOrder(ctx context.Context, loyaltyPoint models.Loyalty) error {
	_, err := PG.connect.ExecContext(ctx, `update public.orders set accrualOrder=$1, statusOrder=$2,
                         uploadedOrder=now() where numberOrder =$3`,
		loyaltyPoint.Accrual, loyaltyPoint.Status, loyaltyPoint.NumberOrder)
	if err != nil {
		return err
	}
	return nil
}

func (PG *PgStorage) GetOrdersProcess(ctx context.Context) ([]models.Order, error) {
	var orders []models.Order
	sliceStatus := []interface{}{consta.OrderStatusPROCESSING, consta.OrderStatusNEW, consta.OrderStatusREGISTERED, consta.OrderStatusInvalid}
	rows, err := PG.connect.QueryContext(ctx, `select numberOrder, login, statusOrder, accrualOrder, uploadedOrder
	from public.orders where statusOrder in ($1, $2, $3,$4)`, sliceStatus...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var order models.Order
		err = rows.Scan(&order.NumberOrder, &order.UserLogin, &order.Status, &order.Accrual, &order.Uploaded)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, rows.Err()
}

func (PG *PgStorage) GetUserBalance(ctx context.Context, userLogin string) (float64, float64, error) {
	var ordersSUM float64
	var withdrawsSUM float64
	err := PG.connect.QueryRowContext(ctx, `select (case when sumOrder is null then 0 else sum_order end) as sum_order, (case when sum_withdraws is null then 0 else sum_withdraws end) as sum_withdraws from
	 (select sum(accrualOrder) as  sum_order from public.orders where login = $1) as orders,
	 (select sum(sum) as  sum_withdraws from public.withdraws where login = $1) as withdraws`, userLogin).
		Scan(&ordersSUM, &withdrawsSUM)
	return ordersSUM, withdrawsSUM, err
}

func (PG *PgStorage) AddWithdraw(ctx context.Context, withdraw models.Withdraw) error {
	result, err := PG.connect.ExecContext(ctx, `
	insert into public.withdraws (login, numberOrder, sum, uploadedOrder)
	select $1, $2, $3, $4
	where (
          select sumOrder >= sumWithdraws + $3 from (
          select (case when sumOrder is null then 0 else sumOrder end ) as sumOrder,
          (case when sumWithdraws is null then 0 else sumWithdraws end ) as sumWithdraws from
          (select sum(accrualOrder) as  sum_order from public.orders where login = $1) as orders,
          (select sum(sum) as  sumWithdraws from public.withdraws where login = $1) as withdraws) as s
          );
	`, withdraw.UserLogin, withdraw.NumberOrder, withdraw.Sum, withdraw.ProcessedAT)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return consta.ErrorStatusShortfallAccount
	}
	return nil
}

func (PG *PgStorage) GetWithdraws(ctx context.Context, userLogin string) ([]models.Withdraw, error) {
	var withdraws []models.Withdraw
	rows, err := PG.connect.QueryContext(ctx, `select login, numberOrder, sum, uploadedOrder from public.withdraws
	where login = $1
	order by uploadedOrder`, userLogin)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	for rows.Next() {
		var withdraw models.Withdraw
		err = rows.Scan(&withdraw.UserLogin, &withdraw.NumberOrder, &withdraw.Sum, &withdraw.ProcessedAT)
		if err != nil {
			return nil, err
		}
		withdraws = append(withdraws, withdraw)
	}
	return withdraws, rows.Err()
}
