exe = tut
all : $(exe).sh

$(exe).sh: tutorial.org
	bash ../scripts/org2nw tutorial.org | notangle -R$(exe).sh > $(exe).sh
test: $(exe).sh
	bash $(exe).sh
clean:
	rm -f $(exe).sh cops.out prim.fasta prim.in prim.out scop.out

