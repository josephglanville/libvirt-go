package libvirt

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"unsafe"
)

/*
#cgo LDFLAGS: -lvirt -ldl
#include <libvirt/libvirt.h>
#include <libvirt/virterror.h>
#include <stdlib.h>
*/
import "C"

type VirConnection struct {
	ptr C.virConnectPtr
}

func NewVirConnection(uri string) (VirConnection, error) {
	cUri := C.CString(uri)
	defer C.free(unsafe.Pointer(cUri))
	ptr := C.virConnectOpen(cUri)
	if ptr == nil {
		return VirConnection{}, errors.New(GetLastError())
	}
	obj := VirConnection{ptr: ptr}
	return obj, nil
}

func GetLastError() string {
	err := C.virGetLastError()
	errMsg := fmt.Sprintf("[Code-%d] [Domain-%d] %s",
		err.code, err.domain, C.GoString(err.message))
	C.virResetError(err)
	return errMsg
}

func (c *VirConnection) CloseConnection() error {
	result := int(C.virConnectClose(c.ptr))
	if result == -1 {
		return errors.New(GetLastError())
	}
	return nil
}

func (c *VirConnection) GetCapabilities() (string, error) {
	str := C.virConnectGetCapabilities(c.ptr)
	if str == nil {
		return "", errors.New(GetLastError())
	}
	capabilities := C.GoString(str)
	C.free(unsafe.Pointer(str))
	return capabilities, nil
}

func (c *VirConnection) GetNodeInfo() (VirNodeInfo, error) {
	ni := VirNodeInfo{}
	var ptr C.virNodeInfo
	result := C.virNodeGetInfo(c.ptr, (*C.virNodeInfo)(unsafe.Pointer(&ptr)))
	if result == -1 {
		return ni, errors.New(GetLastError())
	}
	ni.ptr = ptr
	return ni, nil
}

func (c *VirConnection) GetHostname() (string, error) {
	str := C.virConnectGetHostname(c.ptr)
	if str == nil {
		return "", errors.New(GetLastError())
	}
	hostname := C.GoString(str)
	C.free(unsafe.Pointer(str))
	return hostname, nil
}

func (c *VirConnection) ListDefinedDomains() ([]string, error) {
	var names [1024](*C.char)
	namesPtr := unsafe.Pointer(&names)
	numDomains := C.virConnectListDefinedDomains(
		c.ptr,
		(**C.char)(namesPtr),
		1024)
	if numDomains == -1 {
		return nil, errors.New(GetLastError())
	}
	goNames := make([]string, numDomains)
	for k := 0; k < int(numDomains); k++ {
		goNames[k] = C.GoString(names[k])
		C.free(unsafe.Pointer(names[k]))
	}
	return goNames, nil
}

func (c *VirConnection) ListDomains() ([]uint32, error) {
	domainIds := make([]int, 1024)
	domainIdsPtr := unsafe.Pointer(&domainIds)
	numDomains := C.virConnectListDomains(c.ptr, (*C.int)(domainIdsPtr), 1024)
	if numDomains == -1 {
		return nil, errors.New(GetLastError())
	}

	domains := make([]uint32, numDomains)

	gBytes := C.GoBytes(domainIdsPtr, C.int(numDomains*32))
	buf := bytes.NewBuffer(gBytes)
	for k := 0; k < int(numDomains); k++ {
		binary.Read(buf, binary.LittleEndian, &domains[k])
	}
	return domains, nil
}

func (c *VirConnection) LookupDomainById(id uint32) (VirDomain, error) {
	ptr := C.virDomainLookupByID(c.ptr, C.int(id))
	if ptr == nil {
		return VirDomain{}, errors.New(GetLastError())
	}
	return VirDomain{ptr: ptr}, nil
}

func (c *VirConnection) LookupDomainByName(id string) (VirDomain, error) {
	cName := C.CString(id)
	defer C.free(unsafe.Pointer(cName))
	ptr := C.virDomainLookupByName(c.ptr, cName)
	if ptr == nil {
		return VirDomain{}, errors.New(GetLastError())
	}
	return VirDomain{ptr: ptr}, nil
}
