import re

import numpy as np
from matplotlib import pyplot as plt
from myplot import ctx


def parse_time(x: str) -> float:
    m = re.match('^(.*)m(.*)s\\s*$', x)
    if not m:
        raise ValueError(f'bad time: {x!r}')
    return 60 * float(m[1]) + float(m[2])


def parse_time_table(f: str):
    d = {'n': [], 'InSilicoSeq': [], 'Izzy': []}
    for line in open(f):
        parts = line.split(' ')
        d['n'].append(int(parts[0]))
        d['InSilicoSeq'].append(parse_time(parts[1]))
        d['Izzy'].append(parse_time(parts[2]))
    return {k: np.array(v) for k, v in d.items()}


FILE = '/dfs7/whitesonlab/alavon/Code/izzy/src/publication/izzy_times.txt'

table = parse_time_table(FILE)

plt.style.use('bmh')
with ctx('izzy_times', dpi=300):
    for i, k in enumerate(list(table)[1:]):
        data = table[k]
        offset = -0.05 + i * 0.1
        plt.bar(np.arange(len(data)) + offset, table['n'] / data, width=0.1, label=k)
    plt.xticks(np.arange(len(table['n'])), [f'{x:,}' for x in table['n']])
    plt.yscale('log')
    plt.xlabel('Sample size (reads)')
    plt.ylabel('Reads per second')
    plt.legend()
