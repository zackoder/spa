package utils

type SignupRequest struct {
	NickName  string `json:"nickName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Gender    string `json:"gender"`
	Age       string `json:"age"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type Users struct {
	Nickname         string `json:"nickname"`
	Firstname        string `json:"firstname"`
	Lastname         string `json:"lastname"`
	Gender           string `json:"gender"`
	Age              string `json:"age"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	Confirm_Password string `json:"confirmPassword"`
}

type Posts struct {
	Id         int       `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Poster     string    `json:"poster"`
	CreatedAt  int       `json:"createdAt"`
	Categories []string  `json:"categories"`
	Reactions  Reactinos `json:"reactions"`
}

type Reactinos struct {
	Likes    int    `json:"likes"`
	Dislikes int    `json:"dislikes"`
	Action   string `json:"action"`
}

type Name struct {
	Name string `json:"nickname"`
}

type Category struct {
	Name string `json:"name"`
}

type Message struct {
	To        string `json:"to"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	CreatedAt int    `json:"creationDate"`
}

type Comment struct {
	PostId       int    `json:"id"`
	Comment      string `json:"comment"`
	CreationDate int    `json:"creationDate"`
	Username     string `json:"username"`
}

// type SigninRequest struct {
// 	Userinpt string `json:"email"`
// 	Password string `json:"password"`
// }

type LoginUsers struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	Err  string
	Code int
}
