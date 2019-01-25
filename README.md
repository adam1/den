# `den`

## Tools for calculating the cyclic density of symmetric groups

### Background

See the [paper draft](https://adam-marks.com/den/denSnv20_amarks.pdf).

### To build

```
% source gopath.bash
% make
```

### Usage

```
% bin/sequence -help
Usage of bin/sequence:
  -b int
    	begin index (default 1)
  -e int
    	end index (default 10)
  -l	list available sequence names
  -prof string
    	enabling profiling: cpu or mem

% bin/sequence -l
Density
DensityV3
DensityDelta
DensitySum
MinCardinalityCentralizerMaximalType
MinTotientLcmMaximalType
NumMaximalTypes
NumMaximalTypesV3
NumTypes
TypeStoreSizeWithParts
TypeStoreSizeWithSlots
TypeStoreSortTime
Width
WidthV3
WidthV3Time
WidthV3SuccessiveRatio
WidthV3RatioToPreviousFactorial
WidthV3RatioToPreviousFactorialTimesSquareRoot

% bin/sequence -b 1 -e 15 Density > data
2018/10/08 07:55:20 n=1 den=1 partitiontime=0 gentime=0 widthtime=0
2018/10/08 07:55:20 n=2 den=0.5 partitiontime=0 gentime=0 widthtime=0
2018/10/08 07:55:20 n=3 den=0.6666666666666666 partitiontime=0 gentime=0 widthtime=0
2018/10/08 07:55:20 n=4 den=0.5416666666666666 partitiontime=0 gentime=0 widthtime=0
2018/10/08 07:55:20 n=5 den=0.25833333333333336 partitiontime=0 gentime=0 widthtime=0
2018/10/08 07:55:20 n=6 den=0.3416666666666667 partitiontime=0 gentime=0 widthtime=0
2018/10/08 07:55:20 n=7 den=0.2571428571428571 partitiontime=0 gentime=0 widthtime=0
2018/10/08 07:55:20 n=8 den=0.2672123015873016 partitiontime=0 gentime=0 widthtime=0
2018/10/08 07:55:20 n=9 den=0.22938161375661376 partitiontime=0 gentime=0 widthtime=0
2018/10/08 07:55:20 n=10 den=0.2173776455026455 partitiontime=0 gentime=0 widthtime=0
2018/10/08 07:55:20 n=11 den=0.17123541967291966 partitiontime=0 gentime=0 widthtime=0
2018/10/08 07:55:20 n=12 den=0.16986361632194966 partitiontime=0 gentime=0 widthtime=0
2018/10/08 07:55:20 n=13 den=0.13624148035606368 partitiontime=0 gentime=0 widthtime=0
2018/10/08 07:55:20 n=14 den=0.13103779001348445 partitiontime=0 gentime=0 widthtime=0
2018/10/08 07:55:20 n=15 den=0.11951714574327421 partitiontime=0 gentime=0 widthtime=0

% cat data
#n Density
1 1
2 0.5
3 0.6666666666666666
4 0.5416666666666666
5 0.25833333333333336
6 0.3416666666666667
7 0.2571428571428571
8 0.2672123015873016
9 0.22938161375661376
10 0.2173776455026455
11 0.17123541967291966
12 0.16986361632194966
13 0.13624148035606368
14 0.13103779001348445
15 0.11951714574327421

```


