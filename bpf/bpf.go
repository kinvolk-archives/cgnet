/*
Copyright 2017 Kinvolk GmbH

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package bpf

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"unsafe"

	bpf "github.com/iovisor/gobpf/elf"
)

/*
#include <linux/bpf.h>
*/
import "C"

const (
	assetPath string = "out/cgnet.o"
	mapName   string = "count_map"
)

func Attach(cgroupPath string) (*Controller, error) {
	b, err := initModule()
	if err != nil {
		return nil, err
	}

	for prog := range b.IterCgroupProgram() {
		if err := bpf.AttachCgroupProgram(prog, cgroupPath, bpf.EgressType); err != nil {
			return nil, fmt.Errorf("error attaching to cgroup %s: %s", cgroupPath, err)
		}
	}

	if err := initMap(b); err != nil {
		return nil, err
	}

	ctl := &Controller{
		cgroup: cgroupPath,
		module: b,
		quit:   make(chan struct{}),
	}
	return ctl, nil
}

func initModule() (*bpf.Module, error) {
	var b *bpf.Module
	if path := os.Getenv("BPF_PROG_PATH"); path != "" {
		// for development
		b = bpf.NewModule(path)
	} else {
		// from assets
		reader := bytes.NewReader(MustAsset(assetPath))
		b = bpf.NewModuleFromReader(reader)
	}
	if b == nil {
		return nil, fmt.Errorf("system doesn't seem to support BPF")
	}

	if err := b.Load(nil); err != nil {
		return nil, fmt.Errorf("loading module failed: %s", err)
	}
	return b, nil
}

func initMap(b *bpf.Module) error {
	if err := update(b, packetsKey, 0); err != nil {
		return fmt.Errorf("error updating map: %s", err)
	}
	if err := update(b, bytesKey, 0); err != nil {
		return fmt.Errorf("error updating map: %s", err)
	}
	return nil
}

func update(b *bpf.Module, key uint32, value uint64) error {

	mp := b.Map(mapName)
	if err := b.UpdateElement(mp, unsafe.Pointer(&key), unsafe.Pointer(&value), C.BPF_ANY); err != nil {
		return err
	}
	return nil
}

func lookup(b *bpf.Module, key uint32) (uint64, error) {
	mp := b.Map(mapName)
	var value uint64
	if err := b.LookupElement(mp, unsafe.Pointer(&key), unsafe.Pointer(&value)); err != nil {
		return 0, err
	}
	return value, nil
}
