#![feature(alloc_system)]
extern crate alloc_system;

extern crate regex;

use std::mem;
// use std::os::raw::c_void;
use std::slice;
use std::str;

use regex::bytes::Regex;


#[no_mangle]
pub unsafe extern fn rust_compile(ptr: *const u8, len: u32, out: *mut *mut Regex) {
    let s = {
        // Create a byte slice from the given input
        let sl = slice::from_raw_parts(ptr, len as usize);

        // Turn it into a string; Go guarantees that this works.
        str::from_utf8_unchecked(sl)
    };

    *out = Box::into_raw(rust_compile_internal(s))
}

fn rust_compile_internal(s: &str) -> Box<Regex> {
    // println!("creating regex from string: {}", s);
    Box::new(Regex::new(s).unwrap())
}


#[no_mangle]
pub unsafe extern fn rust_free(ptr: *mut Regex) {
    let ptr: Box<Regex> = Box::from_raw(ptr);
    // println!("freeing regex");
    mem::drop(ptr);
}


fn is_match_internal(re: &Regex, buf: &[u8]) -> bool {
    re.is_match(buf)
}

#[no_mangle]
pub unsafe extern fn is_match(re: *mut Regex, ptr: *const u8, len: u32, out: &mut bool) {
    if let Some(re) = re.as_ref() {
        let buf = slice::from_raw_parts(ptr, len as usize);
        // println!("matching a buffer of length {}: {:?}", len, buf);
        *out = is_match_internal(re, buf);
    }
}


fn find_index_internal(re: &Regex, buf: &[u8], out: &mut [u32]) -> bool {
    if let Some(mat) = re.find(buf) {
        out[0] = mat.start() as u32;
        out[1] = mat.end() as u32;
        true
    } else {
        false
    }
}

#[no_mangle]
pub unsafe extern fn find_index(re: *mut Regex, ptr: *const u8, len: u32, out: *mut u32, success: &mut bool) {
    if let Some(re) = re.as_ref() {
        let buf = slice::from_raw_parts(ptr, len as usize);
        let out = slice::from_raw_parts_mut(out, 2);

        *success = find_index_internal(re, buf, out);
    }
}
