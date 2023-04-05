package main

import (
	"database/sql"
	"fmt"
	"go-helper/internal/helper"
	"time"
)

// UPDATE table_name SET att2=T.att2 (VALUES ($1, $2), ($3, $4)) AS T (att1, att2) WHERE table_name.att1 = T.att1;
// UPDATE table_name SET att2=T.att2 (VALUES (?, ?)) AS T (att1, att2) WHERE table_name.att1 = T.att1;
// INSERT INTO table_name (a, b, c) VALUES (type($1, $2), $3, $4) ON ~
// DELETE FROM test_name WHERE a IN (type(?, ?), type(?, ?)) ON

func batchTest(sql string) error {
	batch, err := helper.NewBatch(sql)
	if err != nil {
		return err
	}
	fmt.Println("[Input]")
	fmt.Println(batch)

	for i := 0; i < 3; i++ {
		if err := batch.Push(i, i+1, i+2, i+3); err != nil {
			return err
		}
	}

	fmt.Println("[Output]")
	fmt.Println(batch)
	fmt.Println("Data       :", batch.Datas)

	return nil
}

func main2() {

	// str := `INSERT INTO table_name (a, b, c) VALUES (type($1, $2), $3, $4)
	// `
	// fmt.Println(str, "!!")

	// fmt.Println(strings.TrimSpace(str), "!!")

	t1 := time.Now()
	if err := batchTest(`INSERT INTO table_name (a, b, c) VALUES (type($1, $2), $3, $4)
		ON CONFLICT (a)
		DO UPDATE
		SET b = EXCLUDED.b, c = EXCLUDED.c`); err != nil {
		panic(err)
	}

	fmt.Println("----", time.Since(t1))

	if err := batchTest(`INSERT INTO table_name (a, b, c) VALUES (type(?, ?), ?, ?)
		on DUPLICATE KEY UPDATE b = VALUES(b), c = VALUES(c)`); err != nil {
		panic(err)
	}
	fmt.Println("----")

	if err := batchTest(`delete FROM test_name WHERE a in (type($1, $2, $3, $4))`); err != nil {
		panic(err)
	}
	fmt.Println("----")

	if err := batchTest(`DELETE FROM test_name WHERE a IN (type(?, ?), type(?, ?))`); err != nil {
		panic(err)
	}
	fmt.Println("----")

	if err := batchTest(`UPDATE table_name SET att2=T.att2 (VALUES ($1, $2, $3, $4)) AS T (att1, att2) WHERE table_name.att1 = T.att1`); err != nil {
		panic(err)
	}
	fmt.Println("----")
}

func example() error {
	db, err := sql.Open("postgres", "config")
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	defer tx.Rollback()

	bc, err := helper.NewBatch(`INSERT INTO table_name (a, b, c) VALUES (type($1, $2), $3, $4)
		ON CONFLICT (a)
		DO UPDATE
		SET b = EXCLUDED.b, c = EXCLUDED.c`)
	if err != nil {
		return err
	}

	fmt.Println("created batch \n", bc)

	if err := bc.Push(1, 2, 3, 4); err != nil {
		return err
	}

	stmt, err := db.Prepare(bc.SQL)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(bc.Datas...)
	if err != nil {
		return err
	}

	return tx.Commit()
}

/* example
tx, err := b.db.Begin()
defer tx.Rollback()

bc, err := NewBatch("INSERT INTO test(att1, att2) VALUES (?, ?)")
if err != nil {
	return err
}

for i := range datas {
	if err := bc.Push(datas[i].att1, datas[i].att2); err != nil {
		return err
	}
}

stmt, err := b.db.Prepare(bc.SQL)
defer stmt.Close()
if err != nil {
	return err
}
_, err = stmt.Exec(bc.Datas...)
if err != nil {
	return err
}

return tx.Commit()
// sql
// INSERT INTO table_name (a, b, c) VALUES (point($1, $2), $3, $4)
// or
// INSERT INTO table_name (a, b, c) VALUES (point(?, ?), ?, ?)
//
// result
// INSERT INTO table_name (a,b, c) VALUES (point($1, $2), $3, $4), (point($1, $2), $3, $4), (point($1, $2), $3, $4)
// or
// INSERT INTO table_name (a,b, c) VALUES (point(?, ?), ?, ?), (point(?, ?), ?, ?), (point(?, ?), ?, ?)
*/
