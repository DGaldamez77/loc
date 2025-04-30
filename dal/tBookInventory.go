package dal

import (
	"encoding/json"
	"fmt"

	"github.com/dgaldamez77/loc/dal/dto"
	"github.com/joomcode/errorx"
)

func (dal DAL) GetBookInventory(params []QueryParams) (out *dto.BookInventory, err error) {
	where, qParams := dal.buildWhere(params, len(params) > 0)

	stmt := fmt.Sprintf(`
			-- GetBookInventory
			SELECT
				row_to_json(bi)
			FROM
				loc.book_inventory bi
			%s
			LIMIT 1`, where)

	var jBi []byte
	if err = dal.queryRow(stmt, qParams...).Scan(&jBi); err != nil {
		return nil, errorx.EnsureStackTrace(err)
	}

	if err = json.Unmarshal(jBi, &out); err != nil {
		return nil, errorx.EnsureStackTrace(err)
	}

	return out, nil
}

func (dal DAL) UpdateBookInventory(id int, changes map[string]interface{}) (out *dto.BookInventory, err error) {
	out = new(dto.BookInventory)
	err = dal.updateTable("loc.book_inventory", map[string]interface{}{"book_id": id}, changes, out)
	if err != nil {
		return nil, err
	}

	return out, nil
}
