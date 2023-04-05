package main

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func test1() error {
	return fmt.Errorf("err1")
}

func test2(err error) error {
	return errors.Wrap(err, "err2")
}

func test3(err error) {
	logrus.Panicf("%+v", err)
}

func main() {
	// err := test2(test1())
	// logrus.Errorf("%+v", err)
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case now := <-ticker.C:
			fmt.Println(now)
		}
	}
}
