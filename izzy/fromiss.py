"""Converts model files from InSilicoSeq to Izzy format."""

import argparse
import json
from typing import Dict, Iterable

import numpy as np

ntoi = {'A': 0, 'C': 1, 'G': 2, 'T': 3}


def subst_to_arrays(s: Iterable[Dict]):
    arrs = []
    for d in s:
        arr = [None] * 4
        for char, subst in d.items():
            probs = np.array([0.0] * 4)
            for schar, prob in zip(*subst):
                probs[ntoi[schar]] = prob
            probs = probs.cumsum()
            if probs[-1] != 0:
                probs /= probs[-1]
            else:
                # All substitution probabilities are 0
                # -> make it 1 for the original char.
                probs[ntoi[char]:] = 1
            arr[ntoi[char]] = probs.tolist()
        arrs.append(arr)
    return arrs


def indel_to_arrays(dicts: Iterable[Dict]):
    arrs = []
    for d in dicts:
        arr = [0] * 4
        for char, prob in d.items():
            arr[ntoi[char]] = prob
        arrs.append(arr)
    return arrs


def hist_to_arrays(hist: np.ndarray):
    hist = hist.tolist()
    for x in hist:
        for i in range(len(x)):
            x[i] = x[i].tolist()
    return hist


def mean_count_to_array(m: np.ndarray):
    m = m.cumsum().astype(float)
    m /= m[-1]
    return m.tolist()


def model_to_dict(m):
    d = {}
    d['name'] = str(m['model'])
    d['readLen'] = m['read_length'].tolist()
    d['insertLen'] = m['insert_size'].tolist()
    d['meanCountForward'] = mean_count_to_array(m['mean_count_forward'])
    d['meanCountReverse'] = mean_count_to_array(m['mean_count_reverse'])
    d['qualityHistForward'] = hist_to_arrays(m['quality_hist_forward'])
    d['qualityHistReverse'] = hist_to_arrays(m['quality_hist_reverse'])
    d['substChoicesForward'] = subst_to_arrays(m['subst_choices_forward'])
    d['substChoicesReverse'] = subst_to_arrays(m['subst_choices_reverse'])
    d['insForward'] = indel_to_arrays(m['ins_forward'])
    d['insReverse'] = indel_to_arrays(m['ins_reverse'])
    d['delForward'] = indel_to_arrays(m['del_forward'])
    d['delReverse'] = indel_to_arrays(m['del_reverse'])

    return d


argp = argparse.ArgumentParser()
argp.add_argument("-i", help="Input file", type=str, required=True)
argp.add_argument("-o", help="Output file", type=str, required=True)
args = argp.parse_args()

model = np.load(args.i, allow_pickle=True)
model = model_to_dict(model)
json.dump(model, open(args.o, 'wt'))
