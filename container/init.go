package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

var old_root = ".pivot_root"

func RunContainerInitProcess() error {
	cmds := readUserCommand()
	if cmds == nil || len(cmds) == 0 {
		return fmt.Errorf("error empty cmds")
	}
	log.Infof("begin setUpMount")
	setUpMount()
	path, err := exec.LookPath(cmds[0])
	if err != nil {
		log.Errorf("path error %v", err)
		return err
	}
	log.Infof("Find path %s", path)
	if err := syscall.Exec(path, cmds[0:], os.Environ()); err != nil {
		log.Errorf(err.Error())
		return err
	}
	return nil
}

func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Errorf("read pipe error")
		return nil
	}
	return strings.Split(string(msg), " ")
}

func setUpMount() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Errorf("PWD err %v", err)
	}
	log.Infof("PWD is %s", pwd)
	syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, "")
	err = pivotRoot(pwd)
	if err != nil {
		log.Errorf("pivotRoot err %v", err)
	}
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")
	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
}

func pivotRoot(root string) error {
	//这个root目录在之前通过cmd.Dir='/path/busybox'设置好了
	// new_root 和put_old 必须不能同时存在当前root 的同一个文件系统中,需要通过--bind重新挂载一下
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("mount --bind PWD PWD error ")
	}
	pivotDir := filepath.Join(root, old_root)
	if err := os.Mkdir(pivotDir, 0777); err != nil {
		return err
	}
	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("privot_root error %v", err)
	}
	log.Infof("now change dir to root")
	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir / %v", err)
	}
	// 更新下文件路径
	pivotDir = filepath.Join("/", old_root)
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("umount pivot_root dir %v", err)
	}
	return os.Remove(pivotDir)
}
