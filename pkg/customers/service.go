package customers

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v4/pgxpool"
)

//ErrNotFound возвращается, когда покупатель не найден.
var ErrNotFound = errors.New("item not found")

//ErrInternal возвращается, когда произошла внутернняя ошибка.
var ErrInternal = errors.New("internal error")

//Service описывает сервис работы с покупателям.
type Service struct {
	pool *pgxpool.Pool
}

//NewService создаёт сервис.
func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

type Quotes struct {
	ID       int64  `json:"id"`
	Author   string `json:"author"`
	Quote    string `json:"quote"`
	Category string `json:"category"`
}

//All ....
func (s *Service) All(ctx context.Context) (quotes []*Quotes, err error) {

	sqlStatement := `SELECT * FROM quote`

	rows, err := s.pool.Query(ctx, sqlStatement)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		item := &Quotes{}
		err := rows.Scan(
			&item.ID,
			&item.Category,
			&item.Quote,
			&item.Author,
		)
		if err != nil {
			log.Println(err)
		}

		quotes = append(quotes, item)
	}

	return quotes, nil
}

//ByID ...
func (s *Service) Bycategory(ctx context.Context, category string) (quotes []*Quotes, err error) {
	sqlStatement := `SELECT * FROM quote WHERE category = $1`
	rows, err := s.pool.Query(ctx, sqlStatement, category)
	fmt.Println(category)
	defer rows.Close()

	for rows.Next() {
		item := &Quotes{}
		err := rows.Scan(
			&item.ID,
			&item.Category,
			&item.Quote,
			&item.Author,
		)
		if err != nil {
			log.Println(err)
		}

		quotes = append(quotes, item)
	}

	return quotes, nil
}

//Delete ...
func (s *Service) Delete(ctx context.Context, id int64) (*Quotes, error) {
	item := &Quotes{}
	fmt.Println(id)
	sqlStatement := `DELETE FROM quote WHERE id = $1 RETURNING *`
	err := s.pool.QueryRow(ctx, sqlStatement, id).Scan(
		&item.ID,
		&item.Category,
		&item.Quote,
		&item.Author,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil
}

//Save ...
func (s *Service) Save(ctx context.Context, quote *Quotes) (c *Quotes, err error) {

	item := &Quotes{}

	sqlStatement := `INSERT INTO quote(author, quote, category) VALUES($1, $2, $3) RETURNING *`
	err = s.pool.QueryRow(ctx, sqlStatement, quote.Author, quote.Quote, quote.Category).Scan(
		&item.ID,
		&item.Category,
		&item.Quote,
		&item.Author,
	)

	if err != nil {
		log.Print(err)
		return nil, ErrInternal
	}

	return item, nil
}
