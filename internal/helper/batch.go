package helper

import (
	"fmt"
	"strconv"
	"strings"
)

type DSN = string

const (
	MySQL      DSN = "mysql"
	PostgreSQL DSN = "postgres"
)

type Method = string

const (
	Delete Method = "DELETE"
	Insert Method = "INSERT"
	Update Method = "UPDATE"
)

// Analyze the sql to figure out the dsn, and return the number of arguments.
// return : dsn name, argment count, error
func FindDSN(sql string) (DSN, int, error) {
	mysqlcnt := strings.Count(sql, "?")
	postgrescnt := strings.Count(sql, "$")
	if mysqlcnt != 0 && postgrescnt != 0 {
		return "", 0, fmt.Errorf("ambiguous symbol, ?:%d, $:%d", mysqlcnt, postgrescnt)
	}

	if mysqlcnt != 0 {
		return MySQL, mysqlcnt, nil
	} else if postgrescnt != 0 {
		return PostgreSQL, postgrescnt, nil
	} else {
		return "", 0, fmt.Errorf("not found placeholder ? or $")
	}
}

// find query method, INSERT or DELETE
func FindMethod(sql string) (Method, error) {
	if strings.Contains(sql, Delete) || strings.Contains(sql, strings.ToLower(Delete)) {
		return Delete, nil
	} else if strings.Contains(sql, Insert) || strings.Contains(sql, strings.ToLower(Insert)) {
		return Insert, nil
	} else if strings.Contains(sql, Update) || strings.Contains(sql, strings.ToLower(Update)) {
		return Update, nil
	} else {
		return "", fmt.Errorf("not found method")
	}
}

// FIND 'ON ~' TEXT
// return header, description
func FindDescription(sql string, method Method) (string, string) {

	if method == Insert || method == Delete {
		str, sub := splitLower(sql, " ON ")
		if len := len(str); len == 1 {
			return str[0], ""
		} else if len == 0 {
			return "", ""
		} else {
			return str[0], " " + sub + str[1]
		}
	}

	str, sub := splitLower(sql, " AS ")
	if len := len(str); len == 1 {
		return str[0], ""
	} else if len == 0 {
		return "", ""
	} else {
		return str[0], " " + sub + str[1]
	}
}

// remove \n and \r
func Serialize(s string) string {
	s = strings.ReplaceAll(s, string(rune(10)), " ")
	s = strings.ReplaceAll(s, string(rune(9)), "")
	s = strings.ReplaceAll(s, "  ", " ")

	return s
}

// TrimRight only once
func trimRightOnce(s string, cutset string) string {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == ' ' {
			continue
		} else if s[i] == ')' {
			return s[0:i]
		} else {
			break
		}
	}

	return s
}

// You must always put sub in lowercase.
// You split when s is uppercase or lowercase.
func splitLower(s string, sub string) ([]string, string) {
	if strings.Contains(s, sub) {
		return strings.Split(s, sub), sub
	} else {
		lower := strings.ToLower(sub)
		return strings.Split(s, lower), lower
	}
}

// example:
//
// func example() error {
// 	db, err := sql.Open("postgres", "config")
// 	if err != nil {
// 		return err
// 	}
// 	tx, err := db.Begin()
// 	defer tx.Rollback()

// 	bc, err := helper.NewBatch(`INSERT INTO table_name (a, b, c) VALUES (type($1, $2), $3, $4)
// 		ON CONFLICT (a)
// 		DO UPDATE
// 		SET b = EXCLUDED.b, c = EXCLUDED.c`)
// 	if err != nil {
// 		return err
// 	}

// 	fmt.Println("created batch \n", bc)

// 	if err := bc.Push(1, 2, 3, 4); err != nil {
// 		return err
// 	}

// 	stmt, err := db.Prepare(bc.SQL)
// 	if err != nil {
// 		return err
// 	}

// 	_, err = stmt.Exec(bc.Datas...)
// 	if err != nil {
// 		return err
// 	}

//		return tx.Commit()
//	}
type Batch struct {
	SQL    string // INSERT INTO TABLE_NAME(ATT1, ATT2) VALUES (?, ?, ?, ?)
	Datas  []any  // x4
	Length int    // 4
	DSN    DSN    // postgres, mysql
	Method Method // DELETE, INSERT
	Desc   string
	parser parser
}

func NewBatch(sql string) (Batch, error) {
	dsn, placeholders, err := FindDSN(sql)
	if err != nil {
		return Batch{}, fmt.Errorf("NewBatch fail : %s", err)
	}

	method, err := FindMethod(sql)
	if err != nil {
		return Batch{}, fmt.Errorf("NewBatch fail : %s", err)
	}

	sql = Serialize(sql)

	nsql, desc := FindDescription(sql, method)
	nsql = strings.TrimSpace(nsql)

	var psr parser
	switch dsn {
	case MySQL:
		psr = &mysqlParser{
			sql:          nsql,
			method:       method,
			placeholders: placeholders,
			datas:        0,
		}
	case PostgreSQL:
		psr = &postgresParser{
			sql:          nsql,
			method:       method,
			placeholders: placeholders,
			datas:        0,
		}
	default:
		return Batch{}, fmt.Errorf("NewBatch fail : unsupported dsn %s", dsn)
	}

	if err := psr.parse(); err != nil {
		return Batch{}, fmt.Errorf("NewBatch fail : %s", err)
	}

	return Batch{
		SQL:    sql,
		Length: placeholders,
		DSN:    dsn,
		Desc:   desc,
		parser: psr,
		Method: method,
		Datas:  []any{},
	}, nil
}

