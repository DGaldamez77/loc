package dal

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/dgaldamez77/loc/dal/dto"
	"github.com/dgaldamez77/loc/log"
	"github.com/joomcode/errorx"
)

func (dal DAL) GetBooks(params []QueryParams) (out []dto.Book, err error) {
	where, qParams := dal.buildWhere(params, len(params) > 0)

	stmt := fmt.Sprintf(`
			-- GetBooks
			SELECT
				row_to_json(b),
				row_to_json(p),
				row_to_json(g),
				(	SELECT 
						array_to_json(array_agg(row_to_json(a))) 
					FROM
						loc.book_author_xref bax
							INNER JOIN loc.author a ON a.author_id = bax.author_id
					WHERE 
						b.book_id = bax.book_id
				) authors
			FROM
				loc.book b
				INNER JOIN loc.publisher p ON p.publisher_id = b.publisher_id
				INNER JOIN loc.genre g ON g.genre_id = b.genre_id
			%s`, where)

	var jB, jP, jG, jAuthors []byte
	var rows *sql.Rows
	if rows, err = dal.query(stmt, qParams...); err != nil {
		return nil, errorx.EnsureStackTrace(err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.LogError(err, "failed to close rows")
		}
	}()

	for rows.Next() {
		var b dto.Book

		if err2 := rows.Scan(&jB, &jP, &jG, &jAuthors); err2 != nil {
			return nil, err2
		}

		if err = json.Unmarshal(jB, &b); err != nil {
			return nil, errorx.EnsureStackTrace(err)
		}

		if err = json.Unmarshal(jP, &b.Publisher); err != nil {
			return nil, errorx.EnsureStackTrace(err)
		}

		if err = json.Unmarshal(jG, &b.Genre); err != nil {
			return nil, errorx.EnsureStackTrace(err)
		}

		if err = json.Unmarshal(jAuthors, &b.Authors); err != nil {
			return nil, errorx.EnsureStackTrace(err)
		}

		out = append(out, b)
	}

	return out, nil
}

func (dal DAL) GetBooksUsingAuthor(params []QueryParams) (out []dto.Book, err error) {
	where, qParams := dal.buildWhere(params, len(params) > 0)

	stmt := fmt.Sprintf(`
			-- GetBooks
			SELECT
				row_to_json(b),
				row_to_json(p),
				row_to_json(g),
				(	SELECT 
						array_to_json(array_agg(row_to_json(a))) 
					FROM
						loc.book_author_xref bax
							INNER JOIN loc.author a ON a.author_id = bax.author_id
					WHERE 
						b.book_id = bax.book_id
				) authors
			FROM
				loc.book b
				INNER JOIN loc.publisher p ON p.publisher_id = b.publisher_id
				INNER JOIN loc.genre g ON g.genre_id = b.genre_id
				INNER JOIN loc.book_author_xref bax ON b.book_id = bax.book_id
				INNER JOIN loc.author a ON bax.author_id = a.author_id

			%s`, where)

	var jB, jP, jG, jAuthors []byte
	var rows *sql.Rows
	if rows, err = dal.query(stmt, qParams...); err != nil {
		return nil, errorx.EnsureStackTrace(err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.LogError(err, "failed to close rows")
		}
	}()

	for rows.Next() {
		var b dto.Book

		if err2 := rows.Scan(&jB, &jP, &jG, &jAuthors); err2 != nil {
			return nil, err2
		}

		if err = json.Unmarshal(jB, &b); err != nil {
			return nil, errorx.EnsureStackTrace(err)
		}

		if err = json.Unmarshal(jP, &b.Publisher); err != nil {
			return nil, errorx.EnsureStackTrace(err)
		}

		if err = json.Unmarshal(jG, &b.Genre); err != nil {
			return nil, errorx.EnsureStackTrace(err)
		}

		if err = json.Unmarshal(jAuthors, &b.Authors); err != nil {
			return nil, errorx.EnsureStackTrace(err)
		}

		out = append(out, b)
	}

	return out, nil
}

func (dal DAL) GetBook(params []QueryParams) (out *dto.Book, err error) {
	where, qParams := dal.buildWhere(params, len(params) > 0)

	stmt := fmt.Sprintf(`
			-- GetBook
			SELECT
				row_to_json(b),
				row_to_json(p),
				row_to_json(g),
				(	SELECT 
						array_to_json(array_agg(row_to_json(a))) 
					FROM
						loc.book_author_xref bax
							INNER JOIN loc.author a ON a.author_id = bax.author_id
					WHERE 
						b.book_id = bax.book_id
				) authors
			FROM
				loc.book b
				INNER JOIN loc.publisher p ON p.publisher_id = b.publisher_id
				INNER JOIN loc.genre g ON g.genre_id = b.genre_id
			%s
			LIMIT 1`, where)

	var jB, jP, jG, jAuthors []byte
	if err = dal.queryRow(stmt, qParams...).Scan(&jB, &jP, &jG, &jAuthors); err != nil {
		return nil, errorx.EnsureStackTrace(err)
	}

	if err = json.Unmarshal(jB, &out); err != nil {
		return nil, errorx.EnsureStackTrace(err)
	}

	if err = json.Unmarshal(jP, &out.Publisher); err != nil {
		return nil, errorx.EnsureStackTrace(err)
	}

	if err = json.Unmarshal(jG, &out.Genre); err != nil {
		return nil, errorx.EnsureStackTrace(err)
	}

	if err = json.Unmarshal(jAuthors, &out.Authors); err != nil {
		return nil, errorx.EnsureStackTrace(err)
	}

	return out, nil
}
