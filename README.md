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
go get github.com/enzoh/go-bls
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
go get github.com/enzoh/go-bls
go get github.com/klauspost/reedsolomon
```
