package server

import (
	"context"
	"errors"
	"os"

	mtd "github.com/husseinelguindi/mtd/lib"
	pb "github.com/husseinelguindi/mtd/protos/mtd"
	"github.com/husseinelguindi/mtd/server/scheduler"
)

// TODO: WRITE SCHEDULER FOR TASKS

type Server struct {
	pb.UnimplementedMtdServer

	scheduler  scheduler.Scheduler
	taskWriter *mtd.Writer
}

func NewServer() Server {
	return Server{
		scheduler:  scheduler.NewScheduler(),
		taskWriter: &mtd.Writer{},
	}
}

func (s *Server) RequestHTTPDownload(ctx context.Context, req *pb.HTTPDownloadRequest) (*pb.HTTPDownloadResponse, error) {
	f, err := os.OpenFile(req.GetFilePath(), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	newTask := mtd.NewTask(req.GetURL(), req.GetChunks(), req.GetBufSize(), f, s.taskWriter, nil)
	// newTask = mtd.Task{
	// 	URL:     req.GetURL(),
	// 	Chunks:  req.GetChunks(),
	// 	BufSize: req.GetBufSize(),

	// 	Dst:    f,
	// 	Writer: s.taskWriter,
	// }

	if len(req.GetHeaders()) > 0 {
		newTask.Headers = make(map[string]string)
		for _, header := range req.GetHeaders() {
			newTask.Headers[header.Key] = header.Val
		}
	}

	// Start the task
	// go func(t *TaskCancel) {
	// 	err := newTask.Download()
	// 	if err != nil {
	// 		t.status = pb.DownloadInfoResponse_ERRORED
	// 	} else {
	// 		t.status = pb.DownloadInfoResponse_COMPLETED
	// 	}
	// }(&taskCancel)

	id := s.scheduler.RunTask(newTask)

	return &pb.HTTPDownloadResponse{ID: id}, nil
}
func (s *Server) RequestDownloadInfo(ctx context.Context, infoReq *pb.DownloadInfoRequest) (*pb.DownloadInfoResponse, error) {
	status, ok := s.scheduler.TaskStatus(infoReq.GetID())
	if !ok {
		return nil, errors.New("task with id not found")
	}

	response := &pb.DownloadInfoResponse{}
	switch status {
	case mtd.IDLE:
		response.Status = pb.DownloadInfoResponse_QUEUED
	case mtd.IN_PROGRESS:
		response.Status = pb.DownloadInfoResponse_INPROGRESS
	case mtd.COMPLETED:
		response.Status = pb.DownloadInfoResponse_COMPLETED
	case mtd.ERRORED:
		response.Status = pb.DownloadInfoResponse_ERRORED
	}

	return response, nil
}
