# Indy SDK Go Bindings

Go bindings for hyperledger indy sdk
Installation go get github.com/joyride9999/IndySdkGoBindings

## Window usage
- Install go (https://golang.org)
- Install gcc (http://mingw-w64.org/doku.php) and add gcc to path
- Optional - install IDE (https://www.jetbrains.com/go/)
- Install IndySDK for windows (https://hyperledger-indy.readthedocs.io/projects/sdk/en/latest/docs/build-guides/windows-build.html). 
    - For windows the dlls/lib files needs to be copied to the lib folder from this project. Files can be downloaded from here https://repo.sovrin.org/windows/libindy/
## Linux
- Install go
- Install indy-sdk (more exactly libindy). Details can be found here https://github.com/hyperledger/indy-sdk
- Optional - install IDE (https://www.jetbrains.com/go/)

## Testing the package
- For testing the blockchain is needed a blockchain test network . (https://github.com/hyperledger/indy-node). Easiest way is to set up the standard indyNode server using docker container.

## TODOs
- review db storage
- optimization
- testing
- 
## LICENSE
This work is licensed under the terms of the Apache License Version 2.0.  See the LICENSE.txt file in the top-level directory.
Header files from "include/" folder were obtained from https://github.com/hyperledger/indy-sdk/tree/master/libindy/include and were 
modified to fit this project needs! 