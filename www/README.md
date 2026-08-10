# DENIS Website

> [!NOTE]
> Unlike the rest of DENIS, everything in `www/` was written by AI (Claude). The DNS server itself remains hand-written.

The DENIS docs site, built with [Bloom](https://github.com/maytees/bloom) — a hand-written signals library. The Bloom runtime is vendored into `www/vendor/bloom` (see its README), so the site builds standalone; nothing outside this repo is needed. To pull in a newer Bloom from a sibling `../bloom` checkout, run `just site-vendor`.

**Note: Bun has issues with segfaults, if the dev server crashes for no reason, just start it up again.**

Install dependencies:

```bash
# from www/
bun install
```

Run the tailwind server in watch mode:

```bash
# from project root
just site-css
```

Then start the dev server:

```bash
# from project root
just site-dev
```

Building:

```bash
# from project root
just site-build
```

Output will be in `www/dist`.
