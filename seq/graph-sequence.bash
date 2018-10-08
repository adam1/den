#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

sequence_name=$1

data_file=$sequence_name.txt
graph_file=$sequence_name.png

echo "set terminal png; set output \"$graph_file\"; set grid xtics lt 0; set grid ytics lt 0; plot \"$data_file\" using 1:2 with lines title '$sequence_name';" | gnuplot

echo wrote $graph_file
