package libvirt

/*
#cgo LDFLAGS: -lvirt -ldl
#include <libvirt/libvirt.h>
#include <libvirt/virterror.h>
#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"unsafe"
)

type VirDomainSnapshot struct {
	ptr C.virDomainSnapshotPtr
}

func (d *VirDomain) CreateSnapshotXML(xml string, flags uint32) (VirDomainSnapshot, error) {
	cXml := C.CString(xml)
	defer C.free(unsafe.Pointer(cXml))
	result := C.virDomainSnapshotCreateXML(d.ptr, cXml, C.uint(flags))
	if result == nil {
		return VirDomainSnapshot{}, errors.New(GetLastError())
	}
	return VirDomainSnapshot{ptr: result}, nil
}

func (d *VirDomain) Save(destFile string) error {
	cPath := C.CString(destFile)
	defer C.free(unsafe.Pointer(cPath))
	result := C.virDomainSave(d.ptr, cPath)
	if result == -1 {
		return errors.New(GetLastError())
	}
	return nil
}

func (d *VirDomain) SaveFlags(destFile string, destXml string, flags uint32) error {
	cDestFile := C.CString(destFile)
	cDestXml := C.CString(destXml)
	defer C.free(unsafe.Pointer(cDestXml))
	defer C.free(unsafe.Pointer(cDestFile))
	result := C.virDomainSaveFlags(d.ptr, cDestFile, cDestXml, C.uint(flags))
	if result == -1 {
		return errors.New(GetLastError())
	}
	return nil
}
