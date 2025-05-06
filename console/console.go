package console

import (
	"auto-audit/common"
	"auto-audit/model"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Read 讀取環境變數
func Read() (*model.Info, error) {
	if _, err := time.LoadLocation("Local"); err != nil {
		return nil, err
	}
	fmt.Println()
	fmt.Println("[" + time.Now().String() + "]取得環境變數")

	projectKey   := flag.String("projectKey", os.Getenv("PROJECTKEY"), "請輸入 sonarqube project key")
	to		     := flag.String("to", os.Getenv("TO"), "請輸入收件人(請使用 , 分隔)")
	url 		 := flag.String("url", os.Getenv("URL"), "請輸入 sonarqube 網址")
	token 	 	 := flag.String("token", os.Getenv("TOKEN"), "請輸入 sonarqube token")
	account  	 := flag.String("account", os.Getenv("ACCOUNT"), "請輸入 sonarqube 帳號")
	password 	 := flag.String("password", os.Getenv("PASSWORD"), "請輸入 sonarqube 密碼")
	checkRequire := flag.String("checkRequire", os.Getenv("CHECKREQUIRE"), "請輸入是否檢查發版要求(true,false)")
	notify 		 := flag.String("notify", os.Getenv("NOTIFY"), "請輸入寄送郵件檢查內容 ex. bug,vulnerability")
	bugReq 		 := flag.String("bugReq", os.Getenv("BUG_REQUEST"), "請輸入 bug 發版要求, 不填則預設為 0")
	vReq 		 := flag.String("vReq", os.Getenv("VULNERABILITY_REQUEST"), "請輸入代碼漏洞發版要求, 不填則預設為 0")

	var (
		crBool    = false
		notifyNew []common.NotifyEnum
	)

	if *checkRequire == "" {
		*checkRequire = "0"
	}
	if *bugReq == "" {
		*bugReq = "0"
	}
	if *vReq == "" {
		*vReq = "0"
	}

	bugReqInt, _ := strconv.Atoi(*bugReq)
	vReqInt, _   := strconv.Atoi(*vReq)
	toArr		 := strings.Split(*to, ",")
	notifyArr 	 := strings.Split(*notify, ",")

	for _, n := range notifyArr {
		notifyNew = append(notifyNew, common.Notify[n])
	}

	if *url == "" {
		return nil, errors.New("url not found")
	}

	if *checkRequire != "0" && *checkRequire != "false" {
		crBool = true
	}

	flag.Parse()
	return &model.Info{
		ProjectKey:   *projectKey,
		To:           toArr,
		Url:          *url,
		Token:        *token,
		Account:      *account,
		Password:     *password,
		CheckRequire: crBool,
		Notify:       notifyNew,
		BugReq:       bugReqInt,
		VReq:         vReqInt,
	}, nil
}
