FROM golang:1.9.2

RUN mkdir /go/src/truora
COPY . /go/src/truora
RUN /bin/bash -c 'ls -la; chmod +x /go/src/truora; ls -la'
WORKDIR /go/src/truora
RUN  go get truora
RUN go install
RUN which truora
ENV PATH ="/go/bin/:${PATH}"

ENTRYPOINT ["/go/bin/truora"]

EXPOSE 3344
