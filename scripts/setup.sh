s=$(which sudo)
if [[ $s == "" ]]; then
    apt update
    apt -y install sudo
fi
h=$(history | tail | grep update)
if [[ $h == "" ]]; then
    sudo apt update
fi
sudo apt -y install golang make ncbi-blast+ phylonium primer3 tar wget
