package models

// AppToken : 공공데이터 포털로부터 발행 된 Token 및 관련 App 정보 처리
type AppToken struct {
	Token     string `json:"token" form:"token" query:"token"`
	NameSpace string `json:"nameSpace" form:"nameSpace" query:"nameSpace"`
}
