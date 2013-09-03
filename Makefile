GO=go
MAIN=gotapiri.go
TARGET=gotapiri

all:
	${GO} build ${MAIN}

clean:
	rm -f ${TARGET} 
