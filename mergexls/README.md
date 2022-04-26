# for ubuntu
## step1
sudo apt-get install python3

sudo apt-get install python3-pip

pip3 install pyexcel pyexcel-xls pyexcel-xlsx

## step2
GOOS=linux go build -o _output/mergexls_linux mergexls/mergexls.go

## step 3
./_output/mergexls_linux  --to-merge-xls-dir "/paas_worksheets/202205_A/shanghai" --output-file-name "shanghai_merged.xlsx"  --python-cmd "python3"

# for windows
## step 1
manually install python3 pkg from internet,  for example https://www.python.org/ftp/python/3.10.4/python-3.10.4-amd64.exe

manually download python excel libs for offline

```pip3 download pyexcel pyexcel-xls pyexcel-xlsx```

and open your git-bash,  install these libs

```pip3 install  *.whl```

## step 2
GOOS=windows go build -o _output/mergexls.exe mergexls/mergexls.go

## step 3
_output\mergexls.exe  --to-merge-xls-dir "D:\paas_worksheets\202205_A\shanghai" --output-file-name "shanghai_merged.xlsx"  --python-cmd "python"
