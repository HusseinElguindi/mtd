package main

import (
	"github.com/husseinelguindi/mtd/client/mtd/cmd"
)

func main() {
	cmd.Execute()

	// conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	// if err != nil {
	// 	log.Fatalln("could not open a tcp connection")
	// }
	// defer conn.Close()

	// client := pb.NewMtdClient(conn)
	// dr := &pb.HTTPDownloadRequest{
	// 	URL:      "https://file-examples-com.github.io/uploads/2017/04/file_example_MP4_1920_18MG.mp4",
	// 	Chunks:   uint32(runtime.NumCPU()),
	// 	BufSize:  1024 * 1024,
	// 	FilePath: "./vid.mp4",
	// }
	// resp, err := client.RequestHTTPDownload(context.Background(), dr)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// log.Println("id:", resp.GetID())
	// resp2, err := client.RequestDownloadInfo(context.Background(), &pb.DownloadInfoRequest{ID: resp.GetID()})
	// log.Printf("%#v %v", resp2.GetStatus(), err)
}
