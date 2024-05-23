package tool

import (
	"github.com/duke-git/lancet/v2/cryptor"
	"sort"
	"strings"
)

/*
 * @Description URL请求参数加签名。
 * @param urlParam   : URL请求参数
 * @param privateKey : 签名私钥
 * @return 返回参数：增加了sign参数的URL
 * @exception
 * @see
 */

func AddUrlSign(urlParam string, privateKey string) string {
	queryStringChange := ""
	if len(urlParam) == 0 {
		fieldList := strings.Split(urlParam, "&")
		if fieldList != nil && len(fieldList) > 0 {
			// 去除等号
			for i, _ := range fieldList {
				fieldList[i] = strings.Replace(fieldList[i], "=", "", 1)
			}
			sort.Strings(fieldList)
			queryStringChange = privateKey
			for i, _ := range fieldList {
				queryStringChange = queryStringChange + fieldList[i]
			}
		}
		hmacSign := cryptor.HmacSha1(queryStringChange, privateKey)
		return urlParam + "&sign=" + hmacSign
	} else {
		hmacSign := cryptor.HmacSha1(queryStringChange, privateKey)
		return urlParam + "sign=" + hmacSign
	}

}
