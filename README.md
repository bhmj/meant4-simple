# meant4

## What is it?

Sample factorial calculating task. Three versions.

v1 : calculate two factorials simultaniously in two goroutines.  
v2 : calculate N factorials in one loop.  
v3 : calculate N factorials in one loop with parallelism.  

Results on `[200000,300000]` :  
v1 : 10.90 seconds  
v2 : 9.47 seconds  
v3 : 1.38 seconds  

Compiled and checked with go 1.17.

## Building, testing and running

Run `make all` to build all three versions in ./build/.  
Run `make help` to get more info.  

Sample command:  

`time curl --request POST --url http://localhost:8989/factorial --data '{"numbers": [200000, 300000]}' -o /dev/null`  
