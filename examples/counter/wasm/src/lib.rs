#![feature(vec_into_raw_parts)]

extern crate svm_sdk;

use svm_sdk::traits::Encoder;
use svm_sdk::CallData;

const VAR_ID: u32 = 0;

#[link(wasm_import_module = "svm")]
extern "C" {
    fn svm_calldata_offset() -> u32;

    fn svm_calldata_len() -> u32;

    fn svm_set_returndata(offset: u32, length: u32);

    fn svm_get32(var_id: u32) -> u32;

    fn svm_set32(var_id: u32, value: u32);

    fn svm_log(offset: u32, length: u32, code: u32);kk
}

#[link(wasm_import_module = "host")]
extern "C" {
    fn add(a: u32, b: u32) -> u32;

    fn mul(a: u32, b: u32) -> u32;
}

#[no_mangle]
pub extern "C" fn svm_alloc(size: i32) -> i32 {
    let ptr = svm_sdk::alloc(size as usize);

    ptr.offset() as i32
}

#[no_mangle]
pub extern "C" fn initialize() {
    let bytes = get_calldata();

    let mut calldata = CallData::new(bytes);

    let initial: u32 = calldata.next_1();

    unsafe {
        svm_set32(VAR_ID, initial);
    }
}

#[no_mangle]
pub extern "C" fn counter_add() {
    let calldata = get_calldata();
    let mut calldata = CallData::new(calldata);

    let arg1: u32 = calldata.next_1();

    unsafe {
        let counter = svm_get32(VAR_ID);

        log("invoking `add` host import function", 100);
        let new_counter = add(counter, arg1);

        svm_set32(VAR_ID, new_counter.clone());

        let mut buf = Vec::new();
        let results = vec![counter, new_counter];
        results.encode(&mut buf);

        let (ptr, len, _cap) = buf.into_raw_parts();
        svm_set_returndata(ptr as u32, len as u32);
    }
}

#[no_mangle]
pub extern "C" fn counter_mul() {
    let calldata = get_calldata();
    let mut calldata = CallData::new(calldata);

    let arg1: u32 = calldata.next_1();

    unsafe {
        let counter = svm_get32(VAR_ID);

        log("invoking `mul` host import function", 100);
        let new_counter = mul(counter, arg1);

        svm_set32(VAR_ID, new_counter.clone());

        let mut buf = Vec::new();
        let results = vec![counter, new_counter];
        results.encode(&mut buf);

        let (ptr, len, _cap) = buf.into_raw_parts();
        svm_set_returndata(ptr as u32, len as u32);
    }
}

fn get_calldata() -> &'static [u8] {
    unsafe {
        let ptr = svm_calldata_offset();
        let len = svm_calldata_len();

        core::slice::from_raw_parts(ptr as *const u8, len as usize)
    }
}

//fn set_returndata() -> &'static [u8] {
//    unsafe {
//        let ptr = svm_calldata_offset();
//        let len = svm_calldata_len();
//
//        core::slice::from_raw_parts(ptr as *const u8, len as usize)
//    }
//}

fn log(msg: &str, code: u8) {
    unsafe {
        let offset = msg.as_ptr() as u32;
        let len = msg.len() as u32;

        svm_log(offset, len, code as u32)
    }
}
