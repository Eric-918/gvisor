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

package proc

import (
	"bytes"
	"fmt"

	"gvisor.dev/gvisor/pkg/abi/linux"
	"gvisor.dev/gvisor/pkg/sentry/context"
	"gvisor.dev/gvisor/pkg/sentry/fsimpl/kernfs"
	"gvisor.dev/gvisor/pkg/sentry/kernel"
	"gvisor.dev/gvisor/pkg/sentry/kernel/auth"
	"gvisor.dev/gvisor/pkg/sentry/vfs"
)

func newSysDir(root *auth.Credentials, inoGen InoGenerator) *kernfs.Dentry {
	return kernfs.NewStaticDir(root, inoGen.NextIno(), dirPerm, map[string]*kernfs.Dentry{
		"kernel": kernfs.NewStaticDir(root, inoGen.NextIno(), dirPerm, map[string]*kernfs.Dentry{
			"hostname": kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &hostnameData{}),
			"shmall":   kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, shmData(linux.SHMALL)),
			"shmmax":   kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, shmData(linux.SHMMAX)),
			"shmmni":   kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, shmData(linux.SHMMNI)),
		}),
		"vm": kernfs.NewStaticDir(root, inoGen.NextIno(), dirPerm, map[string]*kernfs.Dentry{
			"mmap_min_addr":     kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &mmapMinAddrData{}),
			"overcommit_memory": kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "0\n"}),
		}),
		"net": newSysNetDir(root, inoGen),
	})
}

func newSysNetDir(root *auth.Credentials, inoGen InoGenerator) *kernfs.Dentry {
	return kernfs.NewStaticDir(root, inoGen.NextIno(), dirPerm, map[string]*kernfs.Dentry{
		"net": kernfs.NewStaticDir(root, inoGen.NextIno(), dirPerm, map[string]*kernfs.Dentry{
			"ipv4": kernfs.NewStaticDir(root, inoGen.NextIno(), dirPerm, map[string]*kernfs.Dentry{
				// Add tcp_sack.
				// TODO(gvisor.dev/issue/1195): tcp_sack allows write(2)
				// "tcp_sack": newTCPSackInode(ctx, msrc, s),

				// The following files are simple stubs until they are implemented in
				// netstack, most of these files are configuration related. We use the
				// value closest to the actual netstack behavior or any empty file, all
				// of these files will have mode 0444 (read-only for all users).
				"ip_local_port_range":     kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "16000   65535"}),
				"ip_local_reserved_ports": kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{}),
				"ipfrag_time":             kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "30"}),
				"ip_nonlocal_bind":        kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "0"}),
				"ip_no_pmtu_disc":         kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "1"}),

				// tcp_allowed_congestion_control tell the user what they are able to
				// do as an unprivledged process so we leave it empty.
				"tcp_allowed_congestion_control":   kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: ""}),
				"tcp_available_congestion_control": kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "reno"}),
				"tcp_congestion_control":           kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "reno"}),

				// Many of the following stub files are features netstack doesn't
				// support. The unsupported features return "0" to indicate they are
				// disabled.
				"tcp_base_mss":              kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "1280"}),
				"tcp_dsack":                 kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "0"}),
				"tcp_early_retrans":         kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "0"}),
				"tcp_fack":                  kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "0"}),
				"tcp_fastopen":              kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "0"}),
				"tcp_fastopen_key":          kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: ""}),
				"tcp_invalid_ratelimit":     kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "0"}),
				"tcp_keepalive_intvl":       kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "0"}),
				"tcp_keepalive_probes":      kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "0"}),
				"tcp_keepalive_time":        kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "7200"}),
				"tcp_mtu_probing":           kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "0"}),
				"tcp_no_metrics_save":       kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "1"}),
				"tcp_probe_interval":        kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "0"}),
				"tcp_probe_threshold":       kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "0"}),
				"tcp_retries1":              kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "3"}),
				"tcp_retries2":              kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "15"}),
				"tcp_rfc1337":               kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "1"}),
				"tcp_slow_start_after_idle": kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "1"}),
				"tcp_synack_retries":        kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "5"}),
				"tcp_syn_retries":           kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "3"}),
				"tcp_timestamps":            kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "1"}),
			}),
			"core": kernfs.NewStaticDir(root, inoGen.NextIno(), dirPerm, map[string]*kernfs.Dentry{
				"default_qdisc": kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "pfifo_fast"}),
				"message_burst": kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "10"}),
				"message_cost":  kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "5"}),
				"optmem_max":    kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "0"}),
				"rmem_default":  kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "212992"}),
				"rmem_max":      kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "212992"}),
				"somaxconn":     kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "128"}),
				"wmem_default":  kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "212992"}),
				"wmem_max":      kernfs.NewDynamicBytesFile(root, inoGen.NextIno(), filePerm, &vfs.StaticData{Data: "212992"}),
			}),
		}),
	})
}

// mmapMinAddrData implements vfs.DynamicBytesSource for
// /proc/sys/vm/mmap_min_addr.
//
// +stateify savable
type mmapMinAddrData struct {
	k *kernel.Kernel
}

var _ vfs.DynamicBytesSource = (*mmapMinAddrData)(nil)

// Generate implements vfs.DynamicBytesSource.Generate.
func (d *mmapMinAddrData) Generate(ctx context.Context, buf *bytes.Buffer) error {
	fmt.Fprintf(buf, "%d\n", d.k.Platform.MinUserAddress())
	return nil
}

// hostnameData implements vfs.DynamicBytesSource for /proc/sys/kernel/hostname.
//
// +stateify savable
type hostnameData struct{}

// Generate implements vfs.DynamicBytesSource.Generate.
func (*hostnameData) Generate(ctx context.Context, buf *bytes.Buffer) error {
	utsns := kernel.UTSNamespaceFromContext(ctx)
	buf.WriteString(utsns.HostName())
	buf.WriteString("\n")
	return nil
}
