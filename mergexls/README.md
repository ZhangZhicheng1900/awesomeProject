# for ubuntu
## step1
sudo apt-get install python3

sudo apt-get install python3-pip

pip3 install pyexcel pyexcel-xls pyexcel-xlsx

## step2
GOOS=linux go build -o _output/mergexls_linux mergexls/mergexls.go
GOOS=windows go build -o _output/mergexls.exe mergexls/mergexls.go
