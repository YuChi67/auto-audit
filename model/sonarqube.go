package model

// Sonarqube 資料表欄位
type Sonarqube struct {
	Id  		int64 	`gorm:"Column:id;primary_key;AUTO_INCREMENT;comment:id"`
	Department 	string 	`gorm:"column:department;type:varchar(255);not null;comment:部門"`
	Url 		string 	`gorm:"column:url;type:varchar(255);not null;comment:網址"`
	Account 	string 	`gorm:"column:account;type:varchar(255);not null;;comment:帳號"`
	Pass 		string 	`gorm:"column:pass;type:varchar(255);not null;comment:密碼"`
	BugReq 		int     `gorm:"column:bugRequest"`
	VReq 		int     `gorm:"column:vulnerableRequest"`
}

func (e *Sonarqube) TableName() string {
	return "sonarqube"
}
