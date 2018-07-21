# go-vmguestlib

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/xlucas/go-vmguestlib/vmguestlib)
[![GitHub version](https://img.shields.io/github/release/xlucas/go-vmguestlib.svg)](https://github.com/xlucas/go-vmguestlib/releases/tag/v2.0.0)

VMware® vSphere Guest API for Go.

## Disclaimer
This is not an official VMware® product.

## Who is using it ?
[vmgstat](https://github.com/xlucas/vmgstat) : a tool to retrieve CPU & Memory virtualization statistics for guest VMs running on vSphere.

## Requirements
In order to use this library, to build it or to run the tests, you will need to :

* Install the VMware® tools.
* Make sure Guest API runtime components are enabled on vSphere. If needed, add in your virtual machines configurations (.vmx) :

```
isolation.tools.guestlibGetInfo.disable = "FALSE"
```

* Then, depending on your OS, run either :

On Debian/Ubuntu :
```bash
echo /usr/lib/vmware-tools/lib/libvmtools.so >>/etc/ld.so.conf.d/vmware-tools-libraries.conf
```

Or on RedHat/CentOS :
```bash
echo /usr/lib/vmware-tools/lib/libvmtools.so >>/etc/ld.so.conf.d/vmware-tools-guestlib.conf
```

Eventually, independently of the OS your guest VM is running on, update the dynamic loader dependencies with :
```bash
ldconfig
```

## Contributing
Feel free to contribute if you think something is missing. The makefile contains some targets that are used to generate accessors code and unit tests. It uses template files and two input lists that describe the mapping between the native and Go APIs as well as provide documentation. If you need to rebuild the accessor codes and corresponding tests, then run `make gen-accessor`.
