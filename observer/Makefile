SUBDIRS := $(wildcard cmd/*/.)

all:export FP_LOG_CONF=../../configs/config.json

all: build

build: $(SUBDIRS)
$(SUBDIRS):
	$(MAKE) -C $@

.PHONY: all $(SUBDIRS)

clean:
	$(GO) clean
	touch bin/tmp.tmp
	rm bin/*
