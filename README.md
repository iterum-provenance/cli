# iterum-cli
The CLI for the iterum tool


## Makefile commands
* `make build` builds the go application in `./build/iterum`
* `make link` creates a symlink to the created build folder making `iterum` accessible via the terminal
* `make clean` removes the symlink pointing to `/usr/bin/iterum` and the build folder