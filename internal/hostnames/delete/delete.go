package delete

import (
	"context"
	"fmt"

	"github.com/go-errors/errors"
	"github.com/spf13/afero"
	"github.com/noTreeTeam/cli/internal/utils"
)

func Run(ctx context.Context, projectRef string, fsys afero.Fs) error {
	// 1. Sanity checks.
	// 2. delete config
	{
		resp, err := utils.GetSupabase().V1DeleteHostnameConfigWithResponse(ctx, projectRef)
		if err != nil {
			return errors.Errorf("failed to delete custom hostname: %w", err)
		}
		if resp.StatusCode() != 200 {
			return errors.New("failed to delete custom hostname config; received: " + resp.Status())
		}
		fmt.Println("Deleted custom hostname config successfully.")
		return nil
	}
}
