#![feature(alloc_system)]
extern crate alloc_system;

extern crate regex;

use std::slice;

use regex::bytes::Regex;


fn is_match_internal(buf: &[u8]) -> bool {
    let re = Regex::new(r"^\d{4}-\d{2}-\d{2}$").unwrap();
    re.is_match(buf)
}

#[no_mangle]
pub extern fn is_match(ptr: *const u8, len: u32, out: &mut bool) {
    let buf = unsafe {
        slice::from_raw_parts(ptr, len as usize)
    };

    println!("matching a buffer of length: {}", len);
    *out = is_match_internal(buf);
}
