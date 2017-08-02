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
	"context"
	"time"

	bpf "github.com/iovisor/gobpf/elf"
)

/*
#include <linux/bpf.h>
*/
import "C"

const (
	packetsKey uint32 = 0
	bytesKey   uint32 = 1
)

type Controller struct {
	cgroup string
	module *bpf.Module
	quit   chan struct{}

	packetsHandler func(uint64) error
	bytesHandler   func(uint64) error
}

func (c *Controller) Stop() {
	c.quit <- struct{}{}
	if err := c.module.Close(); err != nil {
		Log.Error("error closing bpf.Module", "err", err)
	}
}

func (c *Controller) SetPacketsHandler(h func(uint64) error) {
	c.packetsHandler = h
}
func (c *Controller) SetBytesHandler(h func(uint64) error) {
	c.bytesHandler = h
}

func (c *Controller) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.quit:
			return
		case <-time.After(1 * time.Second):
			// TODO this needs to be solved differently
			// Maybe sending new values on individual channels?
			packets, err := lookup(c.module, packetsKey)
			if err != nil {
				Log.Error("lookup failed", "err", err)
				continue
			}
			if err := c.packetsHandler(packets); err != nil {
				Log.Error("packetsHandler failed", "err", err)
			}

			bytes, err := lookup(c.module, bytesKey)
			if err != nil {
				Log.Error("lookup failed", "err", err)
				continue
			}
			if err := c.bytesHandler(bytes); err != nil {
				Log.Error("bytesHandler failed", "err", err)
			}
		}
	}
}
