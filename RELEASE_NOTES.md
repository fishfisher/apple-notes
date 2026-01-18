# How to Release

## First Release (v0.1.0)
```bash
# Make sure you're on main with latest changes
git pull

# Create and push a tag
git tag -a v0.1.0 -m "Initial release: Apple Notes CLI with smart edit protection"
git push origin v0.1.0

# Run goreleaser (requires HOMEBREW_TAP_TOKEN to be set)
goreleaser release --clean
```

## Subsequent Releases
```bash
git tag -a v0.2.0 -m "Release v0.2.0: [describe changes]"
git push origin v0.2.0
goreleaser release --clean
```

## Test Release Without Publishing
```bash
goreleaser release --snapshot --clean
```

This will:
1. Build binaries for macOS (Intel + ARM)
2. Create GitHub release
3. Update homebrew-tap with new formula
4. Users can install with: `brew install fishfisher/tap/apple-notes`
