package main

// import (
// 	"fmt"
// 	"log"

// 	"strings"
// )

// type Point struct {
// 	x int
// 	y int
// }

// func main() {

// 	i := 45
// 	p := &i
// 	fmt.Println(*p)

// 	point := Point{5, 6}

// 	fmt.Println(point)

// 	log.SetPrefix("greetings: ")
// 	log.SetFlags(0)

// 	// name, error := greeting.Hello("pippo")

// 	// if error == nil {
// 	// 	fmt.Println(name)
// 	// 	fmt.Println(quote.Go())
// 	// } else {
// 	// 	log.Fatal(error)
// 	// }

// 	// names := []string{"marco", "paola"}

// 	// messages, _ := greeting.Hellos(names)

// 	// fmt.Println(messages)

// 	fmt.Println(WordCount("I am learning Go!"))

// }

// func WordCount(s string) map[string]int {

// 	var split []string = strings.Split(s, " ")
// 	var wordsCount = make(map[string]int)

// 	for _, v := range split {
// 		wordsCount[v] += 1
// 	}

// 	return wordsCount
// }

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
	"unicode"
)

type rot13Reader struct {
	r io.Reader
}

func (t rot13Reader) Read(buf []byte) (int, error) {

	n, err := t.r.Read(buf)


	for i := 0; i < len(buf); i++ {
		char := rune(buf[i])
		if unicode.IsLetter(char) {
			if unicode.IsUpper(char) {
				buf[i] = ((buf[i] - 65 + 13) % 26) + 65

			} else {
				buf[i] = ((buf[i] - 97 + 13) % 26) + 97
			}
		}
	}

	return n, err

}

func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Println(s)
	}
}

func sum(n []int,c chan int)  {
	sum := 0
	for i := 0; i < len(n); i++ {
		sum += n[i]
	}
	 c <- sum
}

func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func main() {

	cc := make(chan int ,10)
	quitC := make(chan int,10)

	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<- cc)
		}
		quitC <- 0
	}()

	fibonacci(cc, quitC)

	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	_, err := io.Copy(os.Stdout, &r)
	if err != nil {
		return
	}

	go say("Eccomi")
	say("ciao")

	c := make(chan int)

	numbers := []int{1,2,3,4,5,6,7,8,9,10}

	go sum(numbers[:len(numbers)/2],c)
	go sum(numbers[len(numbers)/2:],c)

	x, y := <- c,<- c

	fmt.Println(x,y,x + y)

	tick := time.Tick(100 * time.Millisecond)
	boom := time.After(500 * time.Millisecond)
	for {
		select {
		case <-tick:
			fmt.Println("tick.")
		case <-boom:
			fmt.Println("BOOM!")
			return
		default:
			fmt.Println("    .")
			time.Sleep(50 * time.Millisecond)
		}
	}

}
