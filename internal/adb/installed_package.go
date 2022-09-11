// Copyright (C) 2015 The Android Open Source Project
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package adb

import (
	"fmt"
	"sort"
	"strings"
)

// https://android.googlesource.com/platform/bionic/+/master/libc/include/sys/system_properties.h#38
// Actually 32, but that includes null-terminator.
const maxPropName = 31

type InstalledPackage struct {
	Name    string    // Name of the package.
	Device  *Device   // The device this package is installed on.
	Actions []*Action // The actions this package supports.
	ABI     string    // The ABI of the package or empty
}

// WrapProperties returns the list of wrap-properties for the given installed
// package.
func (p *InstalledPackage) WrapProperties() ([]string, error) {
	list, err := p.Device.Command("getprop", p.wrapPropName()).Call()
	return strings.Fields(list), err
}

// WrapProperties sets the list of wrap-properties for the given installed
// package.
func (p *InstalledPackage) SetWrapProperties(props ...string) error {
	arg := strings.Join(props, " ")
	return p.Device.Command("setprop", p.wrapPropName(), arg).Run()
}

// Action represents an Android action that can be sent as an intent.
type Action struct {
	Name     string // Example: android.intent.action.MAIN
	Package  *InstalledPackage
	Activity string // Example: .FooBarActivity
}

// String returns the package name.
func (p *InstalledPackage) String() string {
	return p.Name
}
func (p *InstalledPackage) wrapPropName() string {
	name := "wrap." + p.Name
	if len(name) > maxPropName {
		name = name[:maxPropName]
	}
	return name
}
func (a *Action) String() string {
	return a.Name + ":" + a.Package.Name + "/" + a.Activity
}

// InstalledPackages returns the sorted list of installed packages on the device.
func (d *Device) InstalledPackages() (Packages, error) {
	str, err := d.Command("dumpsys", "package").Call()
	if err != nil {
		return nil, err
	}
	return d.parsePackages(str)
}
func (d *Device) parsePackages(str string) (Packages, error) {
	tree := parseTabbedTree(str)
	activities := tree.find("Activity Resolver Table:")
	if activities == nil {
		return nil, fmt.Errorf("Could not find Activity Resolver Table in dumpsys")
	}
	actions := activities.find("Non-Data Actions:")
	if actions == nil {
		return nil, fmt.Errorf("Could not find Non-Data Actions in dumpsys")
	}
	packageMap := map[string]*InstalledPackage{}
	for _, action := range actions.children {
		for _, entry := range action.children {
			// 43178558 com.google.foo/.FooActivity filter 431d7db8
			// 43178558 com.google.foo/.FooActivity
			fields := strings.Fields(entry.text)
			if len(fields) < 2 {
				return nil, fmt.Errorf("Could not parse package: '%v'", entry.text)
			}
			component := fields[1]
			parts := strings.SplitN(component, "/", 2)
			name := parts[0]
			p, ok := packageMap[name]
			if !ok {
				p = &InstalledPackage{
					Name:    name,
					Device:  d,
					Actions: []*Action{},
				}
				packageMap[name] = p
			}
			p.Actions = append(p.Actions, &Action{
				Package:  p,
				Name:     strings.TrimRight(action.text, ":"),
				Activity: parts[1],
			})
		}
	}
	// Read the "Packages:" section if it is present and use it to set ABI
	packSection := tree.find("Packages:")
	if packSection != nil {
		for _, pack := range packSection.children {
			// Package [com.google.foo] (ffffffc):
			fields := strings.Fields(pack.text)
			if len(fields) != 3 {
				continue
			}
			name := strings.Trim(fields[1], "[]")
			ip, ok := packageMap[name]
			if !ok {
				// We didn't find an action for this package
				continue
			}
			for _, attr := range pack.children {
				// primaryCpuAbi=arm64-v8a
				// primaryCpuAbi=null
				av := strings.TrimSpace(attr.text)
				if !strings.HasPrefix(av, "primaryCpuAbi=") {
					continue
				}
				splits := strings.SplitN(av, "=", 2)
				if len(splits) < 2 {
					break
				}
				if splits[1] == "null" {
					// This means the package manager will select the platform ABI
					break
				}
				ip.ABI = splits[1]
			}
		}
	}
	packages := make(Packages, 0, len(packageMap))
	for _, p := range packageMap {
		packages = append(packages, p)
	}
	sort.Sort(packages)
	return packages, nil
}

type Packages []*InstalledPackage

func (l Packages) Len() int           { return len(l) }
func (l Packages) Less(i, j int) bool { return l[i].Name < l[j].Name }
func (l Packages) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

type treeNode struct {
	text     string
	children []*treeNode
	parent   *treeNode
	depth    int
}

func (t *treeNode) find(name string) *treeNode {
	if t == nil {
		return nil
	}
	for _, c := range t.children {
		if c.text == name {
			return c
		}
	}
	return nil
}
func parseTabbedTree(str string) *treeNode {
	head := &treeNode{depth: -1}
	for _, line := range strings.Split(str, "\n") {
		line = strings.TrimRight(line, "\r")
		if line == "" {
			continue
		}
		// Calculate the line's depth
		depth := 0
		for i, r := range line {
			if r == ' ' {
				depth++
			} else {
				line = line[i:]
				break
			}
		}
		// Find the insertion point
		for {
			if head.depth >= depth {
				head = head.parent
			} else {
				node := &treeNode{text: line, depth: depth, parent: head}
				head.children = append(head.children, node)
				head = node
				break
			}
		}
	}
	for head.parent != nil {
		head = head.parent
	}
	return head
}
