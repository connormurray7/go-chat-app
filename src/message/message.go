package message

type Message struct {
	email    string `json:"email"`
	username string `json:"username"`
	body     string `json:"body"`
}
