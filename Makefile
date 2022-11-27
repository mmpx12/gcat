build:
	go build -ldflags="-w -s" 

install:
	mv gcat /usr/bin/gcat

termux-install:
	mv gcat /data/data/com.termux/files/usr/bin/gcat

all: build install

termux-all: build termux-install

clean:
	rm -f gcat /usr/bin/gcat

termux-clean:
	rm -f gcat /data/data/com.termux/files/usr/bin/gcat
