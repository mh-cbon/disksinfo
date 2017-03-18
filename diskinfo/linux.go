package diskinfo

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// LinuxLoader can load disk information for a linux system using df / fidsk.
type LinuxLoader struct {
}

// Load returns the list of partition found and their properties.
func (l *LinuxLoader) Load() ([]*Properties, error) {
	//-
	var ret PropertiesList

	if temp, err := runDf(); err != nil {
		return ret, err
	} else {
		ret = ret.Append(temp)
	}
	//-
	if temp, err := runLsLabel(); err != nil {
		return ret, err
	} else {
		ret = ret.Append(temp)
	}
	//-
	if temp, err := runLsUsb(); err != nil {
		return ret, err
	} else {
		ret = ret.Merge(temp, "IsRemovable")
	}
	//-
	if temp, err := runMount(); err != nil {
		return ret, err
	} else {
		ret = ret.Merge(temp, "Label")
	}
	//-
	return ret, nil
}

func runLsLabel() ([]*Properties, error) {
	var ret []*Properties

	path := "/dev/disk/by-label/"
	disks, err := runLs(path)
	if err != nil {
		return disks, err
	}

	for _, disk := range disks {
		ret = append(ret, disk)
	}

	return ret, nil
}

func runLsUsb() ([]*Properties, error) {
	var ret []*Properties

	path := "/dev/disk/by-id/"
	disks, err := runLs(path)
	if err != nil {
		return disks, err
	}

	for _, disk := range disks {
		if strings.HasPrefix(disk.Label, "usb-") {
			disk.IsRemovable = true
			ret = append(ret, disk)
		}
	}

	return ret, nil
}

// LsReader ...
type LsReader struct {
	r    io.Reader
	line string
}

// NewLsReader ...
func NewLsReader(r io.Reader) *LsReader {
	return &LsReader{r: r}
}

// Read ...
func (l *LsReader) Read(path string) ([]*Properties, error) {

	/*
			   total 0
		     lrwxrwxrwx 1 root root  9 27 févr. 13:21 /dev/disk/by-id/usb-TOSHIBA_External_USB_3.0_20141114026944-0:0 -> ../../sdb
		     lrwxrwxrwx 1 root root 10 27 févr. 13:21 /dev/disk/by-id/usb-TOSHIBA_External_USB_3.0_20141114026944-0:0-part1 -> ../../sdb1
	*/

	var ret []*Properties
	i := 0

	b := NewLineReader(l.r)
	var err error
	for {
		line, err2 := b.ReadLine()
		err = err2

		if i > 0 && line != "" {
			props := strings.Split(line, " ")
			s := []string{}
			for _, p := range props {
				p = strings.TrimSpace(p)
				if p != "" {
					s = append(s, p)
				}
			}
			p := NewProperties()
			name := s[8]
			if strings.HasPrefix(name, "'") {
				name = name[1:]
			}
			if strings.HasSuffix(name, "'") {
				name = name[:len(name)-1]
			}
			p.Label = strings.Replace(name, "\\x20", " ", -1)
			abs, err2 := filepath.Abs(filepath.Join(path, s[10]))
			if err2 != nil {
				panic(err2)
			}
			p.Path = abs
			ret = append(ret, p)
		}

		if err != nil {
			break
		}
		i++
	}

	if err == io.EOF {
		err = nil
	}

	return ret, err
}

func runLs(path string) ([]*Properties, error) {
	var ret []*Properties
	_, err := os.Stat(path)
	// when there is not any usb disk connected on the computer
	// skip err and return to avoid errors about
	// missing directory.
	if os.IsNotExist(err) {
		return ret, nil
	}
	cmd := exec.Command("ls", "-l", path)
	cmd.Stderr = os.Stderr

	sink, err := cmd.StdoutPipe()
	if err != nil {
		return ret, err
	}

	if err2 := cmd.Start(); err2 != nil {
		return ret, err2
	}

	ret, err = NewLsReader(sink).Read(path)
	if err != nil {
		return ret, err
	}

	if err2 := cmd.Wait(); err2 != nil {
		return ret, err2
	}

	return ret, err
}

