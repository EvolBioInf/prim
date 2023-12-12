# `prim`: Designing and Testing PCR Primers
## Description
Programs for designing and testing
primers. Primer design is built on the program
[`primer3`](https://primer3.org/), primer testing on
[Blast](https://blast.ncbi.nlm.nih.gov/).
## Author
[Bernhard Haubold](http://guanine.evolbio.mpg.de/), `haubold@evolbio.mpg.de`
## Make the Programs
Make sure you've installed the packages `git`, `golang`,
`make`, `ncbi-blast+`, `noweb`, `phylonium`, and `primer3`. Then make the programs.  
  `$ make`  
  The directory `bin` now contains the binaries, scripts are in
  `scripts`.
## Make the Documentation
Make sure you've installed the packages `git`, `make`, `noweb`, `texlive-science`,
`texlive-pstricks`, `texlive-latex-extra`,
and `texlive-fonts-extra`. Then make the documentation.  
  `$ make doc`  
  The documentation is now in `doc/primDoc.pdf`.
## License
[GNU General Public License](https://www.gnu.org/licenses/gpl.html)
