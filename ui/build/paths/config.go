// Copyright 2018 Google Inc. All rights reserved.
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

package paths

import "runtime"

type PathConfig struct {
	// Whether to create the symlink in the new PATH for this tool.
	Symlink bool

	// Whether to log about usages of this tool to the soong.log
	Log bool

	// Whether to exit with an error instead of invoking the underlying tool.
	Error bool

	// Whether we use a linux-specific prebuilt for this tool. On Darwin,
	// we'll allow the host executable instead.
	LinuxOnlyPrebuilt bool
}

// These binaries can be run from $PATH, nonhermetically. There should be as
// few as possible of these, since this means that the build depends on tools
// that are not shipped in the source tree and whose behavior is therefore
// unpredictable.
var Allowed = PathConfig{
	Symlink: true,
	Log:     false,
	Error:   false,
}

// This tool is specifically disallowed and calling it will result in an
// "executable no found" error.
var Forbidden = PathConfig{
	Symlink: false,
	Log:     true,
	Error:   true,
}

// This tool is allowed, but access to it will be logged.
var Log = PathConfig{
	Symlink: true,
	Log:     true,
	Error:   false,
}

// The configuration used if the tool is not listed in the config below.
// Currently this will create the symlink, but log and error when it's used. In
// the future, I expect the symlink to be removed, and this will be equivalent
// to Forbidden. This applies to every tool not specifically mentioned in the
// configuration.
var Missing = PathConfig{
	Symlink: true,
	Log:     true,
	Error:   true,
}

// This is used for binaries for which we have prebuilt versions, but only for
// Linux. Thus, their execution from $PATH is only allowed on Mac OS.
var LinuxOnlyPrebuilt = PathConfig{
	Symlink:           false,
	Log:               true,
	Error:             true,
	LinuxOnlyPrebuilt: true,
}

func GetConfig(name string) PathConfig {
	if config, ok := Configuration[name]; ok {
		return config
	}
	return Missing
}

// This list specifies whether a particular binary from $PATH is allowed to be
// run during the build. For more documentation, see path_interposer.go .
var Configuration = map[string]PathConfig{
	"bash":           Allowed,
	"brotli":	  Allowed,
	"ccache": 	  Allowed,
	"cpio":   	  Allowed,
        "curl":    	  Allowed,
	"date":    	  Allowed,
	"diff":           Allowed,
	"dlv":            Allowed,
	"expr":           Allowed,
	"flock":   	  Allowed,
	"fuser":          Allowed,
	"gcert":          Allowed,
	"gcertstatus":    Allowed,
	"gcloud":         Allowed,
	"getopt":         Allowed,
	"git":            Allowed,
	"hexdump":        Allowed,
	"jar":            Allowed,
	"java":           Allowed,
	"javap":          Allowed,
	"locale":	  Allowed,
	"lsof":           Allowed,
	"ld.lld":  	  Allowed,
	"llvm-ar": 	  Allowed,
        "nproc":   	  Allowed,
	"openssl":        Allowed,
	"patch":   	  Allowed,
	"pkg-config": 	  Allowed,
	"pstree":         Allowed,
	"python3": 	  Allowed,
	"python3.6": 	  Allowed,
	"python3.7": 	  Allowed,
	"python3.8": 	  Allowed,
	"python3.9": 	  Allowed,
        "python3.10":     Allowed,
	"repo":    	  Allowed,
	"rsync":          Allowed,
	"sh":             Allowed,
	"stubby":         Allowed,
        "tar":            Allowed,
	"tr":             Allowed,
	"unzip":          Allowed,
	"zip":            Allowed,
	"zipdetails":     Allowed,
	"arm-linux-androidkernel-as": Allowed,
	"arm-linux-androidkernel-ld": Allowed,

	// Host toolchain is removed. In-tree toolchain should be used instead.
	// GCC also can't find cc1 with this implementation.
	"ar":         Forbidden,
	"as":         Forbidden,
	"cc":         Forbidden,
	"clang":      Forbidden,
	"clang++":    Forbidden,
	"gcc":        Forbidden,
	"g++":        Forbidden,
	"ld":         Allowed,// HACK
	"ld.bfd":     Forbidden,
	"ld.gold":    Forbidden,
	"perl":		  Allowed,// HACK

	// These are toybox tools that only work on Linux.
	"pgrep": LinuxOnlyPrebuilt,
	"pkill": LinuxOnlyPrebuilt,
	"ps":    LinuxOnlyPrebuilt,
}

func init() {
	if runtime.GOOS == "darwin" {
		Configuration["sw_vers"] = Allowed
		Configuration["xcrun"] = Allowed

		// We don't have darwin prebuilts for some tools,
		// so allow the host versions.
		for name, config := range Configuration {
			if config.LinuxOnlyPrebuilt {
				Configuration[name] = Allowed
			}
		}
	}
}
