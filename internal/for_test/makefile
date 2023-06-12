
MAKEFILE_DIR := $(dir $(lastword $(MAKEFILE_LIST)))

.PHONY: test
test:
	cd $(MAKEFILE_DIR) && ../../easygpt \
		--apikey-file ./apikey.txt \
		--setting ./easygpt.hjson
