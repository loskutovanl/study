package entity

// User содержит информацию о пользователе: id, имя, возраст, список друзей
type User struct {
	Id      int
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Friends []int  `json:"friends"`
}

// Friends содержит информацию об id двух пользователей, отправивших запрос на дружбу
type Friends struct {
	SourceId int `json:"source_id"`
	TargetId int `json:"target_id"`
}

// NewAge содержит инормацию о новом возрасте пользователя
type NewAge struct {
	Id  int
	Age int `json:"new_age"`
}
