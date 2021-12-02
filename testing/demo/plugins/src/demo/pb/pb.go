package pb

import "fmt"

func OnReadySending(m, c, r map[string][]byte) {
	fmt.Println("OnReadySending message templates", m, "config", c, "request attributes", r)
}

func OnReceived(m, c, r map[string][]byte) {
	fmt.Println("OnReceived message templates", m, "config", c, "response", r)
}

func OnError(m, c, r map[string][]byte) {
	fmt.Println("OnError message templates", m, "config", c, "response", r)
}
