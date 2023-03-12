# go-cdcl
    
CDCL SAT solver made in Go designed for [DIMCAS](http://logic.pdmi.ras.ru/~basolver/dimacs.html) problems

> This project was made as part of the COMP60332 Automated Reasoning and Verification Course as part of the MSc Advanced Computer Science course at the Univeristy of Manchester.

## Get Started


### Requirements

- Built for Ubuntu/Debian based Linux, might work in other linux distributions but might be unstable

### Steps

1. Download Binary

2. Give execution permission
```bash
chmod +x ./gocdcl
```
3. Run `gocdcl -h` to open help menu
```bash
$ ./gocdcl -h
NAME:
   gocdcl - Pass SAT file as stdin pipe or using the -f/--file flag to run SAT solver

USAGE:
   gocdcl [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --file value, -f value  .SAT file to be processed. This option is overridden if input provided by stdin pipe
   --verbose, -v           Switches on detailed logging for cdcl solver (default: false)
   --experimental, -e      use experimental features (default: false)
   --help, -h              show help
```

You can pass your DIMCAS format files via stdin or using the `-f` flag.

## Building From Source

### Requirements**

- Go Version 1.20

### Steps

1. Clone Repository and open directory
```bash
git clone https://github.com/alanpjohn/go-cdcl.git && cd go-cdcl
```

2. Run Makefile
```bash
make
```

3. You will find the binary generated with filename `gocdcl`
