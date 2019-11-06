package micro

import (
	"fmt"
	"log"
	"net/rpc"
	"strconv"
	"sync"
	"testing"
)

func TestClient(t *testing.T) {
	//Client()

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		//go func(i int) {
		//	defer wg.Done()

			client, err := rpc.Dial("tcp", "localhost:1234")
			if err != nil {
				log.Fatal("dialing:", err)
			}


			var reply *string = new(string)
		helloCall := client.Go("HelloService.Hello", "hello" + strconv.Itoa(i), reply, nil)
		helloCall = <-helloCall.Done
		if err := helloCall.Error; err != nil {
			log.Fatal(err)
		}

		args := helloCall.Args.(string)
		reply = helloCall.Reply.(*string)
		fmt.Println(args, reply)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(reply)
		client.Close()
		wg.Done()
		//}(i)
	}
	wg.Wait()
}

func TestClient2(t *testing.T) {

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			client, err := rpc.Dial("tcp", "localhost:1234")
			if err != nil {
				log.Fatal("dialing:", err)
			}
			defer client.Close()

			var reply string
			err = client.Call("HelloService.Hello", "hello", &reply)
			if err != nil {
				fmt.Println(err)
				//log.Fatal(err)
			}

			fmt.Println(reply)
		}(i)
	}
	wg.Wait()
}

func TestClient3(t *testing.T) {

	client, err := rpc.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer client.Close()

	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			var reply string
			err = client.Call("HelloService.Hello", "hello", &reply)
			if err != nil {
				fmt.Println(err)
				//log.Fatal(err)
			}

			fmt.Println(reply)
		}(i)
	}
	wg.Wait()
}