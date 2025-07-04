package bloat

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

//go:embed bloat.sql
var BloatQuery string

type Result struct {
	Type        string
	Schemaname  string
	Object_name string
	Bloat       string
	Waste       string
}

func Run(ctx context.Context, config pgconn.Config, fsys afero.Fs, options ...func(*pgx.ConnConfig)) error {
	conn, err := utils.ConnectByConfig(ctx, config, options...)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	rows, err := conn.Query(ctx, BloatQuery, reset.LikeEscapeSchema(utils.InternalSchemas))
	if err != nil {
		return errors.Errorf("failed to query rows: %w", err)
	}
	result, err := pgxv5.CollectRows[Result](rows)
	if err != nil {
		return err
	}

	table := "|Type|Schema name|Object name|Bloat|Waste\n|-|-|-|-|-|\n"
	for _, r := range result {
		table += fmt.Sprintf("|`%s`|`%s`|`%s`|`%s`|`%s`|\n", r.Type, r.Schemaname, r.Object_name, r.Bloat, r.Waste)
	}
	return list.RenderTable(table)
}