func runDf() ([]*Properties, error) {
	var ret []*Properties
	cmd := exec.Command("df", "-h")
	cmd.Stderr = os.Stderr

	sink, err := cmd.StdoutPipe()
	if err != nil {
		return ret, err
	}

	if err2 := cmd.Start(); err2 != nil {
		return ret, err2
	}

	ret, err = NewDfReader(sink).Read()
	if err != nil {
		return ret, err
	}

	if err2 := cmd.Wait(); err2 != nil {
		return ret, err2
	}

	return ret, err
}

// DfReader ...
type DfReader struct {
	r    io.Reader
	line string
}

// NewDfReader ...
func NewDfReader(r io.Reader) *DfReader {
	return &DfReader{r: r}
}

// Read ...
func (l *DfReader) Read() ([]*Properties, error) {

	/*
	   Sys. de fichiers        Taille Utilisé Dispo Uti% Monté sur
	   devtmpfs                  1,9G       0  1,9G   0% /dev
	   tmpfs                     1,9G     54M  1,9G   3% /dev/shm
	*/

	var ret []*Properties
	i := 0

	b := NewLineReader(l.r)
	var err error
	for {
		line, err2 := b.ReadLine()
		err = err2

		if i > 0 && line != "" {
			props := strings.Split(line, " ")
			s := []string{}
			for _, p := range props {
				p = strings.TrimSpace(p)
				if p != "" {
					s = append(s, p)
				}
			}
			p := NewProperties()
			p.Size = s[1]
			p.SpaceLeft = s[3]
			p.Path = s[0]
			p.MountPath = s[5]
			ret = append(ret, p)
		}

		if err != nil {
			break
		}
		i++
	}

	if err == io.EOF {
		err = nil
	}
	return ret, err
}

func runMount() ([]*Properties, error) {
	var ret []*Properties
	cmd := exec.Command("mount", "-l")
	cmd.Stderr = os.Stderr

	sink, err := cmd.StdoutPipe()
	if err != nil {
		return ret, err
	}

	if err2 := cmd.Start(); err2 != nil {
		return ret, err2
	}

	ret, err = NewMountReader(sink).Read()
	if err != nil {
		return ret, err
	}

	if err2 := cmd.Wait(); err2 != nil {
		return ret, err2
	}

	return ret, err
}

// MountReader ...
type MountReader struct {
	r    io.Reader
	line string
}

// NewMountReader ...
func NewMountReader(r io.Reader) *MountReader {
	return &MountReader{r: r}
}

var lineR = regexp.MustCompile(`\s*([^\s]+)\s+on\s+([^\s]+)\s+type\s+([^\s]+)\s+\(([^)]+)\)(\s+\[[^]]+\])?`)

