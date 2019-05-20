FROM golang:alpine

ENV DX2HTMLREPO=http://git.doggoat.de/dhalucario/docx2cleanhtml.git
ENV DX2HTMLPATH=/opt/
ENV DX2HTMLPACKAGES=make git npm nodejs

WORKDIR $DX2HTMLPATH

RUN apk add $DX2HTMLPACKAGES
RUN git clone http://git.doggoat.de/dhalucario/docx2cleanhtml.git

WORKDIR $DX2HTMLPATH/docx2cleanhtml

RUN go mod vendor
RUN npm install -g webpack webpack-cli
RUN npm install

RUN make

RUN chmod +x /opt/docx2clearhtml/docx2cleanhtml

ENV PATH=/opt/docx2clearhtml:$PATH

EXPOSE 8001

CMD ["docx2clearhtml"]
#CMD ["docx2clearhtml", "-wsrv", "ip:127.0.0.1", "port:8001"]
#CMD ["ls", "-la"]