# [`prim`](https://owncloud.gwdg.de/index.php/s/HoGv0q8FXumX6Vl)
## Description
Programs for designing and testing
primers.
## Author
[Bernhard Haubold](http://guanine.evolbio.mpg.de/), `haubold@evolbio.mpg.de`
## Make the Programs
On an Ubuntu system like Ubuntu on
[wsl](https://learn.microsoft.com/en-us/windows/wsl/install) under
MS-Windows or the [Ubuntu Docker
container](https://hub.docker.com/_/ubuntu), you can clone the
repository and change into it.

`git clone https://github.com/evolbioinf/prim`  
`cd prim`

Then install the additional dependencies by running the script
[`setup.sh`](scripts/setup.sh).

`bash scripts/setup.sh`

Make the programs.

`make`

The directory `bin` now contains the binaries.

## License
[GNU General Public License](https://www.gnu.org/licenses/gpl.html)
