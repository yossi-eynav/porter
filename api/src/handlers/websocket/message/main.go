package message
import "time"

type Message struct {
	Body string
	Color string
	Timestamp time.Time
}

func SendMessage(msg string, color string,  ch chan Message, shouldBlock bool) {
	message := Message{
		msg,
		color,
		time.Now(),
	}
	if shouldBlock {
		ch <- message
	} else {
		go func() {
			ch <- message
		}()
	}

}
