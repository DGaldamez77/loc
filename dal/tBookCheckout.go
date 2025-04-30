package dal

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/dgaldamez77/loc/dal/dto"
	"github.com/dgaldamez77/loc/util"
	"github.com/joomcode/errorx"
)

func (dal DAL) GetBookCheckouts(params []QueryParams) (out []dto.BookCheckout, err error) {
	where, qParams := dal.buildWhere(params, len(params) > 0)

	stmt := fmt.Sprintf(`
			-- GetBookCheckouts
			SELECT
				row_to_json(bco)
			FROM
				loc.book_checkout bco
			%s`, where)

	var jBCO []byte
	var rows *sql.Rows
	if rows, err = dal.query(stmt, qParams...); err != nil {
		return nil, errorx.EnsureStackTrace(err)
	}
	defer rows.Close()

	for rows.Next() {
		var bco dto.BookCheckout

		if err2 := rows.Scan(&jBCO); err2 != nil {
			return nil, err2
		}

		if err = json.Unmarshal(jBCO, &bco); err != nil {
			return nil, errorx.EnsureStackTrace(err)
		}

		out = append(out, bco)
	}

	return out, nil
}

func (dal DAL) InsertBookCheckout(in dto.BookCheckout) (out *dto.BookCheckout, err error) {
	if in.CreatedBy == "" {
		return nil, util.ErrMissingCreatedBy
	}

	var j []byte
	stmt := `
			-- InsertBookCheckout
			INSERT INTO loc.book_checkout
				(book_id, user_id, created_by)
			VALUES ($1,$2, $3)
			RETURNING row_to_json(book_checkout)`
	if err = dal.queryRow(stmt, in.BookID, in.UserID, in.CreatedBy).Scan(&j); err != nil {
		return nil, errorx.EnsureStackTrace(err)
	}

	out = new(dto.BookCheckout)
	if err = json.Unmarshal(j, &out); err != nil {
		return nil, errorx.EnsureStackTrace(err)
	}

	return out, nil
}
