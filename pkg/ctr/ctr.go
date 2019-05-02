/*
Package ctr ...
Package to create a busybox container
*/
package ctr

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

//CreateRootfs ...
//Create BusyBox Pseudo root file system
func CreateRootfs(busyboxPath string) error {
	oldDir, err := os.Getwd()
	if err != nil {
		return err
	}
	err = os.Chdir(busyboxPath)
	if err != nil {
		return err
	}
	fmt.Println("* Setting up container root filesystem in ", busyboxPath)
	dirs := []struct {
		name       string
		permission os.FileMode
	}{
		{"bin", 0755},
		{"proc", 0755},
		{"sys", 0755},
		{"dev", 0755},
	}
	for _, v := range dirs {
		err = os.Mkdir(v.name, v.permission)
		if err != nil {
			return err
		}
	}

	stdOut, err := exec.Command(filepath.Join(busyboxPath, "busybox"), "--list").Output()
	if err != nil {
		fmt.Println("ERROR ", err)
		return err
	}
	fmt.Println("* Creating symlinks inside container rootfs...")
	for _, command := range strings.Split(string(stdOut), "\n") {
		if len(command) != 0 {
			os.Symlink("/busybox", "bin/"+command)
			if err != nil {
				fmt.Println("ERROR ", err)
				return err
			}
		}
	}
	//To get back to the original directory where we invoked the binary
	err = os.Chdir(oldDir)
	if err != nil {
		return err
	}
	return err
}

//SetNameSpaces ...
//Function to change Arg[0] to "fork" and set namespaces
//Finally call the binary itself
func SetNameSpaces(name, rootfs, program string) {
	fmt.Println("* Setting up namespaces...")
	cmd := &exec.Cmd{
		Path: os.Args[0],
		Args: append([]string{"start-ctr"}, name, rootfs, program),
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUSER |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUTS,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}
	fmt.Println("* Exec-ing runtime with new namespace settings...")
	cmd.Run()
}

//StartContainer ...
//
func StartContainer(name, ctrRoot, program string) {
	fmt.Println("* Container setup...")
	//fmt.Println(name, ctrRoot, program)
	//Set a new hostname for our container
	fmt.Println("* Setting up Hostname as", name)
	if err := syscall.Sethostname([]byte(name)); err != nil {
		fmt.Printf("* Setting Hostname failed\n")
	}
	mountData := []struct {
		source string
		target string
		fstype string
		flags  uintptr
		data   string
	}{
		{
			"proc",
			filepath.Join(ctrRoot, "proc"),
			"proc",
			uintptr(syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV),
			"",
		},
		{
			"sysfs",
			filepath.Join(ctrRoot, "sys"),
			"sysfs",
			uintptr(syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV),
			"",
		},
		{
			"tmpfs",
			filepath.Join(ctrRoot, "dev"),
			"tmpfs",
			uintptr(syscall.MS_NOSUID | syscall.MS_STRICTATIME),
			"",
		},
		{
			filepath.Join(ctrRoot),
			filepath.Join(ctrRoot),
			"bind",
			uintptr(syscall.MS_BIND | syscall.MS_REC),
			"",
		},
	}
	fmt.Println("* Mounting virtual filesystems...")
	for _, v := range mountData {
		err := syscall.Mount(v.source, v.target, v.fstype, v.flags, v.data)
		if err != nil {
			fmt.Printf("ERROR: %v", err)
			return
		}
	}

	fmt.Println("* Pivoting rootfs as contiainer root filesystem...")
	pivTmp := filepath.Join(ctrRoot, ".pivot_root")
	if err := os.Mkdir(pivTmp, 0777); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	if err := syscall.PivotRoot(filepath.Join(ctrRoot), pivTmp); err != nil {
		fmt.Printf("Error pivot: %v\n", err)
		return
	}
	if err := syscall.Chdir("/"); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	pivTmp = filepath.Join("/", ".pivot_root")

	if err := syscall.Unmount(pivTmp, syscall.MNT_DETACH); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	err := os.Remove(pivTmp)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	fmt.Println("* Setting up shell environment...")
	err = os.Setenv("USER", "root")
	if err != nil {
		fmt.Println("* Env set failed", err)
	}
	err = os.Setenv("PS1", "${USER}@`hostname` # ")
	if err != nil {
		fmt.Println("* Env set failed", err)
		return
	}
	fmt.Println("* Starting shell...")
	if err := syscall.Exec("/busybox", strings.Fields(program), os.Environ()); err != nil {
		fmt.Println("* Exec failed")
		return
	}

}
