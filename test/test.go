package main

import (
	"encoding/json"
	"fmt"
)

type Obj struct {
	ID    string
	Value int
}

func main() {
	t := make(map[string]Obj)

	t["key"] = Obj{
		ID:    "id",
		Value: 123,
	}
	t["key2"] = Obj{
		ID:    "id2",
		Value: 1234,
	}

	b, _ := json.Marshal(t)
	fmt.Println("Json : ", string(b))
}
