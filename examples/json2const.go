package main

import (
	"encoding/json"
	"fmt"
	"github.com/gohouse/file"
	"log"
)

type Country struct {
	CallingCode    string `json:"calling_code"`
	CallingCodeAll string `json:"calling_code_all"`
	Code           string `json:"code"`
	Emoji          string `json:"emoji"`
	IsoCodes       string `json:"iso_codes"`
	NameEn         string `json:"name"`
	NameCn         string `json:"name_cn"`
	Unicode        string `json:"unicode"`
}

func main() {
	//json2cn()
	json2Struct()
}

func json2Struct() {
	var jsa,_ = file.FileGetContents("country.json")
	var resa = map[string]Country{}
	json.Unmarshal([]byte(jsa), &resa)

	//var res = map[string]Country{}
	f := file.NewFile("xxx.txt").OpenFile()
	defer f.Close()
	//var str string
	for key, item := range resa {
		//str += fmt.Sprintf("CC_%s = %v // %s | %s\n", key, item.Code, item.NameCn,item.NameEn)
		var cnName = item.NameCn
		if cnName == "" {
			cnName = "null"
		}
		// 处理多个
		fmt.Fprintf(f, "CC_%s = %v // %s | %s\n", key, item.CallingCode, item.NameEn, cnName)
	}
	//_,err :=  f.Write([]byte(str))
	//if err != nil {
	//	log.Println("失败:", err.Error())
	//	return
	//}
	log.Println("finish")
}
