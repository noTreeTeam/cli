package list

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-errors/errors"
	"github.com/spf13/afero"
	"github.com/noTreeTeam/cli/internal/migration/list"
	"github.com/noTreeTeam/cli/internal/utils"
	"github.com/noTreeTeam/cli/internal/utils/flags"
	"github.com/noTreeTeam/cli/pkg/api"
)

func Run(ctx context.Context, fsys afero.Fs) error {
	opts := api.V1ListAllSnippetsParams{ProjectRef: &flags.ProjectRef}
	resp, err := utils.GetSupabase().V1ListAllSnippetsWithResponse(ctx, &opts)
	if err != nil {
		return errors.Errorf("failed to list snippets: %w", err)
	}

	if resp.JSON200 == nil {
		return errors.New("Unexpected error listing SQL snippets: " + string(resp.Body))
	}

	table := `|ID|NAME|VISIBILITY|OWNER|CREATED AT (UTC)|UPDATED AT (UTC)|
|-|-|-|-|-|-|
`
	for _, snippet := range resp.JSON200.Data {
		table += fmt.Sprintf(
			"|`%s`|`%s`|`%s`|`%s`|`%s`|`%s`|\n",
			snippet.Id,
			strings.ReplaceAll(snippet.Name, "|", "\\|"),
			strings.ReplaceAll(string(snippet.Visibility), "|", "\\|"),
			strings.ReplaceAll(snippet.Owner.Username, "|", "\\|"),
			utils.FormatTimestamp(snippet.InsertedAt),
			utils.FormatTimestamp(snippet.UpdatedAt),
		)
	}

	return list.RenderTable(table)
}
