package seq_scans

import (
	"context"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/noTreeTeam/cli/internal/db/reset"
	"github.com/noTreeTeam/cli/internal/utils"
	"github.com/noTreeTeam/cli/pkg/pgtest"
)

var dbConfig = pgconn.Config{
	Host:     "127.0.0.1",
	Port:     5432,
	User:     "admin",
	Password: "password",
	Database: "postgres",
}

func TestSequentialScansCommand(t *testing.T) {
	t.Run("inspects sequential scans", func(t *testing.T) {
		// Setup in-memory fs
		fsys := afero.NewMemMapFs()
		// Setup mock postgres
		conn := pgtest.NewConn()
		defer conn.Close(t)
		conn.Query(SeqScansQuery, reset.LikeEscapeSchema(utils.InternalSchemas)).
			Reply("SELECT 1", Result{
				Name:  "test_table",
				Count: 99999,
			})
		// Run test
		err := Run(context.Background(), dbConfig, fsys, conn.Intercept)
		// Check error
		assert.NoError(t, err)
	})
}
