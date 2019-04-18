# docx2cleanhtml

## Installation

To install dependencies you need GoLang and dep.

Tutorials on how to do so can you find here:

+ [Install GoLang](https://golang.org/doc/install) (OS X: `brew install go` | Fedora: `dnf install golang`)
+ [Install dep](https://github.com/golang/dep) (OS X: `brew install dep` | Fedora: `dnf install dep`)

## Dependencies

After setting up GoLang your $GOPATH should point into a subfolder of your workspace.
For example my $GOPATH points to `/home/leon/workspace/gospace`

The project should be located in `$GOPATH/src/leong/docx2cleanhtml`.
GoLang establishes this structure for dependency management.

After that you can run `dep ensure` in the root of the project folder to install the dependencies.

## Building

To run the application you use `go run main.go`.
If you want to ship you can use `go build`
or preferably `make build` which outputs the executable in the `bin` folder.

## Example Usage

`docx2cleanhtml -i /path/to/your/file.docx`
