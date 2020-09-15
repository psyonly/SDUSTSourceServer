package gotools

var logHash = map[string]string{
	"HEAD":  "[PSY]>>> ",
	"VALUE": "Value is ",
}

func AddLogHead(str string) string {
	return logHash["HEAD"] + str
}

func PrintLogHead() string {
	return logHash["HEAD"]
}

func PrintLogValue() string {
	return logHash["VALUE"]
}
