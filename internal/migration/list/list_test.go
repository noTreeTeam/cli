package list

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/noTreeTeam/cli/internal/testing/fstest"
	"github.com/noTreeTeam/cli/internal/utils"
	"github.com/noTreeTeam/cli/pkg/migration"
	"github.com/noTreeTeam/cli/pkg/pgtest"
)

var dbConfig = pgconn.Config{
	Host:     "127.0.0.1",
	Port:     5432,
	User:     "admin",
	Password: "password",
	Database: "postgres",
}

func TestMigrationList(t *testing.T) {
	t.Run("lists remote migrations", func(t *testing.T) {
		// Setup in-memory fs
		fsys := afero.NewMemMapFs()
		// Setup mock postgres
		conn := pgtest.NewConn()
		defer conn.Close(t)
		conn.Query(migration.LIST_MIGRATION_VERSION).
			Reply("SELECT 0")
		// Run test
		err := Run(context.Background(), dbConfig, fsys, conn.Intercept)
		// Check error
		assert.NoError(t, err)
	})

	t.Run("throws error on remote failure", func(t *testing.T) {
		// Setup in-memory fs
		fsys := afero.NewMemMapFs()
		// Run test
		err := Run(context.Background(), pgconn.Config{}, fsys)
		// Check error
		assert.ErrorContains(t, err, "invalid port (outside range)")
	})

	t.Run("throws error on open failure", func(t *testing.T) {
		// Setup in-memory fs
		fsys := &fstest.OpenErrorFs{DenyPath: utils.MigrationsDir}
		// Setup mock postgres
		conn := pgtest.NewConn()
		defer conn.Close(t)
		conn.Query(migration.LIST_MIGRATION_VERSION).
			Reply("SELECT 0")
		// Run test
		err := Run(context.Background(), dbConfig, fsys, conn.Intercept)
		// Check error
		assert.ErrorIs(t, err, os.ErrPermission)
	})
}

func TestRemoteMigrations(t *testing.T) {
	t.Run("loads migration versions", func(t *testing.T) {
		// Setup mock postgres
		conn := pgtest.NewConn()
		defer conn.Close(t)
		conn.Query(migration.LIST_MIGRATION_VERSION).
			Reply("SELECT 1", []interface{}{"20220727064247"})
		// Run test
		versions, err := loadRemoteVersions(context.Background(), dbConfig, conn.Intercept)
		// Check error
		assert.NoError(t, err)
		assert.ElementsMatch(t, []string{"20220727064247"}, versions)
	})

	t.Run("throws error on connect failure", func(t *testing.T) {
		// Run test
		_, err := loadRemoteVersions(context.Background(), pgconn.Config{})
		// Check error
		assert.ErrorContains(t, err, "invalid port (outside range)")
	})

	t.Run("loads empty migrations on missing table", func(t *testing.T) {
		// Setup mock postgres
		conn := pgtest.NewConn()
		defer conn.Close(t)
		conn.Query(migration.LIST_MIGRATION_VERSION).
			ReplyError(pgerrcode.UndefinedTable, "relation \"supabase_migrations.schema_migrations\" does not exist")
		// Run test
		versions, err := loadRemoteVersions(context.Background(), dbConfig, conn.Intercept)
		// Check error
		assert.NoError(t, err)
		assert.Empty(t, versions)
	})

	t.Run("throws error on invalid row", func(t *testing.T) {
		// Setup mock postgres
		conn := pgtest.NewConn()
		defer conn.Close(t)
		conn.Query(migration.LIST_MIGRATION_VERSION).
			Reply("SELECT 1", []interface{}{})
		// Run test
		_, err := loadRemoteVersions(context.Background(), dbConfig, conn.Intercept)
		// Check error
		assert.ErrorContains(t, err, "number of field descriptions must equal number of destinations, got 0 and 1")
	})
}

func TestLocalMigrations(t *testing.T) {
	t.Run("loads migration versions", func(t *testing.T) {
		// Setup in-memory fs
		fsys := afero.NewMemMapFs()
		path := filepath.Join(utils.MigrationsDir, "20220727064246_test.sql")
		require.NoError(t, afero.WriteFile(fsys, path, []byte{}, 0644))
		path = filepath.Join(utils.MigrationsDir, "20220727064248_test.sql")
		require.NoError(t, afero.WriteFile(fsys, path, []byte{}, 0644))
		// Run test
		versions, err := LoadLocalVersions(fsys)
		// Check error
		assert.NoError(t, err)
		assert.ElementsMatch(t, []string{"20220727064246", "20220727064248"}, versions)
	})

	t.Run("ignores outdated and invalid files", func(t *testing.T) {
		// Setup in-memory fs
		fsys := afero.NewMemMapFs()
		path := filepath.Join(utils.MigrationsDir, "20211208000000_init.sql")
		require.NoError(t, afero.WriteFile(fsys, path, []byte{}, 0644))
		path = filepath.Join(utils.MigrationsDir, "20211208000001_invalid.ts")
		require.NoError(t, afero.WriteFile(fsys, path, []byte{}, 0644))
		// Run test
		versions, err := LoadLocalVersions(fsys)
		// Check error
		assert.NoError(t, err)
		assert.Empty(t, versions)
	})

	t.Run("throws error on open failure", func(t *testing.T) {
		// Setup in-memory fs
		fsys := &fstest.OpenErrorFs{DenyPath: utils.MigrationsDir}
		// Run test
		_, err := LoadLocalVersions(fsys)
		// Check error
		assert.ErrorIs(t, err, os.ErrPermission)
	})
}

func TestMakeTable(t *testing.T) {
	t.Run("tabulate version", func(t *testing.T) {
		// Run test
		table := makeTable([]string{"0", "2"}, []string{"0", "1"})
		// Check error
		lines := strings.Split(strings.TrimSpace(table), "\n")
		assert.ElementsMatch(t, []string{
			"|Local|Remote|Time (UTC)|",
			"|-|-|-|",
			"|`0`|`0`|`0`|",
			"|`1`|` `|`1`|",
			"|` `|`2`|`2`|",
		}, lines)
	})

	t.Run("tabulate timestamp", func(t *testing.T) {
		// Run test
		table := makeTable([]string{"20220727064246", "20220727064248"}, []string{"20220727064246", "20220727064247"})
		// Check error
		lines := strings.Split(strings.TrimSpace(table), "\n")
		assert.ElementsMatch(t, []string{
			"|Local|Remote|Time (UTC)|",
			"|-|-|-|",
			"|`20220727064246`|`20220727064246`|`2022-07-27 06:42:46`|",
			"|`20220727064247`|` `|`2022-07-27 06:42:47`|",
			"|` `|`20220727064248`|`2022-07-27 06:42:48`|",
		}, lines)
	})

	t.Run("ignores string values", func(t *testing.T) {
		// Run test
		table := makeTable([]string{"a", "c"}, []string{"a", "b"})
		// Check error
		lines := strings.Split(strings.TrimSpace(table), "\n")
		assert.ElementsMatch(t, []string{
			"|Local|Remote|Time (UTC)|",
			"|-|-|-|",
		}, lines)
	})
}
