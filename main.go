package main

import (
	"auto-audit/common"
	"auto-audit/console"
	"auto-audit/module"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func main() {
	info, err := console.Read()
	if err != nil {
		return
	}

	if strings.HasSuffix(info.Url, "/") {
		info.Url = info.Url[:len(info.Url)-1]	}

	var notifyStr []string
	for _, n := range info.Notify {
		notifyStr = append(notifyStr, common.NotifyStr[n])
	}

	fmt.Println("取得變數資料")
	fmt.Println("sonarqube 網址: " + info.Url)
	fmt.Println("sonarqube project key: " + info.ProjectKey)
	fmt.Println("sonarqube token: " + info.Token)
	fmt.Println("sonarqube 帳號: " + info.Account)
	fmt.Println("是否檢查發版要求: " + strconv.FormatBool(info.CheckRequire))
	fmt.Println("寄送郵件檢查內容: " + strings.Join(notifyStr,","))
	fmt.Println("bug 發版要求: " + strconv.Itoa(info.BugReq))
	fmt.Println("代碼漏洞發版要求: " + strconv.Itoa(info.VReq))

	scan, err := module.GetCodeScan(info)
	if err != nil {
		return
	}

	cover, err := module.SearchHistory(info)
	if err != nil {
		return
	}

	bugScan, err := strconv.Atoi(scan.Bug)
	if err != nil {
		return
	}
	vScan, err := strconv.Atoi(scan.Vulnerability)
	if err != nil {
		return
	}

	coverScan, err := strconv.ParseFloat(scan.Coverage, 64)
	if err != nil {
		log.Println(err)
		return
	}

	var isSend = false
	if len(info.Notify) == 0 {
		isSend = true
	} else {
		for _, notify := range info.Notify {
			switch notify {
			case common.BUG:
				if bugScan > info.BugReq {
					isSend = true
					break
				}
			case common.VULNERABILITY:
				if vScan > info.VReq {
					isSend = true
					break
				}
			case common.COVERAGE:
				if coverScan < cover {
					isSend = true
					break
				}
			}
		}

		if isSend {
			if err = module.SendEmail(info, scan); err != nil {
				return
			}
			fmt.Println("郵件已寄送")
		}
	}

	if info.CheckRequire {
		if bugScan > info.BugReq {
			err = fmt.Errorf("sonarqube 掃描結果 bug數 %d 不符合發版要求 %d", bugScan, info.BugReq)
			panic(err.Error())
		}
		if vScan > info.VReq {
			err = fmt.Errorf("sonarqube 掃描結果 漏洞數 %d 不符合發版要求 %d", vScan, info.VReq)
			panic(err.Error())
		}
		if cover > coverScan {
			err = fmt.Errorf("sonarqube 掃描結果 覆蓋率 %g 低於前次 %g", cover, coverScan)
			panic(err.Error())
		}
	}
	fmt.Println("done")
	
}
