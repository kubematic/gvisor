// Copyright 2018 Google Inc.
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

package fs

import (
	"gvisor.googlesource.com/gvisor/pkg/sentry/context"
	"gvisor.googlesource.com/gvisor/pkg/sentry/usermem"
	"gvisor.googlesource.com/gvisor/pkg/syserror"
)

// MockInodeOperations implements InodeOperations for testing Inodes.
type MockInodeOperations struct {
	InodeOperations

	UAttr UnstableAttr

	createCalled          bool
	createDirectoryCalled bool
	createLinkCalled      bool
	renameCalled          bool
	walkCalled            bool
}

// NewMockInode returns a mock *Inode using MockInodeOperations.
func NewMockInode(ctx context.Context, msrc *MountSource, sattr StableAttr) *Inode {
	return NewInode(NewMockInodeOperations(ctx), msrc, sattr)
}

// NewMockInodeOperations returns a *MockInodeOperations.
func NewMockInodeOperations(ctx context.Context) *MockInodeOperations {
	return &MockInodeOperations{
		UAttr: WithCurrentTime(ctx, UnstableAttr{
			Perms: FilePermsFromMode(0777),
		}),
	}
}

// MockMountSourceOps implements fs.MountSourceOperations.
type MockMountSourceOps struct {
	MountSourceOperations
	keep       bool
	revalidate bool
}

// NewMockMountSource returns a new *MountSource using MockMountSourceOps.
func NewMockMountSource(cache *DirentCache) *MountSource {
	var keep bool
	if cache != nil {
		keep = cache.maxSize > 0
	}
	return &MountSource{
		MountSourceOperations: &MockMountSourceOps{keep: keep},
		fscache:               cache,
		children:              make(map[*MountSource]struct{}),
	}
}

// Revalidate implements fs.MountSourceOperations.Revalidate.
func (n *MockMountSourceOps) Revalidate(context.Context, *Inode) bool {
	return n.revalidate
}

// Keep implements fs.MountSourceOperations.Keep.
func (n *MockMountSourceOps) Keep(dirent *Dirent) bool {
	return n.keep
}

// WriteOut implements fs.InodeOperations.WriteOut.
func (n *MockInodeOperations) WriteOut(context.Context, *Inode) error {
	return nil
}

// UnstableAttr implements fs.InodeOperations.UnstableAttr.
func (n *MockInodeOperations) UnstableAttr(context.Context, *Inode) (UnstableAttr, error) {
	return n.UAttr, nil
}

// IsVirtual implements fs.InodeOperations.IsVirtual.
func (n *MockInodeOperations) IsVirtual() bool {
	return false
}

// Lookup implements fs.InodeOperations.Lookup.
func (n *MockInodeOperations) Lookup(ctx context.Context, dir *Inode, p string) (*Dirent, error) {
	n.walkCalled = true
	return NewDirent(NewInode(&MockInodeOperations{}, dir.MountSource, StableAttr{}), p), nil
}

// SetPermissions implements fs.InodeOperations.SetPermissions.
func (n *MockInodeOperations) SetPermissions(context.Context, *Inode, FilePermissions) bool {
	return false
}

// SetOwner implements fs.InodeOperations.SetOwner.
func (*MockInodeOperations) SetOwner(context.Context, *Inode, FileOwner) error {
	return syserror.EINVAL
}

// SetTimestamps implements fs.InodeOperations.SetTimestamps.
func (n *MockInodeOperations) SetTimestamps(context.Context, *Inode, TimeSpec) error {
	return nil
}

// Create implements fs.InodeOperations.Create.
func (n *MockInodeOperations) Create(ctx context.Context, dir *Inode, p string, flags FileFlags, perms FilePermissions) (*File, error) {
	n.createCalled = true
	d := NewDirent(NewInode(&MockInodeOperations{}, dir.MountSource, StableAttr{}), p)
	return &File{Dirent: d}, nil
}

// CreateLink implements fs.InodeOperations.CreateLink.
func (n *MockInodeOperations) CreateLink(_ context.Context, dir *Inode, oldname string, newname string) error {
	n.createLinkCalled = true
	return nil
}

// CreateDirectory implements fs.InodeOperations.CreateDirectory.
func (n *MockInodeOperations) CreateDirectory(context.Context, *Inode, string, FilePermissions) error {
	n.createDirectoryCalled = true
	return nil
}

// Rename implements fs.InodeOperations.Rename.
func (n *MockInodeOperations) Rename(ctx context.Context, oldParent *Inode, oldName string, newParent *Inode, newName string) error {
	n.renameCalled = true
	return nil
}

// Check implements fs.InodeOperations.Check.
func (n *MockInodeOperations) Check(ctx context.Context, inode *Inode, p PermMask) bool {
	return ContextCanAccessFile(ctx, inode, p)
}

// Release implements fs.InodeOperations.Release.
func (n *MockInodeOperations) Release(context.Context) {}

// Truncate implements fs.InodeOperations.Truncate.
func (n *MockInodeOperations) Truncate(ctx context.Context, inode *Inode, size int64) error {
	return nil
}

// DeprecatedPwritev implements fs.InodeOperations.DeprecatedPwritev.
func (n *MockInodeOperations) DeprecatedPwritev(context.Context, usermem.IOSequence, int64) (int64, error) {
	return 0, nil
}

// DeprecatedReaddir implements fs.InodeOperations.DeprecatedReaddir.
func (n *MockInodeOperations) DeprecatedReaddir(context.Context, *DirCtx, int) (int, error) {
	return 0, nil
}

// Remove implements fs.InodeOperations.Remove.
func (n *MockInodeOperations) Remove(context.Context, *Inode, string) error {
	return nil
}

// RemoveDirectory implements fs.InodeOperations.RemoveDirectory.
func (n *MockInodeOperations) RemoveDirectory(context.Context, *Inode, string) error {
	return nil
}

// Getlink implements fs.InodeOperations.Getlink.
func (n *MockInodeOperations) Getlink(context.Context, *Inode) (*Dirent, error) {
	return nil, syserror.ENOLINK
}
