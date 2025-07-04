package keys

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/go-errors/errors"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/afero"
	"github.com/noTreeTeam/cli/internal/utils"
	"github.com/noTreeTeam/cli/pkg/config"
)

type CustomName struct {
	DbHost         string `env:"db.host,default=NEXT_PUBLIC_SUPABASE_URL"`
	DbPassword     string `env:"db.password,default=SUPABASE_DB_PASSWORD"`
	JWTSecret      string `env:"db.password,default=SUPABASE_AUTH_JWT_SECRET"`
	AnonKey        string `env:"auth.anon_key,default=SUPABASE_AUTH_ANON_KEY"`
	ServiceRoleKey string `env:"auth.service_role_key,default=SUPABASE_AUTH_SERVICE_ROLE_KEY"`
}

func Run(ctx context.Context, projectRef, format string, names CustomName, fsys afero.Fs) error {
	branch := GetGitBranch(fsys)
	if err := GenerateSecrets(ctx, projectRef, branch, fsys); err != nil {
		return err
	}
	return utils.EncodeOutput(format, os.Stdout, map[string]string{
		names.DbHost:         fmt.Sprintf("%s-%s.fly.dev", projectRef, branch),
		names.DbPassword:     utils.Config.Db.Password,
		names.JWTSecret:      utils.Config.Auth.JwtSecret.Value,
		names.AnonKey:        utils.Config.Auth.AnonKey.Value,
		names.ServiceRoleKey: utils.Config.Auth.ServiceRoleKey.Value,
	})
}

func GenerateSecrets(ctx context.Context, projectRef, branch string, fsys afero.Fs) error {
	// Load JWT secret from api
	resp, err := utils.GetSupabase().V1GetPostgrestServiceConfigWithResponse(ctx, projectRef)
	if err != nil {
		return errors.Errorf("failed to get postgrest config: %w", err)
	}
	if resp.JSON200 == nil {
		return errors.New("Unexpected error retrieving JWT secret: " + string(resp.Body))
	}
	utils.Config.Auth.JwtSecret.Value = *resp.JSON200.JwtSecret
	// Generate database password
	key := strings.Join([]string{
		projectRef,
		utils.Config.Auth.JwtSecret.Value,
		branch,
	}, ":")
	hash := sha256.Sum256([]byte(key))
	utils.Config.Db.Password = hex.EncodeToString(hash[:])
	// Generate JWT tokens
	anonToken := config.CustomClaims{
		Issuer: "supabase",
		Ref:    projectRef,
		Role:   "anon",
	}.NewToken()
	if utils.Config.Auth.AnonKey.Value, err = anonToken.SignedString([]byte(utils.Config.Auth.JwtSecret.Value)); err != nil {
		return errors.Errorf("failed to sign anon key: %w", err)
	}
	serviceToken := config.CustomClaims{
		Issuer: "supabase",
		Ref:    projectRef,
		Role:   "service_role",
	}.NewToken()
	if utils.Config.Auth.ServiceRoleKey.Value, err = serviceToken.SignedString([]byte(utils.Config.Auth.JwtSecret.Value)); err != nil {
		return errors.Errorf("failed to sign service_role key: %w", err)
	}
	return nil
}

func GetGitBranch(fsys afero.Fs) string {
	return GetGitBranchOrDefault("main", fsys)
}

func GetGitBranchOrDefault(def string, fsys afero.Fs) string {
	head := os.Getenv("GITHUB_HEAD_REF")
	if len(head) > 0 {
		return head
	}
	opts := &git.PlainOpenOptions{DetectDotGit: true}
	if repo, err := git.PlainOpenWithOptions(".", opts); err == nil {
		if ref, err := repo.Head(); err == nil {
			return ref.Name().Short()
		}
	}
	return def
}
