package model

import "auto-audit/common"

type Info struct {
	ProjectKey   string
	To			 []string
	Url		  	 string
	Token		 string
	Account		 string
	Password	 string
	CheckRequire bool
	Notify 		 []common.NotifyEnum
	BugReq 		 int
	VReq		 int
}

type Component struct {
	Component struct {
		Key      string `json:"key"`
		Name     string `json:"name"`
		Measures []struct {
			Metric    string `json:"metric"`
			Value     string `json:"value"`
			BestValue bool   `json:"bestValue"`
		} `json:"measures"`
	} `json:"component"`
}

type Error struct {
	Errors []struct {
		Msg string `json:"msg"`
	} `json:"errors"`
}

type CodeScan struct {
	ProjectName   string `json:"project_name"`
	Bug           string `json:"bug,omitempty"`
	Vulnerability string `json:"vulnerability,omitempty"`
	Coverage      string `json:"coverage,omitempty"`
	CodeSmell     string `json:"code_smell,omitempty"`
	Duplicate     string `json:"duplicate,omitempty"`
}

type EmailInfo struct {
	UserIds []string `json:"userIds"`
	MsgType string   `json:"msgType"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
}

type SearchHistory struct {
	Paging struct {
		PageIndex int `json:"pageIndex"`
		PageSize  int `json:"pageSize"`
		Total     int `json:"total"`
	} `json:"paging"`
	Measures []struct {
		Metric  string `json:"metric"`
		History []struct {
			Date  string `json:"date"`
			Value string `json:"value"`
		} `json:"history"`
	} `json:"measures"`
}