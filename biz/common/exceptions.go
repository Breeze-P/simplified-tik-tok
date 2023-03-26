package common

// 统一错误处理
const (
	ParamaError    = 1000 // 请求格式错误
	DataBaseError  = 1001 // 数据库请求错误
	FileError      = 1002 // 文件存储出错
	PrivilegeError = 1003 // 权限不足
	OtherError     = 1004 // 其他错误
)

// 错误结构体
type APIException struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

func newAPIException(statusCode int, statusMsg string) *APIException {
	return &APIException{
		StatusCode: statusCode,
		StatusMsg:  statusMsg,
	}
}

func NewParameterError(msg string) *APIException {
	return newAPIException(ParamaError, msg)
}

func NewDatabaseError(msg string) *APIException {
	return newAPIException(DataBaseError, msg)
}

func NewFileError(msg string) *APIException {
	return newAPIException(FileError, msg)
}

func NewPrivilegeError(msg string) *APIException {
	return newAPIException(PrivilegeError, msg)
}

func NewOtherError(msg string) *APIException {
	return newAPIException(OtherError, msg)
}
