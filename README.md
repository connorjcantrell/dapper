# Dapper
#### Decentralized Application Manager for the Algorand Blockchain
Dappman is a Golang CLI toolkit for compiling, deploying, and managing Algorand applications. It is a thin wrapper built on top of [goal app](https://developer.algorand.org/docs/clis/goal/app/app/)

## Fetch from GitHub
The easiest is to clone Dappman in a directory outside of GOPATH, as in the following example:
```
mkdir $HOME/src
cd $HOME/src
git clone https://github.com/connorjcantrell/dapper.git
cd dapper
go install
```

## Create Environment Variables
- `ALGORAND_PASSPHRASE`
- `ALGOD_ADDRESS`
- `ALGOD_TOKEN`
- `DAPPER_TEMPLATES`

## Getting Started
Initialize dapper inside a project directory:
```
dapper init --global-byteslices 0 -global-ints 0 --local-byteslices 0 --local-ints 0 --boilerplate pyteal
```

This will generate the following project structure:
```
project
│
└───.dapper
    │   config.json 
│
└───public
│   
└───src
    │   approval_program.py
    │   clear_state_program.py
```

`.config.json` is a local representation of the application details that exist on the Algorand Blockchain. This file will be referenced/ modified during `create`, `update` and `delete` commands. 
Application ID is initially set to `0` to signify the app has not yet been created. 

**Do not manually modify `.config.json`** 
#### `.config.json`
```
{
	"application_id": 0,
	"block": 0,
	"creator": "",
	"revision": 0,
	"deleted": false,
	"global_state_schema": {
		"num_byte_slice": 0,
		"num_uint": 10
	},
	"local_state_schema": {
		"num_byte_slice": 0,
		"num_uint": 0
	}
}

```

## Compile
```
dapper compile pyteal
```
1. Searches for `approval_program.py` and `clear_state_program.py` in `/src` directory
2. Compiles PyTeal down to TEAL, writes TEAL programs to `/public` directory


```
project
│
└───.dapper
    │   config.json 
│
└───public
│   │   approval_program.teal
│   │   clear_state_program.teal
│
└───src
    │   approval_program.py
    │   clear_state_program.py
```

## Create
Issue a transaction that creates an application
```
dapper create
```
#### `.config.json` modifications
`application_id` will be changed from `0`
`revision` will be increased by `1`


## Update
Issue a transaction that updates an application's ApprovalProgram and ClearStateProgram
```
dapper update
```

#### `.config.json` modifications
`revision` will be increased by `1`


## Delete
```
dapper delete
```
#### `.config.json` modifications
`deleted` changed to `true`

