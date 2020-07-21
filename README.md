# StarchainGo

![StarchainGo](https://github.com/lambda-mike/starchain-go/workflows/Starchain%20Go/badge.svg)

This is a simplified in-memory blockchain implementation written in Go based on the first project from Udacity _Blockchain Developer Nanodegree Program_ course. Udacity kindly provided instructions and boilerplate code - this implementation is based on those great learning resources.

Disclaimer: this is NOT perfect, idiomatic Go, I treat this project as an excuse to learn Go lang and experiment.


## Build

`go build` - builds all packages, produces single executable in root dir: _starchain_


## Test

`go test ./...` - test all packages

`go test -v ./...` - test all packages and print verbose output

`go test -v ./... -run GetBlocks` - test all packages and print verbose output, filter tests to be executed with _run_ flag


## Play

Helpful screenshots can be found in _screenshots/_ directory, where you can find examples of how to query the API. The same example queries can be found in _test.sh_ file. Remember to edit them appropriately - some of them will fail if not changed due to validations. If you change wallet address, change it in all places - otherwise, validations will fail. The timestamp has to be fresh as well (not older than 5 mins).

Sketch of an example scenario:

- check if api works by calling `/hello` endpoint

- get genesis block: `/block/0`

- request message by calling `/requestValidation` endpoint

- submit new star to blockchain by calling `/submitStar` endpoint

- get blocks for a given address by calling `/blocks/:addr` endpoint

Again, you can find examples of queries above in **test.sh** file.
You might find it helpful to edit them and execute interactively in shell, one by one.
