# Supabase CLI (noTreeTeam Fork)

[![Coverage Status](https://coveralls.io/repos/github/supabase/cli/badge.svg?branch=main)](https://coveralls.io/github/supabase/cli?branch=main) [![Bitbucket Pipelines](https://img.shields.io/bitbucket/pipelines/supabase-cli/setup-cli/master?style=flat-square&label=Bitbucket%20Canary)](https://bitbucket.org/supabase-cli/setup-cli/pipelines) [![Gitlab Pipeline Status](https://img.shields.io/gitlab/pipeline-status/sweatybridge%2Fsetup-cli?label=Gitlab%20Canary)
](https://gitlab.com/sweatybridge/setup-cli/-/pipelines)

This is a fork of the [Supabase CLI](https://github.com/supabase/cli) with custom modifications for the noTreeTeam organization.

## Key Differences from Upstream

- **Custom Postgres Image**: Uses `ghcr.io/notreeteam/postgres:latest` instead of the default Supabase Postgres image
- **Published to GitHub Packages**: Available as `@notreeteam/supabase-cli` on GitHub Packages
- **Automated Releases**: Regularly updated to track upstream changes

## Getting started

### Install the CLI (noTreeTeam Fork)

#### Via GitHub Packages (Recommended)

1. Create or update your `.npmrc` file to configure the GitHub Packages registry:

```
@notreeteam:registry=https://npm.pkg.github.com
//npm.pkg.github.com/:_authToken=YOUR_GITHUB_TOKEN
```

2. Generate a GitHub Personal Access Token with `read:packages` scope at: https://github.com/settings/tokens

3. Install the package:

```bash
npm install @notreeteam/supabase-cli --save-dev
```

4. Use the CLI:

```bash
npx supabase --help
```

#### Direct from GitHub Repository

Alternatively, you can install directly from the GitHub repository:

```json
{
  "devDependencies": {
    "supabase": "github:noTreeTeam/cli"
  }
}
```

### Install the Original Supabase CLI

Available via [NPM](https://www.npmjs.com) as dev dependency. To install:

```bash
npm i supabase --save-dev
```

To install the beta release channel:

```bash
npm i supabase@beta --save-dev
```

When installing with yarn 4, you need to disable experimental fetch with the following nodejs config.

```
NODE_OPTIONS=--no-experimental-fetch yarn add supabase
```

You can override the GitHub repository and release used by the installer by setting environment variables.
Specify `SUPABASE_REPO` for the repository and `SUPABASE_VERSION` for the tag name (use `latest` to fetch the most recent release):

```bash
SUPABASE_REPO=myorg/cli SUPABASE_VERSION=latest npm i github:myorg/cli --save-dev
```

Alternatively, reference the GitHub repository directly in your `package.json`
to keep the CLI local to your project:

```json
{
  "devDependencies": {
    "supabase": "github:myorg/cli"
  }
}
```

Then run `npm install` and access the CLI via `npx supabase`.

> **Note**
For Bun versions below v1.0.17, you must add `supabase` as a [trusted dependency](https://bun.sh/guides/install/trusted) before running `bun add -D supabase`.

<details>
  <summary><b>macOS</b></summary>

  Available via [Homebrew](https://brew.sh). To install:

  ```sh
  brew install supabase/tap/supabase
  ```

  To install the beta release channel:
  
  ```sh
  brew install supabase/tap/supabase-beta
  brew link --overwrite supabase-beta
  ```
  
  To upgrade:

  ```sh
  brew upgrade supabase
  ```
</details>

<details>
  <summary><b>Windows</b></summary>

  Available via [Scoop](https://scoop.sh). To install:

  ```powershell
  scoop bucket add supabase https://github.com/supabase/scoop-bucket.git
  scoop install supabase
  ```

  To upgrade:

  ```powershell
  scoop update supabase
  ```
</details>

<details>
  <summary><b>Linux</b></summary>

  Available via [Homebrew](https://brew.sh) and Linux packages.

  #### via Homebrew

  To install:

  ```sh
  brew install supabase/tap/supabase
  ```

  To upgrade:

  ```sh
  brew upgrade supabase
  ```

  #### via Linux packages

  Linux packages are provided in [Releases](https://github.com/noTreeTeam/cli/releases). To install, download the `.apk`/`.deb`/`.rpm`/`.pkg.tar.zst` file depending on your package manager and run the respective commands.

  ```sh
  sudo apk add --allow-untrusted <...>.apk
  ```

  ```sh
  sudo dpkg -i <...>.deb
  ```

  ```sh
  sudo rpm -i <...>.rpm
  ```

  ```sh
  sudo pacman -U <...>.pkg.tar.zst
  ```
</details>

<details>
  <summary><b>Other Platforms</b></summary>

  You can also install the CLI via [go modules](https://go.dev/ref/mod#go-install) without the help of package managers.

  ```sh
  go install github.com/noTreeTeam/cli@latest
  ```

  Add a symlink to the binary in `$PATH` for easier access:

  ```sh
  ln -s "$(go env GOPATH)/bin/cli" /usr/bin/supabase
  ```

  This works on other non-standard Linux distros.
</details>

<details>
  <summary><b>Community Maintained Packages</b></summary>

  Available via [pkgx](https://pkgx.sh/). Package script [here](https://github.com/pkgxdev/pantry/blob/main/projects/supabase.com/cli/package.yml).
  To install in your working directory:

  ```bash
  pkgx install supabase
  ```

  Available via [Nixpkgs](https://nixos.org/). Package script [here](https://github.com/NixOS/nixpkgs/blob/master/pkgs/development/tools/supabase-cli/default.nix).
</details>

### Run the CLI

```bash
supabase bootstrap
```

Or using npx:

```bash
npx supabase bootstrap
```

The bootstrap command will guide you through the process of setting up a Supabase project using one of the [starter](https://github.com/supabase-community/supabase-samples/blob/main/samples.json) templates.
### GitHub Actions

Use the bundled action to install the CLI from any GitHub repository.
Default repository is `noTreeTeam/cli` so the action can replace the official setup step.

```yaml
steps:
  - uses: actions/checkout@v4
  - uses: ./.github/actions/setup-cli
    with:
      repo: <owner>/<repo>
      version: latest
  - run: supabase --version
```


## Docs

Command & config reference can be found [here](https://supabase.com/docs/reference/cli/about).

## Breaking changes

We follow semantic versioning for changes that directly impact CLI commands, flags, and configurations.

However, due to dependencies on other service images, we cannot guarantee that schema migrations, seed.sql, and generated types will always work for the same CLI major version. If you need such guarantees, we encourage you to pin a specific version of CLI in package.json.

## Developing

To run from source:

```sh
# Go >= 1.22
go run . help
```
