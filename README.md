# Izzy

A large-scale metagenomic read simulator.

## Error models

Izzy uses the error models from [InSilicoSeq].

[InSilicoSeq]: https://github.com/HadrienG/InSilicoSeq

## How to use

Run `izzy` for the full list of command-line options.

### Basic use

```
izzy -i genomes.fasta -o my_reads -n 1000000 -m basic
```

Creates 1M reads from the genomes in `genomes.fasta` and outputs
two files (paired-end) with the prefix `my_reads`,
using the basic error model and
the default abundance distribution (log-normal).

### Multiple genome files

```
izzy -i "genomes/*.fasta" -o my_reads -n 1000000 -m basic
```

The input path may contain glob characters,
including `*` for anything,
`?` for any single character,
`[xyz]` for any single character in x, y and z,
or `[x-z]` for any single character from x to z.

Quotes may be needed around the glob pattern,
to prevent the shell from expanding it.

### Multiple contigs per genome

```
izzy -i genomes.fasta -o my_reads -n 1000000 -m basic -g "Species_\d+"
```

To consider each group of fasta entries as one species,
a grouping criterion can be provided as a [regular expression].
The regular expression is matched against the fasta entry name,
and entries that share the same matched value are grouped together.

[regular expression]: https://pkg.go.dev/regexp/syntax

In this example, assume the naming scheme is `Species_[number]_contig_[number]`.
The expression `Species_\d+` will match the `Species_[number]` part and will
group together all entries with `Species_1000`,
all entries with `Species_1001`, etc.

Quotes may be needed around the regular expression,
since some special characters may be picked up by the shell and
trigger unwanted behavior.
