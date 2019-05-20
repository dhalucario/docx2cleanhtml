GOCMD=go
GOBUILD=$(GOCMD) build --mod=vendor
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

NPM=npm
NPMINSTALL=$(NPM) install

WEBPACKCMD=webpack
WEBPACKBUILD=$(WEBPACKCMD) --config webpack-dist.config.js

./bin/docx2clearhtml : ./bin/public/bundle.js ./bin/public/bundle.css
	$(GOBUILD) -o ./bin/docx2clearhtml
./bin/public/bundle.js ./bin/public/bundle.css: ./web/index.js ./web/js/main.js ./package.json ./webpack-dist.config.js
	$(NPMINSTALL)
	$(WEBPACKBUILD)

clean:
	$(GOCLEAN)
	rm ./bin/docx2clearhtml
	rm ./bin/public/bundle.js
	rm ./bin/public/bundle.css

test:
	$(GOTEST)