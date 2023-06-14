package internal

import (
	"log"

	"github.com/hjson/hjson-go/v4"
)

func Unmarshal(b []byte) Setting {
	var setting Setting
	if err := hjson.Unmarshal(b, &setting); err != nil {
		log.Printf("%+v", err)
	}
	return setting
}

func ReadSettingHjson(path *string) []Setting {

	if hjsonStr, err := GetText(path); err != nil {
		log.Printf("Error: %v", err)
		log.Printf("設定ファイルが不正です。: %v", *path)
		return []Setting{}
	} else {
		b := []byte(*hjsonStr)

		//*hjsonStr = fmt.Sprintf("{%v}", *hjsonStr)
		setting := Unmarshal(b)
		//log.Printf("%v", setting)

		var data map[string]interface{}
		if err := hjson.Unmarshal(b, &data); err != nil {
			log.Printf("%v", err)
			//} else {
			/*
				for k, v := range data {
					log.Printf("%v:%v", k, v)
				}
			*/
		}

		//j, _ := json.Marshal([]Setting{setting})
		//log.Printf("setting:\n%v", JsonFormat(j))

		return []Setting{setting}
	}
}
