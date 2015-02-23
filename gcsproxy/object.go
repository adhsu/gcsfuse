// Copyright 2015 Google Inc. All Rights Reserved.
// Author: jacobsa@google.com (Aaron Jacobs)

package gcsproxy

import (
	"io"

	"github.com/jacobsa/gcloud/gcs"
	"google.golang.org/cloud/storage"
)

// A view on an object in GCS that allows random access reads and writes.
//
// Reads may involve reading from a local cache. Writes are buffered locally
// until the Sync method is called, at which time a new generation of the
// object is created.
//
// All methods are safe for concurrent access. Concurrent readers and writers
// within process receive the same guarantees as with POSIX files.
type ProxyObject struct {
}

var _ io.ReaderAt = &ProxyObject{}
var _ io.WriterAt = &ProxyObject{}

// Create a new view on the GCS object with the given name. The remote object
// is assumed to be non-existent, so that the local contents are empty. Use
// NoteLatest to change that if necessary.
func NewProxyObject(
	bucket gcs.Bucket,
	name string) (*ProxyObject, error)

// Inform the proxy object of the most recently observed generation of the
// object of interest in GCS.
//
// If this is no newer than the newest generation that has previously been
// observed, it is ignored. Otherwise, it becomes the definitive source of data
// for the object. Any local-only state is clobbered, including local
// modifications.
func (po *ProxyObject) NoteLatest(o storage.Object) error

// Return the current size in bytes of our view of the content.
func (po *ProxyObject) Size() uint64

// Make a random access read into our view of the content. May block for
// network access.
func (po *ProxyObject) ReadAt(buf []byte, offset int64) (int, error)

// Make a random access write into our view of the content. May block for
// network access. Not guaranteed to be reflected remotely until after Sync is
// called successfully.
func (po *ProxyObject) WriteAt(buf []byte, offset int64) (int, error)

// Truncate our view of the content to the given number of bytes, extending if
// n is greater than Size(). May block for network access. Not guaranteed to be
// reflected remotely until after Sync is called successfully.
func (po *ProxyObject) Truncate(n uint64) error

// Ensure that the remote object reflects the local state, returning a record
// for a generation that does. Clobbers the remote version.
func (po *ProxyObject) Sync() (storage.Object, error)
