#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

sequence_name=$1
begin=$2
end=$3

$DIR/../bin/sequence -b $begin -e $end $sequence_name > $sequence_name.txt 2> $sequence_name.log

