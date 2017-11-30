// #![crate_type="lib"]

// #![feature(no_std)]
// #![no_std]
// #![no_main]

// #![feature(lang_items, compiler_builtins_lib, core_intrinsics)]

// use core::intrinsics;

// #[allow(private_no_mangle_fns)] #[no_mangle] // rust-lang/rust#38281
// #[lang = "panic_fmt"] fn panic_fmt() -> ! { unsafe { intrinsics::abort() } }
// #[lang = "eh_personality"] extern fn eh_personality() {}

extern crate regex;

use regex::bytes::Regex;


#[no_mangle]
pub extern fn is_match(buf: &[u8; 64], out: &mut bool) {
    let re = Regex::new(r"^\d{4}-\d{2}-\d{2}$").unwrap();
    *out = re.is_match(buf);
}
