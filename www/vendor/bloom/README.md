# bloom (vendored)

A vendored copy of the [Bloom](https://github.com/maytees/bloom) runtime (`packages/runtime`, tests excluded), so this repo builds standalone — on GitHub Actions or a fresh clone — without the bloom repo checked out next to it.

Bun compiles the TypeScript source at build time, so no pre-compiled bundle is needed.

To update it from a sibling `../bloom` checkout, run `just site-vendor` from the repo root.
