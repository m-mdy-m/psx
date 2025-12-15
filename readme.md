# PSX - Project Structure Checker

yeah so basically this checks if your project has the right structure or not. like if you forgot to add a README or LICENSE or whatever, it'll tell you.

## what it does

- checks if you have basic files (README, LICENSE, etc)
- looks at your project type (node, go, rust, whatever)
- tells you what's missing
- can fix some stuff automatically if you want
- works on linux, mac, windows

## install

### quick way
```bash
curl -sSL https://raw.githubusercontent.com/m-mdy-m/psx/main/scripts/install.sh | bash
```

### or download binary
go to releases page and download the one for your OS

### or build it yourself
```bash
git clone https://github.com/m-mdy-m/psx
cd psx
make build
```

## how to use

just run it in your project folder:
```bash
psx check
```

it'll tell you what's wrong. if you want more info:
```bash
psx check --verbose
```

fix stuff automatically:
```bash
psx fix
```

or fix stuff one by one (asks you first):
```bash
psx fix --interactive
```

## examples

### check a nodejs project
```bash
cd my-node-app
psx check

# output:
# ✗ README_MISSING
# ✗ tests/ folder not found
# ⚠ No LICENSE file
```

### fix everything
```bash
psx fix

# creates README.md
# creates tests/ folder
# asks which license you want
```

## why i made this

i kept forgetting to add READMEs and licenses to my projects lol. also wanted something simple that just works without tons of config.

## contributing

sure, send PRs. check CONTRIBUTING.md if you care about that stuff.

## license

[MIT](./LICENSE) - do whatever you want

## issues?

open an issue on github or email me: bitsgenix@gmail.com

---

made by [@m-mdy-m](https://github.com/m-mdy-m) | [website](https://m-mdy-m.github.io)