func (b Batch) String() string {
	return fmt.Sprintf("SQL        : %s\nlen(Datas) : %d\nLength     : %d\nDSN        : %s\nMethod     : %s\n%s\nDesc       : %s",
		b.SQL, len(b.Datas), b.Length, b.DSN, b.Method, b.parser, b.Desc)
}

// Batch.Datas length
func (b *Batch) Len() int {
	return len(b.Datas)
}

// append Batch.Datas
// modifIed b.SQL
func (b *Batch) Push(arg ...any) error {
	argments := len(arg)
	if argments != b.Length {
		return fmt.Errorf("Batch.Push unmatched placeholder %d, %d", b.Length, argments)
	}

	sql, err := b.parser.push()
	if err != nil {
		return err
	}

	b.Datas = append(b.Datas, arg...)
	b.SQL = sql + b.Desc

	return nil
}

type parser interface {
	parse() error
	push() (string, error)
	String() string
}

type mysqlParser struct {
	sql          string
	method       Method
	placeholders int
	datas        int
	tail         string
}

func (m mysqlParser) String() string {
	return fmt.Sprintf("Tail       : %s", m.tail)
}

func (m *mysqlParser) parse() error {
	if m.method == Delete {
		tails, _ := splitLower(m.sql, "IN")
		if len(tails) != 2 {
			return fmt.Errorf("Batch.mysqlParser error : 'IN' split fail")
		}
		tail := strings.TrimSpace(tails[1])
		tail = strings.TrimLeft(tail, "(")
		m.tail = trimRightOnce(tail, ")")
	} else if m.method == Insert {
		tails, _ := splitLower(m.sql, "VALUES")
		if len(tails) != 2 {
			return fmt.Errorf("Batch.mysqlParser error : 'VALUES' split fail")
		}

		m.tail = strings.TrimSpace(tails[1])
	} else {
		return fmt.Errorf("Batch.mysqlParser error : method not found %s", m.method)
	}

	return nil
}
func (m *mysqlParser) push() (string, error) {
	m.datas += m.placeholders
	if m.datas == 0 {
		return m.sql, nil
	}

	if m.method == Delete {
		m.sql = trimRightOnce(m.sql, ")") + ", " + m.tail + ")"
	} else if m.method == Insert {
		m.sql += ", " + m.tail
	}

	return m.sql, nil
}

type postgresParser struct {
	sql          string
	method       Method
	placeholders int
	datas        int
	tail         string
}

func (m postgresParser) String() string {
	return fmt.Sprintf("Tail       : %s", m.tail)
}

// (point($11, $22), $3, $4) -> (point(%s, %s), %s, %s)
func replaceSprintf(origin string) string {
	ret := ""
	isConvert := false

	for _, c := range origin {
		if c == '$' {
			ret += "%s"
			isConvert = true
		} else if isConvert {
			if c == ',' || c == ')' {
				isConvert = false
				ret += string(c)
			} else {
				continue
			}
		} else {
			ret += string(c)
		}
	}

	return ret
}

func (m *postgresParser) parse() error {
	if m.method == Delete {
		tails, _ := splitLower(m.sql, "IN")
		if len(tails) != 2 {
			return fmt.Errorf("Batch.postgresParser error : 'IN' split fail")
		}
		tail := strings.TrimSpace(tails[1])
		tail = strings.TrimLeft(tail, "(")
		tail = trimRightOnce(tail, ")")

		m.tail = replaceSprintf(tail)
	} else if m.method == Insert {
		tails, _ := splitLower(m.sql, "VALUES")
		if len(tails) != 2 {
			return fmt.Errorf("Batch.postgresParser error : 'VALUES' split fail : %s", tails)
		}
		tail := strings.TrimSpace(tails[1])

		m.tail = replaceSprintf(tail)
	} else if m.method == Update {
		tails, _ := splitLower(m.sql, "VALUES")
		if len(tails) != 2 {
			return fmt.Errorf("Batch.mysqlParser error : 'VALUES' split fail")
		}

		tail := strings.TrimSpace(tails[1])
		tail = trimRightOnce(tail, ")")

		m.tail = replaceSprintf(tail)

	} else {
		return fmt.Errorf("Batch.postgresParser error : method not found %s", m.method)
	}

	return nil
}
func (m *postgresParser) push() (string, error) {
	if m.datas == 0 {
		m.datas += m.placeholders
		return m.sql, nil
	}

	placeholder := []any{}
	for i := 0; i < m.placeholders; i++ {
		placeholder = append(placeholder, "$"+strconv.Itoa(i+m.datas+1))
	}

	tail := fmt.Sprintf(m.tail, placeholder...)

	if m.method == Delete || m.method == Update {
		m.sql = trimRightOnce(m.sql, ")") + ", " + tail + ")"
	} else if m.method == Insert {
		m.sql += ", " + tail
	}

	m.datas += m.placeholders

	return m.sql, nil
}
