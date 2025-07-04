package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/noTreeTeam/cli/internal/inspect/bloat"
	"github.com/noTreeTeam/cli/internal/inspect/blocking"
	"github.com/noTreeTeam/cli/internal/inspect/cache"
	"github.com/noTreeTeam/cli/internal/utils"
	"github.com/noTreeTeam/cli/internal/utils/flags"

	"github.com/noTreeTeam/cli/internal/inspect"
	"github.com/noTreeTeam/cli/internal/inspect/calls"
	"github.com/noTreeTeam/cli/internal/inspect/index_sizes"
	"github.com/noTreeTeam/cli/internal/inspect/index_usage"
	"github.com/noTreeTeam/cli/internal/inspect/locks"
	"github.com/noTreeTeam/cli/internal/inspect/long_running_queries"
	"github.com/noTreeTeam/cli/internal/inspect/outliers"
	"github.com/noTreeTeam/cli/internal/inspect/replication_slots"
	"github.com/noTreeTeam/cli/internal/inspect/role_configs"
	"github.com/noTreeTeam/cli/internal/inspect/role_connections"
	"github.com/noTreeTeam/cli/internal/inspect/seq_scans"
	"github.com/noTreeTeam/cli/internal/inspect/table_index_sizes"
	"github.com/noTreeTeam/cli/internal/inspect/table_record_counts"
	"github.com/noTreeTeam/cli/internal/inspect/table_sizes"
	"github.com/noTreeTeam/cli/internal/inspect/total_index_size"
	"github.com/noTreeTeam/cli/internal/inspect/total_table_sizes"
	"github.com/noTreeTeam/cli/internal/inspect/unused_indexes"
	"github.com/noTreeTeam/cli/internal/inspect/vacuum_stats"
)

var (
	inspectCmd = &cobra.Command{
		GroupID: groupLocalDev,
		Use:     "inspect",
		Short:   "Tools to inspect your Supabase project",
	}

	inspectDBCmd = &cobra.Command{
		Use:   "db",
		Short: "Tools to inspect your Supabase database",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			ctx, _ := signal.NotifyContext(cmd.Context(), os.Interrupt)
			cmd.SetContext(ctx)
			return cmd.Root().PersistentPreRunE(cmd, args)
		},
	}

	inspectCacheHitCmd = &cobra.Command{
		Use:   "cache-hit",
		Short: "Show cache hit rates for tables and indices",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cache.Run(cmd.Context(), flags.DbConfig, afero.NewOsFs())
		},
	}

	inspectReplicationSlotsCmd = &cobra.Command{
		Use:   "replication-slots",
		Short: "Show information about replication slots on the database",
		RunE: func(cmd *cobra.Command, args []string) error {
			return replication_slots.Run(cmd.Context(), flags.DbConfig, afero.NewOsFs())
		},
	}

	inspectIndexUsageCmd = &cobra.Command{
		Use:   "index-usage",
		Short: "Show information about the efficiency of indexes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return index_usage.Run(cmd.Context(), flags.DbConfig, afero.NewOsFs())
		},
	}

	inspectLocksCmd = &cobra.Command{
		Use:   "locks",
		Short: "Show queries which have taken out an exclusive lock on a relation",
		RunE: func(cmd *cobra.Command, args []string) error {
			return locks.Run(cmd.Context(), flags.DbConfig, afero.NewOsFs())
		},
	}

	inspectBlockingCmd = &cobra.Command{
		Use:   "blocking",
		Short: "Show queries that are holding locks and the queries that are waiting for them to be released",
		RunE: func(cmd *cobra.Command, args []string) error {
			return blocking.Run(cmd.Context(), flags.DbConfig, afero.NewOsFs())
		},
	}

	inspectOutliersCmd = &cobra.Command{
		Use:   "outliers",
		Short: "Show queries from pg_stat_statements ordered by total execution time",
		RunE: func(cmd *cobra.Command, args []string) error {
			return outliers.Run(cmd.Context(), flags.DbConfig, afero.NewOsFs())
		},
	}

	inspectCallsCmd = &cobra.Command{
		Use:   "calls",
		Short: "Show queries from pg_stat_statements ordered by total times called",
		RunE: func(cmd *cobra.Command, args []string) error {
			return calls.Run(cmd.Context(), flags.DbConfig, afero.NewOsFs())
		},
	}

	inspectTotalIndexSizeCmd = &cobra.Command{
		Use:   "total-index-size",
		Short: "Show total size of all indexes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return total_index_size.Run(cmd.Context(), flags.DbConfig, afero.NewOsFs())
		},
	}

	inspectIndexSizesCmd = &cobra.Command{
		Use:   "index-sizes",
		Short: "Show index sizes of individual indexes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return index_sizes.Run(cmd.Context(), flags.DbConfig, afero.NewOsFs())
		},
	}

	inspectTableSizesCmd = &cobra.Command{
		Use:   "table-sizes",
		Short: "Show table sizes of individual tables without their index sizes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return table_sizes.Run(cmd.Context(), flags.DbConfig, afero.NewOsFs())
		},
	}

	inspectTableIndexSizesCmd = &cobra.Command{
		Use:   "table-index-sizes",
		Short: "Show index sizes of individual tables",
		RunE: func(cmd *cobra.Command, args []string) error {
			return table_index_sizes.Run(cmd.Context(), flags.DbConfig, afero.NewOsFs())
		},
	}

	inspectTotalTableSizesCmd = &cobra.Command{
		Use:   "total-table-sizes",
		Short: "Show total table sizes, including table index sizes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return total_table_sizes.Run(cmd.Context(), flags.DbConfig, afero.NewOsFs())
		},
	}

	inspectUnusedIndexesCmd = &cobra.Command{
		Use:   "unused-indexes",
		Short: "Show indexes with low usage",
		RunE: func(cmd *cobra.Command, args []string) error {
			return unused_indexes.Run(cmd.Context(), flags.DbConfig, afero.NewOsFs())
		},
	}

	inspectSeqScansCmd = &cobra.Command{
		Use:   "seq-scans",
		Short: "Show number of sequential scans recorded against all tables",
		RunE: func(cmd *cobra.Command, args []string) error {
			return seq_scans.Run(cmd.Context(), flags.DbConfig, afero.NewOsFs())
		},
	}

	inspectLongRunningQueriesCmd = &cobra.Command{
		Use:   "long-running-queries",
		Short: "Show currently running queries running for longer than 5 minutes",
		RunE: func(cmd *cobra.Command, args []string) error {
			return long_running_queries.Run(cmd.Context(), flags.DbConfig, afero.NewOsFs())
		},
	}

	inspectTableRecordCountsCmd = &cobra.Command{
		Use:   "table-record-counts",
		Short: "Show estimated number of rows per table",
		RunE: func(cmd *cobra.Command, args []string) error {
			return table_record_counts.Run(cmd.Context(), flags.DbConfig, afero.NewOsFs())
		},
	}

	inspectBloatCmd = &cobra.Command{
		Use:   "bloat",
		Short: "Estimates space allocated to a relation that is full of dead tuples",
		RunE: func(cmd *cobra.Command, args []string) error {
			return bloat.Run(cmd.Context(), flags.DbConfig, afero.NewOsFs())
		},
	}

	inspectVacuumStatsCmd = &cobra.Command{
		Use:   "vacuum-stats",
		Short: "Show statistics related to vacuum operations per table",
		RunE: func(cmd *cobra.Command, args []string) error {
			return vacuum_stats.Run(cmd.Context(), flags.DbConfig, afero.NewOsFs())
		},
	}

	inspectRoleConfigsCmd = &cobra.Command{
		Use:   "role-configs",
		Short: "Show configuration settings for database roles when they have been modified",
		RunE: func(cmd *cobra.Command, args []string) error {
			return role_configs.Run(cmd.Context(), flags.DbConfig, afero.NewOsFs())
		},
	}

	inspectRoleConnectionsCmd = &cobra.Command{
		Use:   "role-connections",
		Short: "Show number of active connections for all database roles",
		RunE: func(cmd *cobra.Command, args []string) error {
			return role_connections.Run(cmd.Context(), flags.DbConfig, afero.NewOsFs())
		},
	}

	outputDir string

	reportCmd = &cobra.Command{
		Use:   "report",
		Short: "Generate a CSV output for all inspect commands",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			if len(outputDir) == 0 {
				defaultPath := filepath.Join(utils.CurrentDirAbs, "report")
				title := fmt.Sprintf("Enter a directory to save output files (or leave blank to use %s): ", utils.Bold(defaultPath))
				if dir, err := utils.NewConsole().PromptText(ctx, title); err != nil {
					return err
				} else if len(dir) == 0 {
					outputDir = defaultPath
				}
			}
			return inspect.Report(ctx, outputDir, flags.DbConfig, afero.NewOsFs())
		},
	}
)

