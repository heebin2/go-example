package main

import (
	"fmt"
	"go-helper/internal/postgrecdc/postgrecdc"
)

func main() {

	pcdc := postgrecdc.PostgreCDC{
		DebeziumAddress:    "localhost:8083",
		DatabaseHost:       "postgres",
		DatabasePort:       "5432",
		DatbaseUser:        "tms",
		DatabasePassword:   "nvidia",
		DatabaseName:       "tms",
		DatabaseWatchTable: []string{"test", "test2"},
	}
	conList, err := pcdc.ConnectorList()
	if err != nil {
		fmt.Println("get connector list error : ", err)
		return
	}
	fmt.Println("connectors : ", conList)

	for _, v := range conList {
		pcdc.UnregistOther(v)
	}

	conList, err = pcdc.ConnectorList()
	if err != nil {
		fmt.Println("get connector list error : ", err)
		return
	}
	fmt.Println("connectors : ", conList)

	isRegisted, err := pcdc.IsRegisted()
	if err != nil {
		fmt.Println("registed checking error : ", err)
		return
	}

	if isRegisted {
		fmt.Println("connector is already registed. try unregist.")
		if err := pcdc.Unregist(); err != nil {
			fmt.Println("unregisted error : ", err)
			return
		}
		fmt.Println("connector is unregist success.")
	}

	fmt.Println("server name : ", pcdc.ServerName())

	if err := pcdc.Regist(); err != nil {
		fmt.Println("regist error : ", err)
		return
	}

	p, err := pcdc.RemoteConfig()
	fmt.Println(p)

	topics := pcdc.Topics()
	if len(topics) <= 0 {
		fmt.Println("not found topics")
		return
	}
	fmt.Println("topics : ", topics)

	// reader := kafka.NewReader()
	// if err := reader.Open("localhost:9092,localhost:9093,localhost:9094", "postgrecdc", topics[0]); err != nil {
	// 	fmt.Println(err)
	// }

	// defer reader.Close()

	for {
		// message := reader.Read()
		// if err != nil {
		// 	fmt.Println("err : ", err)
		// }
		// if len(message) > 0 {
		// 	fmt.Println("receive message : ", string(message))
		// }
	}
}
