package globalkey

const (
	RequestIdKey                = "request_id"
	WechatAccessTokenUrl        = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=wx9770b8510291e3bf&secret=19a7ab64b942667f9149766128f12876"
	WechatGetUserPhoneNumberUrl = "https://api.weixin.qq.com/wxa/business/getuserphonenumber?access_token="
)

const (
	TOPIC_PAT      = 0
	TOPIC_Gift     = 1
	TOPIC_WITHDRAW = 2
)

const (
	// 2006-01-02 15:04:05为go的标准格式
	timeFormat1 = "2006-01-02 15:04:05"
	timeFormat2 = "2006-01-02-15-04-05"
	timeFormat3 = "2006-01-02 15:04:05.000" // 可以精确到毫秒ms
	// 非go的标准格式 打印的时间格式为未知时间
	timeFormat4 = "2022-01-02 10:04:05"
	timeFormat5 = "2022-01-02-10-04-05"
	timeFormat6 = "2022-01-02 10:04:05.000" // 可以精确到毫秒ms
)
