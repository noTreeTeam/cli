package restore

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-errors/errors"
	"github.com/noTreeTeam/cli/internal/utils"
	"github.com/noTreeTeam/cli/internal/utils/flags"
	"github.com/noTreeTeam/cli/pkg/api"
)

func Run(ctx context.Context, timestamp int64) error {
	body := api.V1RestorePitrBody{RecoveryTimeTargetUnix: timestamp}
	resp, err := utils.GetSupabase().V1RestorePitrBackupWithResponse(ctx, flags.ProjectRef, body)
	if err != nil {
		return errors.Errorf("failed to restore backup: %w", err)
	} else if resp.StatusCode() != http.StatusCreated {
		return errors.Errorf("unexpected restore backup status %d: %s", resp.StatusCode(), string(resp.Body))
	}
	fmt.Println("Started PITR restore:", flags.ProjectRef)
	return nil
}