func init() {
	inspectFlags := inspectCmd.PersistentFlags()
	inspectFlags.String("db-url", "", "Inspect the database specified by the connection string (must be percent-encoded).")
	inspectFlags.Bool("linked", true, "Inspect the linked project.")
	inspectFlags.Bool("local", false, "Inspect the local database.")
	inspectCmd.MarkFlagsMutuallyExclusive("db-url", "linked", "local")
	inspectDBCmd.AddCommand(inspectCacheHitCmd)
	inspectDBCmd.AddCommand(inspectReplicationSlotsCmd)
	inspectDBCmd.AddCommand(inspectIndexUsageCmd)
	inspectDBCmd.AddCommand(inspectLocksCmd)
	inspectDBCmd.AddCommand(inspectBlockingCmd)
	inspectDBCmd.AddCommand(inspectOutliersCmd)
	inspectDBCmd.AddCommand(inspectCallsCmd)
	inspectDBCmd.AddCommand(inspectTotalIndexSizeCmd)
	inspectDBCmd.AddCommand(inspectIndexSizesCmd)
	inspectDBCmd.AddCommand(inspectTableSizesCmd)
	inspectDBCmd.AddCommand(inspectTableIndexSizesCmd)
	inspectDBCmd.AddCommand(inspectTotalTableSizesCmd)
	inspectDBCmd.AddCommand(inspectUnusedIndexesCmd)
	inspectDBCmd.AddCommand(inspectSeqScansCmd)
	inspectDBCmd.AddCommand(inspectLongRunningQueriesCmd)
	inspectDBCmd.AddCommand(inspectTableRecordCountsCmd)
	inspectDBCmd.AddCommand(inspectBloatCmd)
	inspectDBCmd.AddCommand(inspectVacuumStatsCmd)
	inspectDBCmd.AddCommand(inspectRoleConfigsCmd)
	inspectDBCmd.AddCommand(inspectRoleConnectionsCmd)
	inspectCmd.AddCommand(inspectDBCmd)
	reportCmd.Flags().StringVar(&outputDir, "output-dir", "", "Path to save CSV files in")
	inspectCmd.AddCommand(reportCmd)
	rootCmd.AddCommand(inspectCmd)
}
