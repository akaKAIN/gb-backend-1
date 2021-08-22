package example

type User struct {
	Name   string  `json:"name"`
	Age    int     `json:"age"`
	Wallet float32 `json:"wallet"`
}

type Upload struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

type Uploads struct {
	Uploads []Upload `json:"uploads"`
}
