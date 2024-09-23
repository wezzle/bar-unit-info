# Bar-unit-info

A CLI tool that provides insight into the many units part of [Beyond All Reason](https://github.com/beyond-all-reason/Beyond-All-Reason).

## How to use

There are a few ways to run the CLI tool:

1. Download the correct binary for your system from the [releases](https://github.com/wezzle/bar-unit-info/releases) page. You can then execute the binary from a terminal application or let your system start one for you by double clicking (MacOS for example).

1. Use Nix flakes to build the binary without checking out this repo: `nix build github:wezzle/bar-unit-info` and run the compiled binary: `./result/bin/bar-unit-info`

1. Checkout this repo and run `nix build` in the root directory, then run the compiled binary: `./result/bin/bar-unit-info`

## Development

This repository uses `nix flakes` to setup a development shell. If you have [direnv](https://direnv.net/) enabled on your shell you will automatically get a development shell with the required dependencies (go and a sparse checkout of the Beyond All Reason main repo). Alternatively when you have nix installed you can run `nix develop` in the root repo to enter a development shell.

If you don't want to use `nix flakes` you can use your own `go` version, just make sure it's above the version listed in the `go.mod`.

This tool uses code generation to embed unit data from the Beyond All Reason main repo in the compiled binary. This data is a representation of a point in time and will not automatically be updated. To run a development version against the latest Beyond All Reason repository you can do either of the following:

### Nix flakes

Update the hash (https://github.com/wezzle/bar-unit-info/blob/main/flake.nix#L35) to `hash = pkgs.lib.fakeHash;` and run `nix develop`. Nix will provide you with the correct hash, for example:

```
error: hash mismatch in fixed-output derivation '/nix/store/b3xq0awnfxkahd5d18jvz65x49sgln84-Beyond-All-Reason.drv':
         specified: sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=
            got:    sha256-M7XEEaCyd7Lk1QmcKBkNSlS8BqSJqEhOGwee+T9Hyes=
```

You then replace the fake hash with the correct hash in flake.nix. After running `nix develop` again you will have a shell with `$GAME_REPO` pointing to the new sparse checkout of the Beyond All Reason main repo.

All that is left to do is run `go generate ./...` to update the generated gamedata files.

### Justfile

The justfile contains two targets to update the gamedata files:

* `just bar-repo` to do a sparse checkout of the latest Beyond All Reason main repo

* `just generate` to generate the updated gamedata files
