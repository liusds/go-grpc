package main

import (
	"context"
	"fmt"
	v1 "go_grpc/api/protocol/v1"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
)

const apiVersion = "v1"

func main() {
	conn, err := grpc.Dial(":9000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	c := v1.NewToDoServerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	t := time.Now().In(time.UTC)
	reminder, _ := ptypes.TimestampProto(t)
	pfx := t.Format(time.RFC3339Nano)

	req1 := v1.CreateRequest{Api: apiVersion, ToDo: &v1.ToDo{Title: "title:" + pfx, Description: "description:" + pfx, Reminder: reminder}}
	res1, err := c.Create(ctx, &req1)
	if err != nil {
		log.Fatalln("创建失败：" + err.Error())
	}
	fmt.Println(res1)
	id := res1.Id
	req2 := v1.ReadRequest{Api: apiVersion, Id: id}
	res2, err := c.Read(ctx, &req2)
	if err != nil {
		log.Fatalln("读取失败：" + err.Error())
	}
	fmt.Println("读取结果：", res2)

	req3 := v1.UpdateRequest{Api: apiVersion, ToDo: &v1.ToDo{Id: res2.ToDo.Id, Title: res2.ToDo.Title, Description: res2.ToDo.Description + "updated", Reminder: res2.ToDo.Reminder}}
	res3, err := c.Update(ctx, &req3)
	if err != nil {
		log.Fatalln("更新失败：" + err.Error())
	}
	fmt.Println("更新行数：", res3)

	req4 := v1.ReadAllRequest{Api: apiVersion}
	res4, err := c.ReadAll(ctx, &req4)
	if err != nil {
		log.Fatalln("查询所有失败：" + err.Error())
	}
	fmt.Println("查询结果：", res4)

	req5 := v1.DeleteRequest{Api: apiVersion, Id: id}
	res5, err := c.Delete(ctx, &req5)
	if err != nil {
		log.Fatalln("删除失败:" + err.Error())
	}
	fmt.Println("删除数据行数：", res5)
}
