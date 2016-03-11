BEATNAME=twitterbeat
BEAT_DIR=github.com/buehler
SYSTEM_TESTS=false
TEST_ENVIRONMENT=false
ES_BEATS=./vendor/github.com/elastic/beats
GOPACKAGES=$(shell glide novendor)
PREFIX?=.

# Path to the libbeat Makefile
-include $(ES_BEATS)/libbeat/scripts/Makefile

.PHONY: init
init:
	glide update  --no-recursive
	make update
	git init
	git add .

.PHONY: install-cfg
install-cfg:
	mkdir -p $(PREFIX)
	cp etc/twitterbeat.template.json     $(PREFIX)/twitterbeat.template.json
	cp etc/twitterbeat.yml               $(PREFIX)/twitterbeat.yml
	cp etc/twitterbeat.yml               $(PREFIX)/twitterbeat-linux.yml
	cp etc/twitterbeat.yml               $(PREFIX)/twitterbeat-binary.yml
	cp etc/twitterbeat.yml               $(PREFIX)/twitterbeat-darwin.yml
	cp etc/twitterbeat.yml               $(PREFIX)/twitterbeat-win.yml

.PHONY: update-deps
update-deps:
	glide update  --no-recursive

.PHONY: update
update:
	echo "Update not supported due to custom template"