package main

import (
	"flag"
	"fmt"
	"github.com/ansilh/mincrd/pkg/ctr"
	"github.com/ansilh/mincrd/pkg/downloader"
	"github.com/ansilh/mincrd/pkg/rootfs"
	"os"
	"path/filepath"
)

func main() {

	//Simple hack to load namespaces prior to starting container
	if os.Args[0] == "start-ctr" {
		ctr.StartContainer(os.Args[1], os.Args[2], os.Args[3])
	} else {
		ctrName := flag.String("name", "", "Container Name")
		flag.Parse()
		if len(*ctrName) == 0 {
			flag.PrintDefaults()
			os.Exit(1)
		}
		//Generate Container ID
		ctrID := rootfs.GenCtrID(*ctrName)

		//Base directory to store all container rootfs
		baseDir := "/var/tmp/mincrd"

		//Container Root filesystem
		ctrRoot := filepath.Join(baseDir, ctrID)

		//Busybox URL - Keeping only Busybox for simplicity
		imageURL := "https://busybox.net/downloads/binaries/1.28.1-defconfig-multiarch/busybox-x86_64"

		//Create container root filesystem
		err := rootfs.CreateRootDir(ctrRoot)
		if err != nil {
			fmt.Println("ERROR", err)
		}

		//Download BusyBox binary and extract it
		_, err = downloader.GetFile("busybox", ctrRoot, imageURL)
		if err != nil {
			fmt.Println("ERROR", err)
		}

		//Create all needed directories and symlinks in container rootfs
		err = ctr.CreateRootfs(ctrRoot)
		if err != nil {
			fmt.Println("ERROR", err)
		}
		//Set namespaces and pass arguments to exec'd child
		ctr.SetNameSpaces(*ctrName, ctrRoot, "sh")
	}
}
