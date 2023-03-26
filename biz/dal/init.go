package dal

import (
	"simplified-tik-tok/biz/dal/mongodb"
	"simplified-tik-tok/biz/dal/mysql"
)

func Init() {
	mysql.Init()
	mongodb.Init()
}
