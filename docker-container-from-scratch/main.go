package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"path/filepath"
	"strconv"
)

func main() {

	switch os.Args[1] {
	case "run":
		run()
	case "child":
        child()
	default:
		panic("invalid command")
	}
}


func run() {
	fmt.Printf("Parent running %v as PID %d\n", os.Args[2:], os.Getpid())

	// Connect container standard I/O streams directly to your terminal
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
    }

	must(cmd.Run())
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func child() {
    fmt.Printf("Container running %v as PID %d\n", os.Args[2:], os.Getpid())

    // Because we are already inside the new UTS namespace, this only sets
    // the hostname for our container, completely isolated from the host!
    must(syscall.Sethostname([]byte("container")))

	// Configure cgroup resource limits before entering the jail
    cg()


	// 2. Change the root directory for this process and its children to /container-fs
    must(syscall.Chroot("/container-fs"))
    must(syscall.Chdir("/"))

	// Mount a fresh, private /proc pseudo-filesystem
    must(syscall.Mount("proc", "proc", "proc", 0, ""))

    // Now launch the actual user command (e.g. /bin/bash) inside the configured container
    cmd := exec.Command(os.Args[2], os.Args[3:]...)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    must(cmd.Run())

	must(syscall.Unmount("/proc", 0))
}

func cg() {
    cgroup := "/sys/fs/cgroup/demo-container"

    // 1. Create a cgroup node directory
    must(os.Mkdir(cgroup, 0755))

    // 2. Limit the container process hierarchy to at most 20 concurrent PIDs
    must(os.WriteFile(filepath.Join(cgroup, "pids.max"), []byte("20"), 0700))

    // 3. Attach our current process (os.Getpid()) to this cgroup controller
    must(os.WriteFile(filepath.Join(cgroup, "cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}