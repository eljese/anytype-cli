package upload

import (
	"context"
	"fmt"

	"github.com/anyproto/anytype-heart/pb"
	"github.com/anyproto/anytype-heart/pb/service"
	"github.com/spf13/cobra"

	"github.com/anyproto/anytype-cli/core"
	"github.com/anyproto/anytype-cli/core/output"
)

func NewUploadCmd() *cobra.Command {
	var spaceId string

	cmd := &cobra.Command{
		Use:   "upload <file_path>",
		Short: "Upload a file to Anytype",
		Long:  "Upload a local file to a specific Anytype space",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]

			if spaceId == "" {
				// Default to tech space if not provided
				// But usually we want to specify a space.
				// For now, let's just fail if spaceId is missing or ask the user.
				return fmt.Errorf("space-id is required")
			}

			err := core.GRPCCall(func(ctx context.Context, client service.ClientCommandsClient) error {
				req := &pb.RpcFileUploadRequest{
					SpaceId:   spaceId,
					LocalPath: filePath,
				}
				resp, err := client.FileUpload(ctx, req)
				if err != nil {
					return fmt.Errorf("failed to upload file: %w", err)
				}
				if resp.Error != nil && resp.Error.Code != pb.RpcFileUploadResponseError_NULL {
					return fmt.Errorf("file upload error: %s", resp.Error.Description)
				}

				output.Success("File uploaded successfully!")
				output.Info("Object ID: %s", resp.ObjectId)
				return nil
			})

			return err
		},
	}

	cmd.Flags().StringVar(&spaceId, "space-id", "", "ID of the space to upload to")

	return cmd
}
