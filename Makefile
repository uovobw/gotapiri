GO=go
MAIN=gotapiri.go
TARGET=gotapiri

all:
	${GO} build github.com/uovobw/gotapiri

clean:
	rm -f ${TARGET} 
