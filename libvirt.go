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
	connection _Ctype_virConnectPtr
}

type VirDomain struct {
	domain _Ctype_virDomainPtr
}

func NewVirConnection(uri string) (VirConnection, error) {
	cUri := C.CString(uri)
	defer C.free(unsafe.Pointer(cUri))
	ptr := C.virConnectOpen(cUri)
	if ptr == nil {
		return VirConnection{}, errors.New(GetLastError())
	}
	obj := VirConnection{connection: ptr}
	return obj, nil
}

func GetLastError() string {
	err := C.virGetLastError()
	errMsg := fmt.Sprintf("[Code-%d] [Domain-%d] %s",
		err.code, err.domain, C.GoString(err.message))
	C.virResetError(err)
	return errMsg
}

func (c *VirConnection) ListDomains() ([]uint32, error) {
	domainIds := make([]int, 1024)
	domainIdsPtr := unsafe.Pointer(&domainIds)
	numDomains := C.virConnectListDomains(c.connection, (*C.int)(domainIdsPtr), 1024)
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
	ptr := C.virDomainLookupByID(c.connection, C.int(id))
	if ptr == nil {
		return VirDomain{}, errors.New(GetLastError())
	}
	return VirDomain{domain: ptr}, nil
}

func (d *VirDomain) GetName() (string, error) {
	name := C.virDomainGetName(d.domain)
	if name == nil {
		return "", errors.New(GetLastError())
	}
	return C.GoString(name), nil
}
