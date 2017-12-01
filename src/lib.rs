extern crate regex;

use std::slice;

use regex::bytes::Regex;


#[no_mangle]
pub extern fn is_match(ptr: *const u8, len: u32, out: &mut bool) {
    let buf = unsafe {
        slice::from_raw_parts(ptr, len as usize)
    };

    println!("matching a buffer of length: {}", len);

    let re = Regex::new(r"^\d{4}-\d{2}-\d{2}$").unwrap();
    *out = re.is_match(buf);
}
