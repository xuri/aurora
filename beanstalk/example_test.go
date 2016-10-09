package beanstalk_test

import (
	"fmt"
	"github.com/kr/beanstalk"
	"time"
)

var conn, _ = beanstalk.Dial("tcp", "127.0.0.1:11300")

func Example_reserve() {
	id, body, err := conn.Reserve(5 * time.Second)
	if err != nil {
		panic(err)
	}
	fmt.Println("job", id)
	fmt.Println(string(body))
}

func Example_reserveOtherTubeSet() {
	tubeSet := beanstalk.NewTubeSet(conn, "mytube1", "mytube2")
	id, body, err := tubeSet.Reserve(10 * time.Hour)
	if err != nil {
		panic(err)
	}
	fmt.Println("job", id)
	fmt.Println(string(body))
}

func Example_put() {
	id, err := conn.Put([]byte("myjob"), 1, 0, time.Minute)
	if err != nil {
		panic(err)
	}
	fmt.Println("job", id)
}

func Example_putOtherTube() {
	tube := &beanstalk.Tube{Conn: conn, Name: "mytube"}
	id, err := tube.Put([]byte("myjob"), 1, 0, time.Minute)
	if err != nil {
		panic(err)
	}
	fmt.Println("job", id)
}
