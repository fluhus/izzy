# Builds release binaries.

set -e

VERSION=0.1
OUTDIR=../release
FLAGS="-ldflags=-s -X main.version=$VERSION"

rm -fr $OUTDIR
mkdir $OUTDIR

# Linux.
go build "$FLAGS" -o $OUTDIR ./izzy
zip -j $OUTDIR/izzy_linux.zip $OUTDIR/izzy
rm $OUTDIR/izzy

# Mac.
GOOS=darwin go build "$FLAGS" -o $OUTDIR ./izzy
zip -j $OUTDIR/izzy_mac.zip $OUTDIR/izzy
rm $OUTDIR/izzy

# Windows.
GOOS=windows go build "$FLAGS" -o $OUTDIR ./izzy
zip -j $OUTDIR/izzy_win.zip $OUTDIR/izzy.exe
rm $OUTDIR/izzy.exe

ls -lh $OUTDIR
