package cmd

import (
	"context"
	"fmt"
	"runtime"

	pb "github.com/husseinelguindi/mtd/protos/mtd"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var getCmd = &cobra.Command{
	Use:   "get <url> <filename>",
	Short: "HTTP GETs the passed URL using default settings, unless overriding flags are set",
	Args:  cobra.ExactArgs(2),
	RunE:  get,
}

func init() {
	rootCmd.AddCommand(getCmd)
}

func get(cmd *cobra.Command, args []string) error {
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		// log.Fatalln("could not open a tcp connection")
		return err
	}
	defer conn.Close()

	client := pb.NewMtdClient(conn)
	dr := &pb.HTTPDownloadRequest{
		URL:      args[0],
		Chunks:   uint32(runtime.NumCPU()),
		BufSize:  1024 * 1024,
		FilePath: args[1],
	}
	resp, err := client.RequestHTTPDownload(context.Background(), dr)
	if err != nil {
		return err
	}

	fmt.Println("task id:", resp.GetID())
	return nil
}
