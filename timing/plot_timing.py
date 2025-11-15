import re
from collections import defaultdict
from glob import glob

import numpy as np
from matplotlib import pyplot as plt

DEFAULT_FIGSIZE = np.array(plt.rcParams['figure.figsize'])
"""For scaling the figure contents."""


def stylize_number(n: int):
    ss = ''
    for i, c in enumerate(reversed(str(n))):
        if i > 0 and i % 3 == 0:
            ss += ','
        ss += c
    return ''.join(reversed(ss))


def sort_dict(d: dict):
    return {k: d[k] for k in sorted(d)}


def parse_time_file(f: str) -> list[float]:
    """Parses a single time file."""
    result = []
    for line in open(f):
        if not line.startswith('real '):
            continue
        result.append(float(line[5:]))
    return result


def parse_time_files() -> dict:
    """Parses all the timimg files."""
    splitter = re.compile(r'^(.*)\.(.*)\.(.*)\.txt$')
    d = {'b': defaultdict(dict), 'v': defaultdict(dict)}
    for f in glob('*.txt'):
        (group, tool, n) = splitter.findall(f)[0]
        n = int(n) * 100000
        d[group][n][tool] = parse_time_file(f)
    d = {k: sort_dict(d[k]) for k in d}
    return d


def plot_times(data: dict, fout: str):
    """Generates the bar plot."""
    plt.style.use('ggplot')
    plt.figure(dpi=400, figsize=DEFAULT_FIGSIZE * 0.75 * [2, 1])

    group_names = {
        'b': 'Bacterial Dataset (10GB)',
        'v': 'Viral Dataset (0.5GB)',
    }
    tool_names = {
        'iss': 'InSilicoSeq',
        'izzy': 'Izzy',
    }

    for i, (g, gdata) in enumerate(data.items()):
        plt.subplot(121 + i)

        xx = np.array(list(range(len(gdata))))
        for v in gdata.values():
            labels = sorted(v)
            break
        print('Labels:', labels)

        BAR_WIDTH = 0.2
        for i, lab in enumerate(labels):
            x = xx - BAR_WIDTH * (len(labels) - 1) / 2 + BAR_WIDTH * i
            throughput = [[n / y for y in x[lab]] for n, x in gdata.items()]
            means = [np.mean(x) for x in throughput]
            stds = [np.std(x) for x in throughput]
            plt.bar(
                x,
                means,
                BAR_WIDTH * 0.9,
                yerr=stds,
                capsize=5,
                label=tool_names[lab],
            )
        plt.xticks(xx, [stylize_number(i) for i in gdata])
        plt.legend()
        plt.xlabel('Sample size (reads)')
        plt.ylabel('Reads per second')
        plt.title(group_names[g])

    plt.tight_layout()

    if fout:
        plt.savefig(fout)
        plt.close()
    else:
        plt.show()


def print_averages(data: dict):
    for g, gdata in data.items():
        print(g)
        for n, nv in gdata.items():
            for k, v in nv.items():
                print(f'{n} {k} {np.mean(v):.1f}+-{np.std(v):.1f}')


def print_throughput(data: dict):
    for g, gdata in data.items():
        print(g)
        for n, nv in gdata.items():
            for k, v in nv.items():
                print(f'{n} {k} {n/np.mean(v):.1f}')


times = parse_time_files()
print(times)

plot_times(times, 'time')
print_averages(times)
print_throughput(times)
