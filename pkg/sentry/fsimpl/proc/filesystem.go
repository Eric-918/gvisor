// Copyright 2019 The gVisor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package proc implements a partial in-memory file system for procfs.
package proc

import (
	"fmt"

	"gvisor.dev/gvisor/pkg/sentry/context"
	"gvisor.dev/gvisor/pkg/sentry/fsimpl/kernfs"
	"gvisor.dev/gvisor/pkg/sentry/kernel"
	"gvisor.dev/gvisor/pkg/sentry/kernel/auth"
	"gvisor.dev/gvisor/pkg/sentry/vfs"
)

const (
	filePerm = 0444
	dirPerm  = 0555
)

// procFSType is the factory class for procfs.
//
// +stateify savable
type procFSType struct{}

var _ vfs.FilesystemType = (*procFSType)(nil)

// GetFilesystem implements vfs.FilesystemType.
func (ft *procFSType) GetFilesystem(ctx context.Context, vfsObj *vfs.VirtualFilesystem, creds *auth.Credentials, source string, opts vfs.GetFilesystemOptions) (*vfs.Filesystem, *vfs.Dentry, error) {
	k := kernel.KernelFromContext(ctx)
	if k == nil {
		return nil, nil, fmt.Errorf("procfs requires a kernel")
	}
	pidns := kernel.PIDNamespaceFromContext(ctx)
	if pidns == nil {
		return nil, nil, fmt.Errorf("procfs requires a PID namespace")
	}

	procfs := &kernfs.Filesystem{}
	procfs.VFSFilesystem().Init(vfsObj, procfs)

	var cgroups map[string]string
	if opts.InternalData != nil {
		cgroups = opts.InternalData.(map[string]string)
	}

	_, dentry := newTasksInode(procfs, k, pidns, cgroups)
	return procfs.VFSFilesystem(), dentry.VFSDentry(), nil
}
