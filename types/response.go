package types

import "encoding/xml"

type BaseResponse struct{
	XMLName xml.Name `xml:"RESPONSE"`
	RsCode string `xml:"RSCODE"`
	Msg string `xml:"MSG"`
}

type QueryCountResponse struct {
	BaseResponse
	Count int `xml:"COUNT"`
}

type QueryMoneyCountResponse struct {
	BaseResponse
	Count StringMap `xml:"COUNT"`
}



type QueryMoreResponse struct {
	BaseResponse
	Qualifiers []StringMap `xml:"QUALIFIERS>LISTS>LIST"`
}


type QueryStsCountResponse struct {
	QueryMoreResponse
}

type QueryPageDataResponse struct {
	QueryMoreResponse
}



type QueryOneResponse struct {
	BaseResponse
	Qualifiers  StringMap `xml:"QUALIFIERS"`
}



type QuerySubAccountStsResponse struct {
	QueryMoreResponse
	Total StringMap `xml:"QUALIFIERS>TOTAL"`
}



