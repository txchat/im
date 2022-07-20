#!/bin/bash
default_url="http://172.16.101.107:16060"
default_prefix="debug/pprof"
default_type="profile"

# shellcheck disable=SC2162
read -p "input the request url (default url is '${default_url}'): " url
if [ -z "${url}" ]; then
    url=${default_url}
fi
# shellcheck disable=SC2162
read -p "input the pprof prefix (default prefix is '${default_prefix}'): " prefix
if [ -z "${prefix}" ]; then
    prefix=${default_prefix}
fi

function show() {
    echo "chose the pprof type:"
    echo "  0: break the loop"
    echo "  1: profile"
    echo "  2: heap"
    echo "  3: allocs"
    echo "  4: goroutine"
    echo "  5: mutex"
    echo "  6: block"
    echo "  7: trace"
    echo "  8: threadcreate"
    echo "  9: cmdline"
    echo -n "input the number: "
}

while true; do
    show
    # shellcheck disable=SC2162
    read num
    case $num in
        0)
            break
            ;;
        1)
            type="profile"
            ;;
        2)
            type="heap"
            ;;
        3)
            type="allocs"
            ;;
        4)
            type="goroutine"
            ;;
        5)
            type="mutex"
            ;;
        6)
            type="block"
            ;;
        7)
            type="trace"
            ;;
        8)
            type="threadcreate"
            ;;
        9)
            type="cmdline"
            ;;
        *)
            type=${default_type}
            ;;
    esac
    #    echo 'file://wsl$/Ubuntu-18.04/tmp/' | clip.exe
    #    echo "request: ${url}/${prefix}/${type}"
    go tool pprof "${url}/${prefix}/${type}"
done
