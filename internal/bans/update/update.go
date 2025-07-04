package update

import (
	"context"
	"fmt"
	"net"

	"github.com/go-errors/errors"
	"github.com/spf13/afero"
	"github.com/noTreeTeam/cli/internal/utils"
	"github.com/noTreeTeam/cli/pkg/api"
)

func validateIps(ips []string) error {
	for _, ip := range ips {
		if net.ParseIP(ip) == nil {
			return errors.Errorf("invalid IP address: %s", ip)
		}
	}
	return nil
}

func Run(ctx context.Context, projectRef string, dbIpsToUnban []string, fsys afero.Fs) error {
	// 1. sanity checks
	if err := validateIps(dbIpsToUnban); err != nil {
		return err
	}

	// 2. remove bans
	{
		resp, err := utils.GetSupabase().V1DeleteNetworkBansWithResponse(ctx, projectRef, api.RemoveNetworkBanRequest{
			Ipv4Addresses: dbIpsToUnban,
		})
		if err != nil {
			return errors.Errorf("failed to remove network bans: %w", err)
		}
		if resp.StatusCode() != 200 {
			return errors.New("Unexpected error removing network bans: " + string(resp.Body))
		}
		fmt.Printf("Successfully removed bans for %+v.\n", dbIpsToUnban)
		return nil
	}
}
