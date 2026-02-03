package file

import (
	"github.com/spf13/cobra"
	"github.com/anyproto/anytype-cli/cmd/file/upload"
)

func NewFileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "file <command>",
		Short: "Manage files",
		Long:  "Upload, download, and manage files in Anytype",
	}

	cmd.AddCommand(upload.NewUploadCmd())

	return cmd
}
