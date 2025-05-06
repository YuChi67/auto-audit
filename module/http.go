package module

import (
	"auto-audit/model"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"strconv"
	"time"
)

const MetricKey = "coverage,bugs,vulnerabilities,code_smells,duplicated_lines_density"

func GetCodeScan(info *model.Info) (*model.CodeScan, error) {
	var (
		url   = info.Url
		token = info.Token
	)

	time.Sleep(10 * time.Second)
	serviceUrl := url + "/api/measures/component?component=" + info.ProjectKey + "&metricKeys=" + MetricKey
	client := &http.Client{}
	req, err := http.NewRequest("GET", serviceUrl, nil)
	if err != nil {
		log.Println("http.NewRequest error:" + err.Error())
		return nil, err
	}
	if token == "" {
		req.Header.Add("Authorization", "Basic " +
			base64.StdEncoding.EncodeToString([]byte(info.Account + ":" + info.Password)))
	} else {
		req.SetBasicAuth(token, "")
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("client.Do error", err.Error())
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Println(err)
		}
	}(resp.Body)
	if resp.StatusCode == 401 {
		log.Println("401 Unauthorized")
		return nil, errors.New("401 Unauthorized")
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ioutil.ReadAll error", err.Error())
		return nil, err
	}
	if resp.StatusCode != 200 {
		var e *model.Error
		err = json.Unmarshal(responseData, &e)
		if err != nil {
			log.Println("ioutil.ReadAll error", err.Error())
			return nil, err
		}
		log.Println("code scan error: " + e.Errors[0].Msg)
		return nil, errors.New(e.Errors[0].Msg)
	}
	var component *model.Component
	err = json.Unmarshal(responseData, &component)
	if err != nil {
		log.Println("ioutil.ReadAll error", err.Error())
		return nil, err
	}

	var codeScan = new(model.CodeScan)
	codeScan.ProjectName = component.Component.Name
	for _, r := range component.Component.Measures {
		if r.Metric == "bugs" {
			codeScan.Bug = r.Value
		} else if r.Metric == "vulnerabilities" {
			codeScan.Vulnerability = r.Value
		} else if r.Metric == "coverage" {
			codeScan.Coverage = r.Value
		} else if r.Metric == "code_smells" {
			codeScan.CodeSmell = r.Value
		} else if r.Metric == "duplicated_lines_density" {
			codeScan.Duplicate = r.Value
		}
	}

	fmt.Println()
	fmt.Println("掃描結果: ")
	fmt.Println("專案名稱: " + codeScan.ProjectName)
	fmt.Println("bug數: " + codeScan.Bug)
	fmt.Println("漏洞數: " + codeScan.Vulnerability)
	fmt.Println("Code Smells: " + codeScan.CodeSmell)
	fmt.Println("測試覆蓋率: " + codeScan.Coverage)
	fmt.Println("代碼重複率: " + codeScan.Duplicate)

	return codeScan, nil
}

func SearchHistory(info *model.Info) (float64, error) {
	var (
		url   = info.Url
		token = info.Token
		cover float64
	)

	serviceUrl := url + "/api/measures/search_history?component=" + info.ProjectKey + "&metrics=coverage&ps=1000"
	client := &http.Client{}
	req, err := http.NewRequest("GET", serviceUrl, nil)
	if err != nil {
		log.Println("http.NewRequest error:" + err.Error())
		return 0, err
	}
	if token == "" {
		req.Header.Add("Authorization", "Basic " +
			base64.StdEncoding.EncodeToString([]byte(info.Account + ":" + info.Password)))
	} else {
		req.SetBasicAuth(token, "")
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("client.Do error", err.Error())
		return 0, err
	}

	defer func(Body io.ReadCloser) {
		if err = Body.Close(); err != nil {
			log.Println(err)
		}
	}(resp.Body)
	if resp.StatusCode == 401 {
		log.Println("401 Unauthorized")
		return 0, errors.New("401 Unauthorized")
	}
	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ioutil.ReadAll error", err.Error())
		return 0, err
	}
	if resp.StatusCode != 200 {
		var e *model.Error
		err = json.Unmarshal(responseData, &e)
		if err != nil {
			log.Println("ioutil.ReadAll error", err.Error())
			return 0, err
		}
		log.Println("code scan error: " + e.Errors[0].Msg)
		return 0, errors.New(e.Errors[0].Msg)
	}
	var history *model.SearchHistory
	err = json.Unmarshal(responseData, &history)
	if err != nil {
		log.Println("ioutil.ReadAll error", err.Error())
		return 0, err
	}

	for _, r := range history.Measures {
		if r.Metric == "coverage" {
			fmt.Println("前次推送的代碼覆蓋率: " + r.History[len(r.History)-2].Value)
			s, err := strconv.ParseFloat(r.History[len(r.History)-2].Value, 64)
			if err != nil {
				return 0, err
			}
			cover = s
		}
	}

	return cover, nil
}

func SendEmail(info *model.Info, scan *model.CodeScan) error {
	// Gmail 資訊
	from := "yuchi060703@gmail.com"
	password := "xxxxxx"

	// Gmail SMTP 設定
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := "稽核掃描結果\\r\\n" + "\\r\\n" + "div style=\"max-width: 1024px;padding: 0 20px;overflow: auto;\">" +
		"<table style=\"border-collapse: collapse; color: #4a4a4d; width: 100%;\">" +
		"<thead style=\"background: #395870; color: #fff\"><tr>" +
		"<th scope=\"col\" style=\"padding: 6px 10px;vertical-align: middle;\">專案名稱</th>" +
		"<th scope=\"col\" style=\"padding: 6px 10px;vertical-align: middle;\">bug數</th>" +
		"<th scope=\"col\" style=\"padding: 6px 10px;vertical-align: middle;\">漏洞數</th>" +
		"<th scope=\"col\" style=\"padding: 6px 10px;vertical-align: middle;\">Code Smells</th>" +
		"<th scope=\"col\" style=\"padding: 6px 10px;vertical-align: middle;\">測試覆蓋率</th>" +
		"<th scope=\"col\" style=\"padding: 6px 10px;vertical-align: middle;\">代碼重複率</th></tr></thead><tbody><tr>" +
		"<td style=\"padding: 6px 10px;vertical-align: middle;border-bottom: 1px solid #cecfd5;border-right: 1px solid #cecfd5;\">" + scan.ProjectName + "</td>" +
		"<td style=\"padding: 6px 10px;vertical-align: middle;border-bottom: 1px solid #cecfd5;border-right: 1px solid #cecfd5;\">" + scan.Bug + "</td>" +
		"<td style=\"padding: 6px 10px;vertical-align: middle;border-bottom: 1px solid #cecfd5;border-right: 1px solid #cecfd5;\">" + scan.Vulnerability + "</td>" +
		"<td style=\"padding: 6px 10px;vertical-align: middle;border-bottom: 1px solid #cecfd5;border-right: 1px solid #cecfd5;\">" + scan.CodeSmell + "</td>" +
		"<td style=\"padding: 6px 10px;vertical-align: middle;border-bottom: 1px solid #cecfd5;border-right: 1px solid #cecfd5;\">" + scan.Coverage + "</td>" +
		"<td style=\"padding: 6px 10px;vertical-align: middle;border-bottom: 1px solid #cecfd5;border-right: 1px solid #cecfd5;\">" + scan.Duplicate + "</td>" +
		"</tr></tbody></table></div>"

	// 認證資訊
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// 發送郵件
	fmt.Println("郵件寄送")
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, info.To, []byte(message))
	if err != nil {
		log.Fatal("發送失敗：", err)
	}

	return nil
}
