image:
	docker build --target autobot -t ingoda/observer/autobot .
	docker build --target nats-sniffer -t ingoda/observer/nats-sniffer .
	docker build --target selfcheck-syslog -t ingoda/observer/selfcheck-syslog .

