#!/bin/bash
# This downloads and install python3.4 to the directory passed in

VERSION=3.4.9
DEST=$1

if [ "$DEST" = "" ]; then
    echo "Syntax: $0 <install-directory>"
    exit 1
fi

if [ -e "$DEST/bin/python3.4" ]; then
    echo "Python already installed in $DEST - skipping install"
    exit 0
fi

mkdir -p $DEST

cd /tmp
curl https://www.python.org/ftp/python/${VERSION}/Python-${VERSION}.tar.xz --output Python-${VERSION}.tar.xz --silent
tar Jxf Python-${VERSION}.tar.xz
cd Python-${VERSION}
./configure --prefix=$DEST
make
make install
echo "Python $VERSION installed in $DEST"
