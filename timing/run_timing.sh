n=1000000
u=10
infile=~/Downloads/representatives.fa
outdir=~/Downloads
infile=~/Downloads/viral.1.1.genomic.fna

for i in 1 1 1 1 1; do
/usr/bin/time -p -a -o v.iss.$((n/100000)).txt \
  iss generate -p 1 -m NovaSeq -n $n -u $u -z \
  -g $infile -o $outdir/iss
done

for i in 1 1 1 1 1; do
/usr/bin/time -p -a -o v.izzy.$((n/100000)).txt \
  izzy -m novaseq -n $n -u $u \
  -i $infile -o $outdir/izz
done
# -g '^.*_'
