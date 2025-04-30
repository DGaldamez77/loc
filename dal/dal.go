package dal

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/dgaldamez77/oloc/dal/dto"
	"github.com/lib/pq"
)

type IDAL interface {
	query(string, ...interface{}) (*sql.Rows, error)
	queryRow(string, ...interface{}) *sql.Row

	GetBooks([]QueryParams) ([]dto.Book, error)
	GetBooksUsingAuthor([]QueryParams) ([]dto.Book, error)
}

type DAL struct {
	DB *sql.DB
	Tx *sql.Tx
}

func NewDAL() IDAL {
	db, err := sql.Open("postgres", os.Getenv("DB_URI"))
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute * 60)
	db.SetConnMaxIdleTime(time.Second * 300)

	ctx, cf := context.WithTimeout(context.Background(), time.Second*60)
	defer cf()

	err = db.PingContext(ctx)
	if err != nil {
		panic("Database ping failed: " + err.Error())
	}

	return &DAL{DB: db}
}

type QueryParams struct {
	FieldName           string
	ComparisonFieldName string
	Value               interface{}
	OmitIfEmpty         bool
	Comparison          string

	StartGroup      bool
	EndGroup        bool
	BooleanOperator string
}

func (dal DAL) buildWhere(params []QueryParams, extra ...interface{}) (where string, queryParams []interface{}) {
	var conditions []string
	addWhere := false
	addOperatorIfNoFields := ""

	if len(extra) > 0 {
		// nolint:gocritic
		switch e := extra[0].(type) {
		case bool:
			addWhere = e
		}
	}

	if len(extra) > 1 {
		// nolint:gocritic
		switch e := extra[1].(type) {
		case string:
			addOperatorIfNoFields = e
		}
	}

	i := 1
	for j, p := range params {
		if p.OmitIfEmpty && p.Value == "" {
			continue
		}

		if p.Comparison == "" {
			p.Comparison = "="
		}

		if p.BooleanOperator == "" {
			p.BooleanOperator = "AND"
		}

		// Do not add a bool operator when this is the last query parameter
		if j == len(params)-1 {
			p.BooleanOperator = ""
		}

		var openParenthesis string
		var closeParenthesis string

		if p.StartGroup {
			openParenthesis = "("
		}

		if p.EndGroup {
			closeParenthesis = ")"
		}

		switch strings.ToUpper(p.Comparison) {
		case "IS NULL", "IS NOT NULL":
			conditions = append(conditions, fmt.Sprintf("%s%s %s%s %s ", openParenthesis, p.FieldName, p.Comparison, closeParenthesis, p.BooleanOperator))
		case "ANY":
			conditions = append(conditions, fmt.Sprintf("%s%s = %s($%d)%s %s ", openParenthesis, p.FieldName, p.Comparison, i, closeParenthesis, p.BooleanOperator))
			queryParams = append(queryParams, pq.Array(p.Value))
			i++
		case "NOT IN", "IN":
			t := reflect.TypeOf(p.Value).Kind()
			var indexes string
			if t == reflect.Slice {
				s := reflect.ValueOf(p.Value)
				for idx := 0; idx < s.Len(); idx++ {
					fmt.Println(s.Index(idx))
					if idx == 0 {
						indexes += fmt.Sprintf("$%d", i)
					} else {
						indexes += fmt.Sprintf(",$%d", i)
					}
					queryParams = append(queryParams, s.Index(idx).Interface())
					i++
				}
			} else {
				s := reflect.ValueOf(p.Value)
				indexes += fmt.Sprintf("$%d", i)
				queryParams = append(queryParams, s.Interface())
				i++
			}
			conditions = append(conditions, fmt.Sprintf("%s%s %s (%s)%s %s ", openParenthesis, p.FieldName,
				p.Comparison, indexes, closeParenthesis, p.BooleanOperator))
		default:
			if p.ComparisonFieldName != "" {
				conditions = append(conditions, fmt.Sprintf("%s%s %s %s%s %s ", openParenthesis, p.FieldName,
					p.Comparison, p.ComparisonFieldName, closeParenthesis, p.BooleanOperator))
			} else {
				conditions = append(conditions, fmt.Sprintf("%s%s %s $%d%s %s ", openParenthesis, p.FieldName,
					p.Comparison, i, closeParenthesis, p.BooleanOperator))
				queryParams = append(queryParams, p.Value)
				i++
			}
		}
	}

	if addWhere {
		where = " WHERE "
	}

	if addOperatorIfNoFields != "" {
		where = " " + addOperatorIfNoFields + " "
	}

	if i > 1 {
		where += strings.Join(conditions, "")
	}

	return where, queryParams
}

func (dal DAL) queryRow(query string, args ...interface{}) *sql.Row {
	if dal.Tx != nil {
		return dal.Tx.QueryRow(query, args...)
	}

	return dal.DB.QueryRow(query, args...)
}

func (dal DAL) query(query string, args ...interface{}) (*sql.Rows, error) {
	if dal.Tx != nil {
		return dal.DB.Query(query, args...)
	}

	return dal.DB.Query(query, args...)
}
