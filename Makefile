packs = util
progs = cops fa2prim prim2tab scop
all:
	test -d bin || mkdir bin
	for pack in $(packs); do \
		make -C $$pack; \
	done
	for prog in $(progs); do \
		make -C $$prog; \
		cp $$prog/$$prog bin; \
	done
test: data
	test -d bin || mkdir bin
	for prog in $(packs) $(progs); do \
		make test -C $$prog; \
	done
data:
	wget https://owncloud.gwdg.de/index.php/s/VtVN3IcZsNSpxei/download
	tar -xvzf download
	rm download
tangle:
	for pack in $(packs) $(progs); do \
		make tangle -C $$pack; \
	done
.PHONY: doc test docker
doc:
	make -C doc
docker:
	make -C docker
clean:
	for pack in $(packs); do \
		make clean -C $$pack; \
	done
	for prog in $(progs) $(packs) doc; do \
		make clean -C $$prog; \
	done
