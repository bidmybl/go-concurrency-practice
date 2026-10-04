package postman

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func postman(
	ctx context.Context,
	wg *sync.WaitGroup,
	transferPoint chan<- string,
	id int,
	mail string,
) {
	defer wg.Done()

	for {
		fmt.Println("I'm postman", id, "I took the letter!")
		select {
		case <-ctx.Done():
			fmt.Println("I'm postman", id, "my workday is over")
			return
		case <-time.After(1 * time.Second):
			fmt.Println("I'm postman", id, "I delivered the letter to the post office: ", mail)
		}

		select {
			case <-ctx.Done():
			fmt.Println("I'm postman", id, "my workday is over")
			return

			case transferPoint <- mail:
				fmt.Println("I'm postman", id, "I handed the letter to the post office: ", mail)
		}
	}
}

func PostmanPool(ctx context.Context, postmanCount int) <-chan string {
	mailTransferPoint := make(chan string)

	wg := &sync.WaitGroup{}

	for i := 1; i <= postmanCount; i++ {
		wg.Add(1)
		go Postman(ctx, wg, mailTransferPoint, i, postmanToMail(i))
	}

	go func() {
		wg.Wait()
		close(mailTransferPoint)
	}()

	return mailTransferPoint
}

func postmanToMail(postmanID int) string {
	ptm := map[int]string{
		1: "auto repair shop letter",
		2: "lottery letter",
		3: "letter from the family",
		4: "tax letter",
		5: "newspaper",
	}

	mail, ok := ptm[postmanID]
	if !ok {
		return "Spam"
	}

	return mail
}
