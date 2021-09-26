package v1

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	v1 "go_grpc/api/protocol/v1"

	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const apiVersion = "v1"

type ToDoServer struct {
	db *sql.DB
}

func NewToDoServer(db *sql.DB) *ToDoServer {
	return &ToDoServer{db: db}
}

func (s *ToDoServer) checkAPI(api string) error {
	if len(api) > 0 {
		if apiVersion != api {
			return status.Error(codes.Unimplemented, fmt.Sprintf("unsupported api version: service implements api version '%s',", api))
		}
	}
	return nil
}

func (s *ToDoServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Error(codes.Unknown, "链接数据库失败："+err.Error())
	}
	return c, nil
}

func (s *ToDoServer) Create(ctx context.Context, crq *v1.CreateRequest) (*v1.CreateResponse, error) {
	if err := s.checkAPI(crq.Api); err != nil {
		return nil, err
	}
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	reminder, err := ptypes.Timestamp(crq.ToDo.Reminder)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "参数错误："+err.Error())
	}
	result, err := c.ExecContext(ctx, "insert into todo (`title`,`description`,`reminder`) value (?,?,?)", crq.ToDo.Title, crq.ToDo.Description, reminder)
	if err != nil {
		return nil, status.Error(codes.Unknown, "添加失败："+err.Error())
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, status.Error(codes.Unknown, "获取添加的id失败："+err.Error())
	}
	return &v1.CreateResponse{Api: apiVersion, Id: id}, nil
}

func (s *ToDoServer) Read(ctx context.Context, rrq *v1.ReadRequest) (*v1.ReadResponse, error) {
	if err := s.checkAPI(rrq.Api); err != nil {
		return nil, err
	}
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	rows, err := c.QueryContext(ctx, "select * from todo where `id` = ?", rrq.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "查询失败："+err.Error())
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("获取数据失败：,id = '%d'查找不到", rrq.Id)+err.Error())
	}
	var td v1.ToDo
	var reminder time.Time
	if err := rows.Scan(&td.Id, &td.Title, &td.Description, &reminder); err != nil {
		return nil, status.Error(codes.Unknown, "查询数据失败："+err.Error())
	}
	if rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("查询到多条数据：id = %d", rrq.Id))
	}
	return &v1.ReadResponse{Api: apiVersion, ToDo: &td}, nil
}

func (s *ToDoServer) Update(ctx context.Context, urq *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	if err := s.checkAPI(urq.Api); err != nil {
		return nil, err
	}
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	t := time.Now().In(time.UTC)
	reminder, _ := ptypes.TimestampProto(t)
	r, _ := ptypes.Timestamp(reminder)

	result, err := c.ExecContext(ctx, "update todo set `title` = ?, `description` = ?, `reminder` = ? where `id` = ?", urq.ToDo.Title, urq.ToDo.Description, r, urq.ToDo.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "更新失败："+err.Error())
	}
	i, err := result.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "行数更新失败："+err.Error())
	}
	if i == 0 {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("更新行数的id = %d 错误", urq.ToDo.Id))
	}
	return &v1.UpdateResponse{Api: apiVersion, Updated: i}, nil
}

func (s *ToDoServer) Delete(ctx context.Context, dreq *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	if err := s.checkAPI(dreq.Api); err != nil {
		return nil, err
	}
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	result, err := c.ExecContext(ctx, "delete from todo where `id` = ?", dreq.Id)
	if err != nil {
		return nil, status.Error(codes.Unknown, "删除失败："+err.Error())
	}
	i, err := result.RowsAffected()
	if err != nil {
		return nil, status.Error(codes.Unknown, "删除更新行数失败："+err.Error())
	}
	if i == 0 {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("删除的id = %d不正确", dreq.Id))
	}
	return &v1.DeleteResponse{Api: apiVersion, Deleted: dreq.Id}, nil
}

func (s *ToDoServer) ReadAll(ctx context.Context, rreq *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	if err := s.checkAPI(rreq.Api); err != nil {
		return nil, err
	}
	c, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	result, err := c.QueryContext(ctx, "select * from todo")
	if err != nil {
		return nil, status.Error(codes.Unknown, "查询所有数据失败："+err.Error())
	}
	defer result.Close()

	list := []*v1.ToDo{}
	for result.Next() {
		td := &v1.ToDo{}
		reminder, _ := ptypes.Timestamp(td.Reminder)
		if err := result.Scan(&td.Id, &td.Title, &td.Description, &reminder); err != nil {
			return nil, status.Error(codes.Unknown, "查询失败："+err.Error())
		}
		list = append(list, td)
	}
	return &v1.ReadAllResponse{Api: apiVersion, ToDos: list}, nil
}
