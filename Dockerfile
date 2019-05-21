FROM golang:alpine

ENV DX2HTMLREPO="https://git.doggoat.de/dhalucario/docx2cleanhtml.git"
ENV DX2HTMLPATH="/opt"
ENV DX2HTMLPACKAGES="make git npm nodejs"

WORKDIR $DX2HTMLPATH

RUN apk add $DX2HTMLPACKAGES
RUN git clone http://git.doggoat.de/dhalucario/docx2cleanhtml.git

WORKDIR $DX2HTMLPATH/docx2cleanhtml

RUN go mod vendor
RUN npm install -g webpack webpack-cli
RUN npm install

RUN make

RUN rm $DX2HTMLPATH/docx2cleanhtml/node_modules -rf
RUN chmod +x $DX2HTMLPATH/docx2cleanhtml/bin/docx2cleanhtml
ENV PATH=$DX2HTMLPATH/docx2cleanhtml/bin:$PATH

RUN apk del $DX2HTMLPACKAGES

EXPOSE 8001

WORKDIR $DX2HTMLPATH/docx2cleanhtml/bin
CMD ["docx2cleanhtml", "-wsrv", "ip:0.0.0.0", "port:8001"]