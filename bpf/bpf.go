package bpf

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"
	"unsafe"

	"github.com/iovisor/gobpf/elf"
)

/*
#include <linux/bpf.h>
*/
import "C"

const assetpath string = "out/cgnet.o"
const mapname string = "count"

var (
	zero       uint32 = 0
	packetsKey uint32 = zero
	bytesKey   uint32 = 1

	b *elf.Module
)

func Setup(cgroupPath string) error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// for quick development on the bpf program
	if path := os.Getenv("BPF_PROG_PATH"); path != "" {
		fmt.Printf("using: %s\n", path)
		b = elf.NewModule(path)
	} else {
		reader := bytes.NewReader(MustAsset(assetpath))
		b = elf.NewModuleFromReader(reader)
	}
	if b == nil {
		return fmt.Errorf("System doesn't support BPF")
	}

	if err := b.Load(nil); err != nil {
		return fmt.Errorf("%s", err)
	}

	for prog := range b.IterCgroupProgram() {
		if err := elf.AttachCgroupProgram(prog, cgroupPath, elf.EgressType); err != nil {
			return fmt.Errorf("%s", err)
		}
	}

	mp := b.Map(mapname)
	if mp == nil {
		return fmt.Errorf("Can't find map '%s'", mapname)
	}

	if err := b.UpdateElement(mp, unsafe.Pointer(&packetsKey), unsafe.Pointer(&zero), C.BPF_ANY); err != nil {
		return fmt.Errorf("error updating map: %s", err)
	}

	if err := b.UpdateElement(mp, unsafe.Pointer(&bytesKey), unsafe.Pointer(&zero), C.BPF_EXIST); err != nil {
		return fmt.Errorf("error updating map: %s", err)
	}

	fmt.Println("Ready.")
	return nil
}

func UpdateLoop(quit chan struct{}) error {
	mp := b.Map(mapname)
	if mp == nil {
		return fmt.Errorf("Can't find map '%s'", mapname)
	}

	var packets, bytes uint64

	for {
		select {
		case <-quit:
			return nil
		case <-time.After(1000 * time.Millisecond):
			if err := updateElements(mp, packets, bytes); err != nil {
				return err
			}
			fmt.Printf("cgroup received %d packets (%d bytes)\n", packets, bytes)
		}
	}
}

func updateElements(mp *elf.Map, packets, bytes uint64) error {
	if err := b.LookupElement(mp, unsafe.Pointer(&packetsKey), unsafe.Pointer(&packets)); err != nil {
		return fmt.Errorf("error looking up in map: %s", err)
	}

	if err := b.LookupElement(mp, unsafe.Pointer(&bytesKey), unsafe.Pointer(&bytes)); err != nil {
		return fmt.Errorf("error looking up in map: %s", err)
	}

	return nil
}
