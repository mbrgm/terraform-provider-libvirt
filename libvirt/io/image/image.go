package image

import (
	"bufio"
	xioutil "github.com/dmacvicar/terraform-provider-libvirt/libvirt/ioutil"
	"github.com/libvirt/libvirt-go-xml"
	"io"
)

type Format int

const (
	QCOW2 Format = iota
	Raw
)

type Source int

const (
	File = iota
	Vagrant
)

type sized interface {
	Size() (int64, error)
}

const (
	qcow2Magic = "QFI\xfb\x00\x00\x00\x03"
)

type Image struct {
	io.Reader
	io.Closer
	sized
	Format Format
}

func NewImageFromSource(src string) (*Image, error) {
	// network transparent reader
	r, err := xioutil.NewURLReader(src)
	if err != nil {
		return nil, err
	}

	// compression
	a, err := xioutil.NewAnyReader(r)
	if err != nil {
		return nil, err
	}

	// figure out format
	format := Raw
	buf := bufio.NewReader(a)
	b, err := buf.Peek(len(qcow2Magic))
	if err != nil {
		return nil, err
	}
	if string(b) == qcow2Magic {
		format = QCOW2
	}
	return &Image{buf, a, a, format}, nil
}

func Import(src string, vol libvirtxml.StorageVolume) error {
	return nil
}