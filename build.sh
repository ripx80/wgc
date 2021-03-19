#!/bin/sh

# compress with upx
#make clean && make && docker run -it --rm -v $(pwd)/bin:/in ripx80/upx -9 -o in/brewman_upx in/brewman

# without upx
make clean && make test && make
exit 0