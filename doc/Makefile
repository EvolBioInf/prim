SRC1 =   cops.tex fa2prim.tex prim2tab.tex scop.tex tutorial.tex util.tex
SRC2 =   header.tex intro.tex primDoc.tex

date = $(shell git log | grep Date | head -n 1 | sed -r 's/Date: +[A-Z][a-z]+ ([A-Z][a-z]+) ([0-9]+) [^ ]+ ([0-9]+) .+/\2_\1_\3/')
version = $(shell git describe)

prog = prim

all: $(prog)Doc.pdf 
$(prog)Doc.pdf: $(prog)Doc.tex $(SRC1) $(SRC2) version.txt date.txt
	latex $(prog)Doc.tex
	bibtex $(prog)Doc
	latex $(prog)Doc.tex
	latex $(prog)Doc.tex
	dvips $(prog)Doc -o -q
	ps2pdf -dALLOWPSTRANSPARENCY $(prog)Doc.ps
%.tex: ../*/%.org $(prog)Doc.tex
	bash ../scripts/org2nw $< | awk -f ../scripts/preWeave.awk | noweave -n -x | sed 's/_test/\\_test/g' > $@
version.txt:
	echo $(version) | tr '-' ' ' | awk '{printf "%s", $$1; if ($$2) printf "-%s", $$2; printf "\n"}' > version.txt
date.txt:
	echo $(date) | tr '_' ' ' > date.txt
clean:
	rm -f $(prog)Doc.pdf $(prog)Doc.ps $(prog)Doc.dvi $(prog)Doc.aux $(prog)Doc.bbl  $(prog)Doc.blg  $(prog)Doc.dvi  $(prog)Doc.log $(SRC1) version.txt date.txt