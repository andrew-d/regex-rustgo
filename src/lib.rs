extern crate regex;

use regex::bytes::Regex;


#[no_mangle]
pub extern fn is_match(buf: &[u8; 64], out: &mut bool) {
    let re = Regex::new(r"^\d{4}-\d{2}-\d{2}$").unwrap();
    println!("is_match = {}", re.is_match(buf));
    *out = re.is_match(buf);
}
