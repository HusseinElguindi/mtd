package cmd

import (
	"context"
	"fmt"

	pb "github.com/husseinelguindi/mtd/protos/mtd"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var statusCmd = &cobra.Command{
	Use:   "status <id>",
	Short: "returns the status of the task with the passed id",
	Args:  cobra.ExactArgs(1),
	RunE:  status,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func status(cmd *cobra.Command, args []string) error {
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		// log.Fatalln("could not open a tcp connection")
		return err
	}
	defer conn.Close()

	client := pb.NewMtdClient(conn)
	dr := &pb.DownloadInfoRequest{ID: args[0]}
	resp, err := client.RequestDownloadInfo(context.Background(), dr)
	if err != nil {
		return err
	}

	fmt.Println("task status: ", resp.GetStatus().String())
	return nil
}
