package create

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
	"github.com/noTreeTeam/cli/internal/testing/apitest"
	"github.com/noTreeTeam/cli/internal/utils"
	"github.com/noTreeTeam/cli/pkg/api"
)

func TestOrganizationCreateCommand(t *testing.T) {
	orgName := "Test Organization"

	t.Run("create an organization", func(t *testing.T) {
		// Setup valid access token
		token := apitest.RandomAccessToken(t)
		t.Setenv("SUPABASE_ACCESS_TOKEN", string(token))
		// Flush pending mocks after test execution
		defer gock.OffAll()
		gock.New(utils.DefaultApiHost).
			Post("/v1/organizations").
			Reply(http.StatusCreated).
			JSON(api.OrganizationResponseV1{
				Id:   "combined-fuchsia-lion",
				Name: orgName,
			})
		// Run test
		assert.NoError(t, Run(context.Background(), orgName))
		// Validate api
		assert.Empty(t, apitest.ListUnmatchedRequests())
	})

	t.Run("throws error on network error", func(t *testing.T) {
		// Setup valid access token
		token := apitest.RandomAccessToken(t)
		t.Setenv("SUPABASE_ACCESS_TOKEN", string(token))
		// Flush pending mocks after test execution
		defer gock.OffAll()
		gock.New(utils.DefaultApiHost).
			Post("/v1/organizations").
			ReplyError(errors.New("network error"))
		// Run test
		assert.Error(t, Run(context.Background(), orgName))
		// Validate api
		assert.Empty(t, apitest.ListUnmatchedRequests())
	})

	t.Run("throws error on server unavailable", func(t *testing.T) {
		// Setup valid access token
		token := apitest.RandomAccessToken(t)
		t.Setenv("SUPABASE_ACCESS_TOKEN", string(token))
		// Flush pending mocks after test execution
		defer gock.OffAll()
		gock.New(utils.DefaultApiHost).
			Post("/v1/organizations").
			Reply(http.StatusServiceUnavailable).
			JSON(map[string]string{"message": "unavailable"})
		// Run test
		assert.Error(t, Run(context.Background(), orgName))
		// Validate api
		assert.Empty(t, apitest.ListUnmatchedRequests())
	})
}
