package outliers

import (
	"context"
	_ "embed"
	"fmt"
	"regexp"

	"github.com/go-errors/errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/spf13/afero"
	"github.com/noTreeTeam/cli/internal/migration/list"
	"github.com/noTreeTeam/cli/internal/utils"
	"github.com/noTreeTeam/cli/pkg/pgxv5"
)

//go:embed outliers.sql
var OutliersQuery string

type Result struct {
	Total_exec_time string
	Prop_exec_time  string
	Ncalls          string
	Sync_io_time    string
	Query           string
}

func Run(ctx context.Context, config pgconn.Config, fsys afero.Fs, options ...func(*pgx.ConnConfig)) error {
	conn, err := utils.ConnectByConfig(ctx, config, options...)
	if err != nil {
		return err
	}
	defer conn.Close(context.Background())
	rows, err := conn.Query(ctx, OutliersQuery)
	if err != nil {
		return errors.Errorf("failed to query rows: %w", err)
	}
	result, err := pgxv5.CollectRows[Result](rows)
	if err != nil {
		return err
	}
	// TODO: implement a markdown table marshaller
	table := "|Query|Execution Time|Proportion of exec time|Number Calls|Sync IO time|\n|-|-|-|-|-|\n"
	for _, r := range result {
		re := regexp.MustCompile(`\s+|\r+|\n+|\t+|\v`)
		query := re.ReplaceAllString(r.Query, " ")

		re = regexp.MustCompile(`\|`)
		query = re.ReplaceAllString(query, `\|`)
		table += fmt.Sprintf("|`%s`|`%s`|`%s`|`%s`|`%s`|\n", query, r.Total_exec_time, r.Prop_exec_time, r.Ncalls, r.Sync_io_time)
	}
	return list.RenderTable(table)
}
