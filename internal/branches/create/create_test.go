package create

import (
	"context"
	"net"
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/noTreeTeam/cli/internal/testing/apitest"
	"github.com/noTreeTeam/cli/internal/utils"
	"github.com/noTreeTeam/cli/internal/utils/flags"
	"github.com/noTreeTeam/cli/pkg/api"
	"github.com/noTreeTeam/cli/pkg/cast"
)

func TestCreateCommand(t *testing.T) {
	// Setup valid project ref
	flags.ProjectRef = apitest.RandomProjectRef()

	t.Run("creates preview branch", func(t *testing.T) {
		// Setup valid access token
		token := apitest.RandomAccessToken(t)
		t.Setenv("SUPABASE_ACCESS_TOKEN", string(token))
		// Setup in-memory fs
		fsys := afero.NewMemMapFs()
		// Setup mock api
		defer gock.OffAll()
		gock.New(utils.DefaultApiHost).
			Post("/v1/projects/" + flags.ProjectRef + "/branches").
			Reply(http.StatusCreated).
			JSON(api.BranchResponse{
				Id: "test-uuid",
			})
		// Run test
		err := Run(context.Background(), api.CreateBranchBody{
			Region: cast.Ptr("sin"),
		}, fsys)
		// Check error
		assert.NoError(t, err)
	})

	t.Run("throws error on network disconnected", func(t *testing.T) {
		// Setup in-memory fs
		fsys := afero.NewMemMapFs()
		require.NoError(t, afero.WriteFile(fsys, utils.ProjectRefPath, []byte(flags.ProjectRef), 0644))
		// Setup mock api
		defer gock.OffAll()
		gock.New(utils.DefaultApiHost).
			Post("/v1/projects/" + flags.ProjectRef + "/branches").
			ReplyError(net.ErrClosed)
		// Run test
		err := Run(context.Background(), api.CreateBranchBody{
			Region: cast.Ptr("sin"),
		}, fsys)
		// Check error
		assert.ErrorIs(t, err, net.ErrClosed)
	})

	t.Run("throws error on service unavailable", func(t *testing.T) {
		// Setup in-memory fs
		fsys := afero.NewMemMapFs()
		require.NoError(t, afero.WriteFile(fsys, utils.ProjectRefPath, []byte(flags.ProjectRef), 0644))
		// Setup mock api
		defer gock.OffAll()
		gock.New(utils.DefaultApiHost).
			Post("/v1/projects/" + flags.ProjectRef + "/branches").
			Reply(http.StatusServiceUnavailable)
		// Run test
		err := Run(context.Background(), api.CreateBranchBody{
			Region: cast.Ptr("sin"),
		}, fsys)
		// Check error
		assert.ErrorContains(t, err, "Unexpected error creating preview branch:")
	})
}
