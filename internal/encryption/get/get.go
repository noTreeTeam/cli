package get

import (
	"context"
	"fmt"

	"github.com/go-errors/errors"
	"github.com/noTreeTeam/cli/internal/utils"
)

func Run(ctx context.Context, projectRef string) error {
	resp, err := utils.GetSupabase().V1GetPgsodiumConfigWithResponse(ctx, projectRef)
	if err != nil {
		return errors.Errorf("failed to retrieve pgsodium config: %w", err)
	}

	if resp.JSON200 == nil {
		return errors.New("Unexpected error retrieving project root key: " + string(resp.Body))
	}

	fmt.Println(resp.JSON200.RootKey)
	return nil
}
