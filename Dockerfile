FROM golang:1.19 as base

RUN useradd -m praveen


FROM base as dev

RUN apt-get update 

# image processing library 
RUN apt-get install -y libvips-dev       

WORKDIR /src/app
COPY . .

RUN chown praveen /src/app

USER praveen


EXPOSE 8080


RUN go build -tags netgo -ldflags '-s -w' -o shopee

CMD [ "./shopee" ]