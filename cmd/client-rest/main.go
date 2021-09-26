package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	address := flag.String("server", "http://localhost:9000", "")
	flag.Parse()

	t := time.Now().In(time.UTC)
	pfx := t.Format(time.RFC3339Nano)

	var body string

	res, err := http.Post(*address+"/v1/todo", "application/json", strings.NewReader(fmt.Sprintf(`
	{
		"api":"v1",
		"toDo": {
			"title":"title:%s"
			"description":"description:%s",
			"reminder":"reminder:%s",
		}
	}
	`, pfx, pfx, pfx)))
	if err != nil {
		log.Fatalln("创建失败：" + err.Error())
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		body = fmt.Sprintf("create response 读取失败:%v", err)
	} else {
		body = string(bodyBytes)
	}
	fmt.Printf("create response data:%s", body)

	var Created struct {
		API string `json:"api"`
		ID  string `json:"id"`
	}

	if err := json.Unmarshal(bodyBytes, &Created); err != nil {
		log.Fatalln("create bodyBytes data json unmarshal failed")
	}

	//Call Read
	res, err = http.Get(fmt.Sprintf(*address + "/v1/todo/" + Created.ID))
	if err != nil {
		log.Fatalln("Read failed:" + err.Error())
	}
	bodyBytes, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		body = fmt.Sprintln("read response read failed:" + err.Error())
	} else {
		body = string(bodyBytes)
	}
	fmt.Printf("read response:%s", body)

	//Call Update
	req, err := http.NewRequest(http.MethodPut, *address+"/v1/todo/"+Created.ID, strings.NewReader(fmt.Sprintf(`
	{
		"api":"v1",
		"toDo": {
			"title":"title:%s + updated",
			"description":"description:%s + updated",
			"reminder":"reminder:%s + updated",
		}
	}
	`, pfx, pfx, pfx)))
	req.Header.Set("Content-Type", "application/json")
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln("update method failed:", err)
	}
	bodyBytes, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		body = fmt.Sprintln("update response read body failed:", err)
	} else {
		body = string(bodyBytes)
	}
	fmt.Printf("update response data:%s", body)

	//Call ReadAll
	res, err = http.Get(*address + "/v1/todo/all")
	if err != nil {
		log.Fatalln("readall method failed:", err)
	}

	bodyBytes, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		body = fmt.Sprintln("readall response data read failed:", err)
	} else {
		body = string(bodyBytes)
	}
	fmt.Printf("readall response data:%s", body)

	//Call Delete
	req, err = http.NewRequest(http.MethodDelete, *address+"/v1/todo/"+Created.ID, nil)
	if err != nil {
		log.Fatalln("delete method failed:", err)
	}
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln("delete defaultclient failed:", err)
	}

	bodyBytes, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		body = fmt.Sprintln("delete readall failed:", err)
	} else {
		body = string(bodyBytes)
	}
	fmt.Printf("delete response:%s", body)
}
