FROM golang:latest

RUN mkdir -p /opt/docx2html/public/output
RUN mkdir -p /opt/docx2html/uploads

WORKDIR /opt/docx2html

COPY bin/docx2clearhtml ./
COPY bin/public/bundle.js ./bin/public/
COPY bin/public/bundle.css ./bin/public/
COPY bin/public/index.html ./bin/public/index.html

CMD ["/opt/docx2clearhtml -wsrv ip:127.0.0.1 port:8001"]