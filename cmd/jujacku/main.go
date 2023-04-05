package main

import (
	"database/sql"
	"fmt"
	"go-helper/cmd/jujacku/jujacku"
	"go-helper/internal/helper"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
)
type Report struct {
	Current int
	Want int
	Rand int
}
func (r Report) String() string {
	return fmt.Sprintf("Current : %d\nWand    : %d\nRand    : %d", r.Current, r.Want, r.Rand)
}

func main() {

	ret := make(map[int]Report)

	config, err := jujacku.NewConfig()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("[Config]\n", config)

	dbinfo := helper.DBInfo{
		Host: "10.101.0.17",
		User: "tas",
		Password: "nvidia",
		Database: "sctm",
		Port: 3306,
	}

	db, err := sql.Open("mysql", dbinfo.Config("mysql"))
	if err != nil {
		fmt.Println(err)
		return
	}
	
	rows1, err := db.Query(`SELECT A.AvenueSeq, count(*) FROM CollectionData AS A 
		JOIN Avenue AS B 
		ON A.AvenueSeq = B.AvenueSeq
		WHERE CollectedDate BETWEEN ? AND ?
		AND B.IntersectionSeq = ?
		AND A.MovementType = ?
		GROUP BY A.AvenueSeq`, helper.ToUTCDatetime(config.Start), helper.ToUTCDatetime(config.End), config.InterID, config.Flow)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	for rows1.Next() {
		av := 0
		cnt := 0
		if err := rows1.Scan(&av, &cnt); err != nil {
			fmt.Println(err)
			return
		}
		ret[av] = Report{
			Current: cnt,
		}
		fmt.Printf("Avenue %d : %d\n", av, cnt)
	}

	for k, v := range ret {
		v.Want = v.Current * int(config.Percentage) / 100
		ret[k] = v
	}

	rows2, err := db.Query(`SELECT A.AvenueSeq, A.DataSeq FROM CollectionData AS A 
		JOIN Avenue AS B 
		ON A.AvenueSeq = B.AvenueSeq
		WHERE CollectedDate BETWEEN ? AND ?
		AND B.IntersectionSeq = ?
		AND A.MovementType = ?`, helper.ToUTCDatetime(config.Start), helper.ToUTCDatetime(config.End), config.InterID, config.Flow)
	if err != nil {
		fmt.Println(err)
		return
	}

	rand.Seed(time.Now().UnixNano())
	
	for rows2.Next() {
		av := 0
		seq := 0
		if err := rows2.Scan(&av, &seq); err != nil {
			fmt.Println(err)
			return
		}
		succ := rand.Intn(100) < int(config.Percentage)
		succstr := "true"
		if !succ {
			succstr = "false"
		}
		fmt.Printf("Seq %d : %s\n", seq, succstr)
		if succ {
			v, exist := ret[av]
			if !exist {
				continue
			}
			v.Rand += 1
			ret[av] = v
		}
	}

	fmt.Println("Report --")
	for k, v := range ret {
		fmt.Println("[", k, "]")
		fmt.Println(v)
	}

	fmt.Println("done.")
}