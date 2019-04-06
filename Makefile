GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

./bin/docx2clearhtml.exe: main.go ./settingsStorage/programSettings.go
	$(GOBUILD) -o ./bin/docx2clearhtml.exe

clean:
	$(GOCLEAN)
	rm ./bin/docx2clearhtml.exe

test:
	$(GOTEST)