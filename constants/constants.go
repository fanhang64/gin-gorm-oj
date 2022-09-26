package constants

var (
	DefaultPage = "1"
	DefaultSize = "20"
)

var SecretKey = []byte("123")

func GetStatusInfo(status int) string {
	return map[int]string{
		-1: "待判断",
		1:  "答案正确",
		2:  "答案错误",
		3:  "运行时超时",
		4:  "运行时超内存",
		5:  "编译错误",
	}[status]
}
