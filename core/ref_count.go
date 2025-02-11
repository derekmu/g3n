package core

type IRefCount interface {
	Incref()
	Decref() bool
}

var _ IRefCount = &RefCount{}

type RefCount struct {
	permanent bool
	refcount  int
}

// Incref increments the reference count for this object.
func (r *RefCount) Incref() {
	if !r.permanent {
		r.refcount++
	}
}

// Decref decrements the reference count for this object and returns whether the object should be disposed.
func (r *RefCount) Decref() bool {
	if r.permanent {
		return false
	} else {
		r.refcount = max(0, r.refcount-1)
		return r.refcount == 0
	}
}

// SetPermanent sets the permanency flag for this object.
// Permanent objects aren't to be disposed of.
// This is meant for objects that may be reused or stored in caches.
func (r *RefCount) SetPermanent(permanent bool) {
	r.permanent = permanent
}
