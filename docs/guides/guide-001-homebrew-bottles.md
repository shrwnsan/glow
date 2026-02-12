# Homebrew Bottles for Fork Distribution

This guide covers generating and distributing Homebrew bottles for the forked `glow` project. Bottles provide pre-built binaries for faster user installations.

## Overview

There are two approaches to distributing Homebrew bottles:

1. **Manual bottle generation** - Using Homebrew's built-in tools
2. **GoReleaser automation** - Automated release pipeline

## Manual Bottle Generation

When to use: Occasional releases, one-off distributions

### Prerequisites

- Homebrew installed
- `gh` CLI authenticated with GitHub
- Write access to `shrwnsan/glow` and `shrwnsan/homebrew-tap`

### Process

1. **Install with build-bottle flag**

   ```bash
   brew install --build-bottle shrwnsan/tap/glow
   ```

2. **Generate bottle**

   ```bash
   brew bottle shrwnsan/tap/glow \
     --root-url=https://github.com/shrwnsan/glow/releases/download/v{VERSION}
   ```

   This creates a bottle file (e.g., `glow--{VERSION}.{arch}_{os}.bottle.tar.gz`)

3. **Rename bottle** (if needed)

   Homebrew may name bottles with internal OS names (e.g., `tahoe`).
   Rename to match formula expectations (e.g., `sequoia`):

   ```bash
   mv glow--2.1.2.arm64_tahoe.bottle.1.tar.gz \
      glow-2.1.2.arm64_sequoia.bottle.tar.gz
   ```

4. **Upload to GitHub release**

   ```bash
   gh release upload v{VERSION} --repo shrwnsan/glow \
     ./glow-{VERSION}.{arch}_{os}.bottle.tar.gz
   ```

5. **Update formula with bottle block**

   Edit `~/Developer/personal/homebrew-tap/Formula/glow.rb`:

   ```ruby
   bottle do
     root_url "https://github.com/shrwnsan/glow/releases/download/v{VERSION}"
     sha256 cellar: :any_skip_relocation, arm64_sequoia: "{SHA256}"
   end
   ```

   Get SHA256 from bottle JSON output or generate with `shasum -a 256`.

6. **Commit and push formula**

   ```bash
   cd ~/Developer/personal/homebrew-tap
   git add Formula/glow.rb
   git commit -m "feat(glow): add {os}_{arch} bottle for v{VERSION}"
   git push
   ```

### Verification

Users can now install instantly without compilation:

```bash
brew install shrwnsan/tap/glow
```

## GoReleaser Automation (TODO)

When to use: Frequent releases, multiple platforms/architectures

### Prerequisites

- GoReleaser installed: `brew install goreleaser`
- GitHub token with `repo` scope
- `.goreleaser-fork.yml` configured

### Setup

1. **Install GoReleaser**

   ```bash
   brew install goreleaser
   ```

2. **Configure GitHub token**

   ```bash
   export GITHUB_TOKEN=$(gh auth token)
   ```

   Or store in `~/.config/goreleaser/goreleaser.yml`:
   ```yaml
   github:
     token: ${{ secrets.GITHUB_TOKEN }}
   ```

### Release Process

```bash
# Tag and push
git tag v{VERSION}
git push origin v{VERSION}

# Run GoReleaser
goreleaser release --config .goreleaser-fork.yml --clean
```

GoReleaser will:
- Build binaries for all configured platforms
- Create tar.gz archives
- Generate Homebrew bottles
- Upload assets to GitHub release
- Auto-commit updated formula to `shrwnsan/homebrew-tap`

## Troubleshooting

### Bottle not found

If users see "Failed to download resource" errors:

1. Verify bottle exists on release:
   ```bash
   gh release view v{VERSION} --repo shrwnsan/glow --json assets
   ```

2. Check formula `root_url` matches release URL
3. Verify SHA256 matches the uploaded file

### Formula updates not visible

Users may need to update tap:

```bash
brew update
```

### Build-bottle flag not recognized

Use fully-qualified formula name:

```bash
brew install --build-bottle shrwnsan/tap/glow
```

## References

- [Homebrew Bottles Documentation](https://docs.brew.sh/Bottles)
- [GoReleaser Homebrew Integration](https://goreleaser.com/customization/homebrew/)
- [`.goreleaser-fork.yml`](../.goreleaser-fork.yml) - Current configuration
