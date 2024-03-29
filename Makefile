NAME=    getpass
PREFIX?= /usr/local

all: ${NAME}.cmd

${NAME}.cmd: ./${NAME}/main.go
	go build -o ${NAME}.cmd ${NAME}/main.go

help:
	@echo "The following targets are available:"
	@echo "all        build the executable"
	@echo "clean      remove build files"
	@echo "doc        format man page into .txt"
	@echo "install    install ${NAME} into ${PREFIX}"
	@echo "test       run 'go test'"
	@echo "uninstall  uninstall ${NAME} from ${PREFIX}"

install:
	mkdir -p ${PREFIX}/bin ${PREFIX}/share/man/man1
	install -c -m 0555 ./${NAME}.cmd ${PREFIX}/bin/${NAME}
	install -c -m 0444 doc/${NAME}.1 ${PREFIX}/share/man/man1/${NAME}.1

uninstall:
	rm -f ${PREFIX}/bin/${NAME}
	rm -f ${PREFIX}/share/man/man1/${NAME}.1

test:
	go test

clean:
	rm -f ${NAME}.cmd

doc: doc/${NAME}.1.txt

doc/${NAME}.1.txt: doc/${NAME}.1
	mandoc -c -O width=80 $? | col -b >$@
