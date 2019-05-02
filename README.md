# A mini Container runtime

This project is for educational purpose only.
I've developed this program while learning golang

Intentionally avoided many guardrails to keep the program simple

## Documentation

* Installation 

```shell
$ go get github.com/ansilh/mincrd
$ go build
```

* Usage 
```shell
$ mkdir /var/tmp/mincrd
$./mincrd -name <name of the container>
```
* Demo

```bash
$ ./mincrd -name container-1
* Downloading busybox...Done , size -  1001112 Bytes
* Setting up container root filesystem in  /var/tmp/mincrd/container-1-7813417024579691178
* Creating symlinks inside container rootfs...
* Setting up namespaces...
* Exec-ing runtime with new namespace settings...
* Container setup...
* Setting up Hostname as container-1
* Mounting virtual filesystems...
* Pivoting rootfs as contiainer root filesystem...
* Setting up shell environment...
* Starting shell...
root@container-1 # 
root@container-1 # df
Filesystem           1K-blocks      Used Available Use% Mounted on
/dev/mapper/ubuntu--vg-root
                     101694448   8928648  87556920   9% /
tmpfs                  4084660         0   4084660   0% /dev
root@container-1 # ip a
1: lo: <LOOPBACK> mtu 65536 qdisc noop state DOWN qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
root@container-1 # exit
```

* Packages

    pkg/downloader
    * Package to download a given file from URL and return its size after download

    pkg/rootfs
    * Package to support contianer root filesystem creation

    pkg/ctr 
    * Package to create a busybox container

