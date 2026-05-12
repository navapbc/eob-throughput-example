# testing fhir throughput

This small Golang application demonstrates the use of `jq` for processing FHIR data.

## running

Assuming you have a golang environment setup:

```
make generate
./timing.bash
```

## tooling

This application uses a pure golang implementation of `jq`, a command-line tool/language for processing JSON documents. It allows for concise expression of queries/filters against JSON documents, and provides fast, serial processing (meaning the JSON is processed in-array without conversion) of JSON.

For example:

```
del((.entry // [])[] | .resource | (.item // [])[] | (.adjudication // [])[] | .amount)
```

The above expression takes an EoB JSON document and retrieves all of the `entry` objects; if there are none, it optionally `//` continues using the empty list. (This makes sure we don't have errors in the rest of the expression if there is not an `entry` key in the EoB.) Then, the `resource` is extracted from the `entry` (it must be present), and then the `items` are extracted as a list (or we return the empty list). From each of those, we extract the `adjudication` field, and the from that structure, we delete the `amount` object.

When executed, this modifies the JSON byte array representing the ExplanationOfBenefit. It does not transform the JSON into a data structure; instead, it manipulates the byte array directly, removing the entire `amount` object from the JSON tree.

## output

This will perform a series of tests. It aims to answer the following questions regarding processing of FHIR data serially:

1. What is an optimal bundle size?
2. What is an optimal level of concurrency?

When I ask "what is an optimal bundle size?," I'm asking "how many EoBs should be retrieved from the database at once?" Because each EoB is around 500K, this means that amount of RAM needed scales with bundle size... ultimately making horizontal scaling more expensive. Bundle sizes of 10, 20, 50, 100, 200, and 500 are tested.

When I ask "what is an optimal level of concurrency?," I'm asking how many parallel processes (or "go funcs") should I run at one time. I run 1,2, 4, and 8 processing threads concurrently, and for my test environment (a Mac M4 Pro). It has 10 performance cores, but my machine might be doing a lot of other things, so I find an optimal balance at around 4 gofuncs.

## results

The ultimate result is that it is possible to process a 500K EoB with a go-native `jq` in around 1.5ms (or thereabouts; results as low as 0.5ms are possible; the absolute timing is less important than to demonstrate that it is possible to manipulate 500K of JSON four different ways in under 2 milliseconds on an ARM CPU).

```
concurrency,rows,bundle,writing,total time (ms),avg time (us),alloc,totalalloc,sys,peak,numgc
1,10000,10,false,14737,1473,13,14766,64,64,1013
1,10000,20,false,15679,1567,14,15160,76,76,756
1,10000,50,false,15491,1549,21,15189,146,146,434
1,10000,100,false,15666,1566,21,15161,266,266,262
1,10000,200,false,15851,1585,21,15201,488,488,164
1,10000,500,false,16483,1648,21,15202,945,945,119
2,10000,10,false,8808,880,33,14775,84,84,636
2,10000,20,false,8935,893,43,15189,142,142,403
2,10000,50,false,8914,891,123,15234,278,278,206
2,10000,100,false,8980,898,234,15259,566,566,112
2,10000,200,false,9131,913,459,15407,981,981,63
2,10000,500,false,9488,948,708,15896,2013,2013,43
4,10000,10,false,5634,563,51,14794,121,121,424
4,10000,20,false,5674,567,93,15217,208,208,243
4,10000,50,false,5554,555,219,15331,480,480,115
4,10000,100,false,6184,618,399,15426,950,950,65
4,10000,200,false,5904,590,745,15696,1704,1704,43
4,10000,500,false,6878,687,2251,16273,4035,4036,30
8,10000,10,false,5561,556,53,14796,126,126,431
8,10000,20,false,5411,541,87,15210,216,216,240
8,10000,50,false,5336,533,240,15352,522,522,113
8,10000,100,false,5424,542,416,15442,1045,1045,64
8,10000,200,false,5690,569,733,15684,1730,1730,43
8,10000,500,false,6266,626,2962,16989,4318,4318,26
1,10000,10,true,22819,2281,21,25688,89,89,1017
1,10000,20,true,24346,2434,32,25642,130,130,707
1,10000,50,true,23779,2377,64,25154,269,269,347
1,10000,100,true,23472,2347,117,24915,492,492,189
1,10000,200,true,24433,2443,224,25054,949,949,111
1,10000,500,true,24585,2458,545,25302,2275,2275,53
2,10000,10,true,18244,1824,32,25216,105,105,973
2,10000,20,true,17827,1782,53,25545,174,174,561
2,10000,50,true,18272,1827,117,25148,368,368,269
2,10000,100,true,18501,1850,224,25062,689,689,152
2,10000,200,true,18815,1881,438,25239,1315,1316,86
2,10000,500,true,19791,1979,1080,25817,2504,2504,44
4,10000,10,true,15756,1575,54,25429,146,146,616
4,10000,20,true,15361,1536,100,25467,245,245,347
4,10000,50,true,16311,1631,235,25146,542,542,154
4,10000,100,true,15864,1586,457,25035,1020,1020,88
4,10000,200,true,16540,1654,907,25232,2006,2006,54
4,10000,500,true,17845,1784,2251,25813,4507,4507,34
8,10000,10,true,13586,1358,99,25234,237,237,345
8,10000,20,true,14131,1413,190,25498,446,446,194
8,10000,50,true,14033,1403,459,25203,1007,1007,89
8,10000,100,true,14308,1430,896,25346,1950,1950,54
8,10000,200,true,15592,1559,1782,25859,3752,3752,37
8,10000,500,true,33286,3328,4390,27959,9236,9236,30
```