// Read ...
func (l *MountReader) Read() ([]*Properties, error) {

	/*
		sysfs on /sys type sysfs (rw,nosuid,nodev,noexec,relatime)
		proc on /proc type proc (rw,nosuid,nodev,noexec,relatime)
		devtmpfs on /dev type devtmpfs (rw,nosuid,size=1966944k,nr_inodes=491736,mode=755)
		securityfs on /sys/kernel/security type securityfs (rw,nosuid,nodev,noexec,relatime)
		tmpfs on /dev/shm type tmpfs (rw,nosuid,nodev)
		devpts on /dev/pts type devpts (rw,nosuid,noexec,relatime,gid=5,mode=620,ptmxmode=000)
		tmpfs on /run type tmpfs (rw,nosuid,nodev,mode=755)
		tmpfs on /sys/fs/cgroup type tmpfs (ro,nosuid,nodev,noexec,mode=755)
		cgroup on /sys/fs/cgroup/systemd type cgroup (rw,nosuid,nodev,noexec,relatime,xattr,release_agent=/usr/lib/systemd/systemd-cgroups-agent,name=systemd)
		pstore on /sys/fs/pstore type pstore (rw,nosuid,nodev,noexec,relatime)
		cgroup on /sys/fs/cgroup/net_cls,net_prio type cgroup (rw,nosuid,nodev,noexec,relatime,net_cls,net_prio)
		cgroup on /sys/fs/cgroup/freezer type cgroup (rw,nosuid,nodev,noexec,relatime,freezer)
		cgroup on /sys/fs/cgroup/hugetlb type cgroup (rw,nosuid,nodev,noexec,relatime,hugetlb)
		cgroup on /sys/fs/cgroup/devices type cgroup (rw,nosuid,nodev,noexec,relatime,devices)
		cgroup on /sys/fs/cgroup/cpu,cpuacct type cgroup (rw,nosuid,nodev,noexec,relatime,cpu,cpuacct)
		cgroup on /sys/fs/cgroup/cpuset type cgroup (rw,nosuid,nodev,noexec,relatime,cpuset)
		cgroup on /sys/fs/cgroup/perf_event type cgroup (rw,nosuid,nodev,noexec,relatime,perf_event)
		cgroup on /sys/fs/cgroup/memory type cgroup (rw,nosuid,nodev,noexec,relatime,memory)
		cgroup on /sys/fs/cgroup/blkio type cgroup (rw,nosuid,nodev,noexec,relatime,blkio)
		cgroup on /sys/fs/cgroup/pids type cgroup (rw,nosuid,nodev,noexec,relatime,pids)
		configfs on /sys/kernel/config type configfs (rw,relatime)
		/dev/mapper/fedora-root on / type ext4 (rw,relatime,data=ordered)
		systemd-1 on /proc/sys/fs/binfmt_misc type autofs (rw,relatime,fd=42,pgrp=1,timeout=0,minproto=5,maxproto=5,direct,pipe_ino=14678)
		mqueue on /dev/mqueue type mqueue (rw,relatime)
		hugetlbfs on /dev/hugepages type hugetlbfs (rw,relatime)
		debugfs on /sys/kernel/debug type debugfs (rw,relatime)
		nfsd on /proc/fs/nfsd type nfsd (rw,relatime)
		tmpfs on /tmp type tmpfs (rw,nosuid,nodev)
		/dev/sda5 on /home type ext4 (rw,relatime,data=ordered)
		sunrpc on /var/lib/nfs/rpc_pipefs type rpc_pipefs (rw,relatime)
		binfmt_misc on /proc/sys/fs/binfmt_misc type binfmt_misc (rw,relatime)
		tmpfs on /run/user/1001 type tmpfs (rw,nosuid,nodev,relatime,size=396012k,mode=700,uid=1001,gid=1001)
		fusectl on /sys/fs/fuse/connections type fusectl (rw,relatime)
		gvfsd-fuse on /run/user/1001/gvfs type fuse.gvfsd-fuse (rw,nosuid,nodev,relatime,user_id=1001,group_id=1001)
		/dev/sdb1 on /run/media/mh-cbon/whatever type fuseblk (rw,nosuid,nodev,relatime,user_id=0,group_id=0,default_permissions,allow_other,blksize=4096,uhelper=udisks2) [whatever]
	*/

	var ret []*Properties
	i := 0

	b := NewLineReader(l.r)
	var err error
	for {
		line, err2 := b.ReadLine()
		err = err2

		if line != "" {
			x := lineR.FindAllStringSubmatch(line, -1)
			if len(x) > 0 {
				s := x[0][1:]
				if s[0][:1] == "/" {
					p := NewProperties()
					p.MountPath = s[1]
					p.Path = s[0]
					if len(s) > 3 && len(s[4]) > 0 {
						p.Label = strings.TrimSpace(s[4])
						p.Label = p.Label[1 : len(p.Label)-1]
					}
					ret = append(ret, p)
				}
			}
		}

		if err != nil {
			break
		}
		i++
	}

	if err == io.EOF {
		err = nil
	}
	return ret, err
}
