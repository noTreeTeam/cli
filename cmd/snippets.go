package cmd

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/noTreeTeam/cli/internal/snippets/download"
	"github.com/noTreeTeam/cli/internal/snippets/list"
	"github.com/noTreeTeam/cli/internal/utils/flags"
)

var (
	snippetsCmd = &cobra.Command{
		GroupID: groupManagementAPI,
		Use:     "snippets",
		Short:   "Manage Supabase SQL snippets",
	}

	snippetsListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all SQL snippets",
		Long:  "List all SQL snippets of the linked project.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return list.Run(cmd.Context(), afero.NewOsFs())
		},
	}

	snippetsDownloadCmd = &cobra.Command{
		Use:   "download <snippet-id>",
		Short: "Download contents of a SQL snippet",
		Long:  "Download contents of the specified SQL snippet.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return download.Run(cmd.Context(), args[0], afero.NewOsFs())
		},
	}
)

func init() {
	snippetsCmd.PersistentFlags().StringVar(&flags.ProjectRef, "project-ref", "", "Project ref of the Supabase project.")
	snippetsCmd.AddCommand(snippetsListCmd)
	snippetsCmd.AddCommand(snippetsDownloadCmd)
	rootCmd.AddCommand(snippetsCmd)
}
