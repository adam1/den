#!/bin/bash

set -x
set -e

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR

begin=1
end=10

# ./run-sequence.bash NumMaximalTypes $begin $end
# ./run-sequence.bash NumMaximalTypesV3 $begin $end
# diff NumMaximalTypes.txt NumMaximalTypesV3.txt

# ./run-sequence.bash Width $begin $end
# ./run-sequence.bash WidthV3 $begin $end
# diff Width.txt WidthV3.txt

./run-sequence.bash Density $begin $end
./run-sequence.bash DensityV3 $begin $end
diff Density.txt DensityV3.txt

