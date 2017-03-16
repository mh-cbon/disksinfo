package diskinfo

import (
	"bufio"
	"bytes"
	"testing"
)

type parseTable struct {
	in        string
	expectErr error
	expectOut []*Properties
	path      string // for ls test
}

func TestDfParser(t *testing.T) {

	testsTable := []parseTable{
		parseTable{
			in: `Sys. de fichiers        Taille Utilisé Dispo Uti% Monté sur
devtmpfs                  1,9G       0  1,9G   0% /dev
tmpfs                     1,9G     67M  1,9G   4% /dev/shm
tmpfs                     1,9G    1,6M  1,9G   1% /run
tmpfs                     1,9G       0  1,9G   0% /sys/fs/cgroup
/dev/mapper/fedora-root    32G     26G  4,1G  87% /
`,
			expectErr: nil,
			expectOut: []*Properties{
				&Properties{
					Size:      "1,9G",
					SpaceLeft: "1,9G",
					Path:      "devtmpfs",
					MountPath: "/dev",
				},
				&Properties{
					Size:      "1,9G",
					SpaceLeft: "1,9G",
					Path:      "tmpfs",
					MountPath: "/dev/shm",
				},
				&Properties{
					Size:      "1,9G",
					SpaceLeft: "1,9G",
					Path:      "tmpfs",
					MountPath: "/run",
				},
				&Properties{
					Size:      "1,9G",
					SpaceLeft: "1,9G",
					Path:      "tmpfs",
					MountPath: "/sys/fs/cgroup",
				},
				&Properties{
					Size:      "32G",
					SpaceLeft: "4,1G",
					Path:      "/dev/mapper/fedora-root",
					MountPath: "/",
				},
				&Properties{
					Size:      "32G",
					SpaceLeft: "4,1G",
					Path:      "/dev/mapper/fedora-root",
					MountPath: "/",
				},
			},
		},
	}

	for i, testTable := range testsTable {

		var b bytes.Buffer
		r := NewDfReader(bufio.NewReader(&b))
		b.WriteString(testTable.in)

		res, err := r.Read()
		if err != nil && testTable.expectErr != err {
			t.Fatalf("Test(%v): Unexpected error %v", i, err)
		}

		for _, p := range res {
			found := PropertiesList(testTable.expectOut).FindByPath(p.Path)
			if found == nil {
				t.Errorf("Test(%v): Unexpected property %q not found\n%#v\ntestTable.in=\n%v", i, p.Path, p, testTable.in)
			}
		}

		for _, p := range testTable.expectOut {
			found := PropertiesList(res).FindByPath(p.Path)
			if found == nil {
				t.Errorf("Test(%v): Property %q not found\n%#v\ntestTable.in=\n%v", i, p.Path, p, testTable.in)
			}
		}
	}
}

func TestLsParser(t *testing.T) {

	testsTable := []parseTable{
		parseTable{
			in: `total 0
lrwxrwxrwx 1 root root 10 27 févr. 11:04 Recovery -> ../../sda1
lrwxrwxrwx 1 root root 10 27 févr. 11:04 stockage -> ../../sda6
lrwxrwxrwx 1 root root 10 27 févr. 11:04 'System\x20Reserved' -> ../../sda2
`,
			path:      "/dev/disk/by-label/",
			expectErr: nil,
			expectOut: []*Properties{
				&Properties{
					Label: "Recovery",
					Path:  "/dev/sda1",
				},
				&Properties{
					Label: "stockage",
					Path:  "/dev/sda6",
				},
				&Properties{
					Label: "System Reserved",
					Path:  "/dev/sda2",
				},
			},
		},
	}

	for i, testTable := range testsTable {

		var b bytes.Buffer
		r := NewLsReader(bufio.NewReader(&b))
		b.WriteString(testTable.in)

		res, err := r.Read(testTable.path)
		if err != nil && testTable.expectErr != err {
			t.Fatalf("Test(%v): Unexpected error %v", i, err)
		}

		for _, p := range res {
			found := PropertiesList(testTable.expectOut).FindByPath(p.Path)
			if found == nil {
				t.Errorf("Test(%v): Unexpected property %q not found\n%#v\ntestTable.in=\n%v", i, p.Path, p, testTable.in)
			}
		}

		for _, p := range testTable.expectOut {
			found := PropertiesList(res).FindByPath(p.Path)
			if found == nil {
				t.Errorf("Test(%v): Property %q not found\n%#v\ntestTable.in=\n%v", i, p.Path, p, testTable.in)
			}
		}
	}
}

func TestMountParser(t *testing.T) {

	testsTable := []parseTable{
		parseTable{
			in: `
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
`,
			expectErr: nil,
			expectOut: []*Properties{
				&Properties{
					MountPath: "/",
					Path:      "/dev/mapper/fedora-root",
				},
				&Properties{
					MountPath: "/home",
					Path:      "/dev/sda5",
				},
				&Properties{
					Label:     "Whatever",
					MountPath: "/run/media/mh-cbon/whatever",
					Path:      "/dev/sdb1",
				},
			},
		},
	}

	for i, testTable := range testsTable {

		var b bytes.Buffer
		r := NewMountReader(bufio.NewReader(&b))
		b.WriteString(testTable.in)

		res, err := r.Read()
		if err != nil && testTable.expectErr != err {
			t.Fatalf("Test(%v): Unexpected error %v", i, err)
		}

		for _, p := range res {
			found := PropertiesList(testTable.expectOut).FindByPath(p.Path)
			if found == nil {
				t.Errorf("Test(%v): Unexpected property %q not found\n%#v\ntestTable.in=\n%v", i, p.Path, p, testTable.in)
			}
		}

		for _, p := range testTable.expectOut {
			found := PropertiesList(res).FindByPath(p.Path)
			if found == nil {
				t.Errorf("Test(%v): Property %q not found\n%#v\ntestTable.in=\n%v", i, p.Path, p, testTable.in)
			}
		}
	}
}
