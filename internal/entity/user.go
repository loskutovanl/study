package entity

// User содержит информацию о пользователе: имя, возраст, список друзей
type User struct {
	Name    string   `json:"name"`
	Age     int      `json:"age"`
	Friends []string `json:"friends"`
}

// Friends содержит информацию об id двух пользователей, отправивших запрос на дружбу
type Friends struct {
	SourceId int `json:"source_id"`
	TargetId int `json:"target_id"`
}

// DeleteUser содержит информацию об id пользователя, на которого запрашивается удаление
type DeleteUser struct {
	TargetId int `json:"target_id"`
}

// NewAge содержит инормацию о новом возрасте пользователя
type NewAge struct {
	Age string `json:"new_age"`
}
