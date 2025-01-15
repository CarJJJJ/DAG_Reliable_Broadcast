# DAG_Reliable_Broadcast


While we are using SigBRB, we need to install pbc library.

## MACOS install

First we install gmp and pbc library.

```bash
brew install gmp
```

install_pbc.sh is a script to install pbc library.
```bash
wget https://crypto.stanford.edu/pbc/files/pbc-0.5.14.tar.gz
if [[ $(shasum -a 256 pbc-0.5.14.tar.gz) = \
    772527404117587560080241cedaf441e5cac3269009cdde4c588a1dce4c23d2* ]]
then
    tar xf pbc-0.5.14.tar.gz
    pushd pbc-0.5.14
    # 添加 GMP 库路径
    export LDFLAGS="-L/opt/homebrew/lib"
    export CPPFLAGS="-I/opt/homebrew/include"
    # 使用 bash 运行 configure
    bash ./configure
    make
    sudo make install
    popd
else
    echo 'Cannot download PBC library from crypto.stanford.edu.'
    exit 1
fi
```

```bash
./install_pbc.sh
```

and then

```bash
go get github.com/CarJJJJ/go-bls
go get github.com/klauspost/reedsolomon
```

## Linux install

```bash
sudo yum install gmp-devel
```

```bash
wget https://crypto.stanford.edu/pbc/files/pbc-0.5.14.tar.gz
if [[ $(sha256sum pbc-0.5.14.tar.gz) = \
	772527404117587560080241cedaf441e5cac3269009cdde4c588a1dce4c23d2* ]]
then
	tar xf pbc-0.5.14.tar.gz
	pushd pbc-0.5.14
	sh configure
	make
	sudo make install
	popd
else
	echo 'Cannot download PBC library from crypto.stanford.edu.'
	exit 1
fi
```

```bash
./install_pbc.sh
```

and then

```bash
go get github.com/CarJJJJ/go-bls
go get github.com/klauspost/reedsolomon
```


if your device can not find libpbc.so.1
```bash
# creat config
echo "usr/local/lib" | sudo tee /etc/ld.so.conf.d/usr-local-lib.conf

# update config
sudo ldconfig
```

# How to read the DAG_Graph in Json

you can see the code in dag_broadcast/dag_server_process_response.go

```go
writeResponseToJSON(node.Node.Id, responseMessagehash, msg)
```

and you can see the DAG_Graph json file in DAG_Reliable_Broadcast/DAG_Graph

and then run the bash

```bash
go test -v dag_test.go
```

```bash
dot -Tpng DAG_Graph.dot -o DAG_Graph.png
```