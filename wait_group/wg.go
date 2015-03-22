package main

import "sync"
import "fmt"

import "errors"
import "time"

//import "github.com/GoogleCloudPlatform/kubernetes/pkg/util"

type Object struct {
	//data
}

func (obj *Object) Update(wg *sync.WaitGroup, i int) error {
	//update data
	fmt.Println(i)
	if i == 4 {
		return fmt.Errorf("This is an error")
	}
	time.Sleep(time.Second)
	fmt.Println("Update done")
	wg.Done()
	return nil
}

func main() {
	//e := make(chan error)
	//go func() {
	//	e <- errors.New("T")
	//}()
	//fmt.Printf("error %+v", <-e)
	var wg sync.WaitGroup
	e := make(chan error)
	i := make(chan bool)
	wg.Add(1)
	go func() {
		defer wg.Done()
		e <- errors.New("This is an error")
		//		wg.Done()
		i <- false
	}()
	//	wg.Wait()
	err := <-e
	fmt.Printf("Type of e %T value of e %+v", err, err)
	fmt.Println(<-i)
	wg.Wait()
	/*
		var wg sync.WaitGroup
		rep := make(chan int, 1)
		defer close(rep)
		rep <- 1
		err := make(chan error, 100)
		defer close(err)
		wg.Add(1)
		go func(r chan int, e chan error) {
			i := 1
			fmt.Printf("sleeping")
			time.Sleep(1)
			if i == 1 {
				e <- fmt.Errorf("I")
			}
			r <- i
			defer wg.Done()
			return
		}(rep, err)
		wg.Wait()
		fmt.Printf("Group done, rep %+v error %+v", <-rep, <-err)
			list := make([]Object, 5)
			num := 0
			for i := 0; i < 5; i++ {
				fmt.Printf("value of i %+v", i)
				for _, object := range list {
					wg.Add(1)
					go func(object Object) {
						err := object.Update(&wg, i)
						if err {
							num := 9
						}
						num = 10
					}(object)
				}
				//now everything has been updated. start again
				wg.Wait()
				fmt.Println("Group done")
			}
	*/
}
