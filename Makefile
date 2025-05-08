all: clean mc-server bin

mc-server:
	go build -o $@ cmd/mc-server/main.go

bin:
	mkdir -p /opt/bin/mc-server/conf
	cp mc-server /opt/bin/mc-server
	cp cmd/mc-server/conf/config.yaml /opt/bin/mc-server/conf

clean:
	rm -f mc-server
	rm -rf /opt/bin/mc-server


