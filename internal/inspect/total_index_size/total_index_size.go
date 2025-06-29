package total_index_size

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/go-errors/errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/spf13/afero"
	"github.com/noTreeTeam/cli/internal/db/reset"
	"github.com/noTreeTeam/cli/internal/migration/list"
	"github.com/noTreeTeam/cli/internal/utils"
	"github.com/noTreeTeam/cli/pkg/pgxv5"
)

//go:embed total_index_size.sql
var TotalIndexSizeQuery string

type Result struct {
	Size string
}

func Run(ctx context.Context, config pgconn.Config, fsys afero.Fs, options ...func(*pgx.ConnConfig)) error {
	conn, err := utils.ConnectByConfig(ctx, config, options...)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	rows, err := conn.Query(ctx, TotalIndexSizeQuery, reset.LikeEscapeSchema(utils.InternalSchemas))
	if err != nil {
		return errors.Errorf("failed to query rows: %w", err)
	}
	result, err := pgxv5.CollectRows[Result](rows)
	if err != nil {
		return err
	}

	table := "|Size|\n|-|\n"
	for _, r := range result {
		table += fmt.Sprintf("|`%s`|\n", r.Size)
	}
	return list.RenderTable(table)
}
