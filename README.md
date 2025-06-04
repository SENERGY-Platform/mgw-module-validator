mgw-module-validator
=======

![Image](https://img.shields.io/github/v/tag/SENERGY-Platform/mgw-module-validator?filter=v%2A&label=latest)


## Run from terminal

    mgw-mod-validator [OPTION]... [PATH]

#### Available options:

    -b string   base path
    -d          check dependencies
    -f string   output format [text, json] (default "text")
    -m          validate multiple modules
    -o string   output file path
    -t string   target path
    -v          print version

## Run as docker container

    docker run --rm --mount type=bind,src=[DIR CONTAING ONE OR MANY MODULES],dst=/mnt/data,ro ghcr.io/senergy-platform/mgw-module-validator:[TAG|latest] [OPTION]...

#### Available options:

    -d          check dependencies
    -f string   output format [text, json] (default "text")
    -m          validate multiple modules
    -o string   output file path
    -t string   target path
    -v          print version