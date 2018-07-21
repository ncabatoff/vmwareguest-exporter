/*
 * Product : go-vmguestlib
 * Date : 2015
 * Author : Xavier Lucas
 * Licence : GNU GPL2
 *
 * These functions permits to workaround the missing implementation of C function pointers
 * in Go. We define a new pointer type for each native function return type and declare
 * a proxy function per return type. All native target functions must have the same
 * parameters than the proxy function calling them.
 */

#ifndef _VM_GUEST_LIB_PROXY_H
#define _VM_GUEST_LIB_PROXY_H

#include <vm_basic_types.h>
#include <vmGuestLib.h>

#ifdef __cplusplus
extern "C" {
#endif

typedef VMGuestLibError (*p_uint32_f) ();
typedef VMGuestLibError (*p_uint64_f) ();

static inline VMGuestLibError proxy_uint32_f(p_uint32_f f, VMGuestLibHandle h, uint32* v)
{
    return f(h, v);
}

static inline VMGuestLibError proxy_uint64_f(p_uint64_f f, VMGuestLibHandle h, uint64* v)
{
    return f(h, v);
}

#ifdef __cplusplus
}
#endif

#endif /* _VM_GUEST_LIB_PROXY_H */
