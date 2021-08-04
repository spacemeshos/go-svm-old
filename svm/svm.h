#ifndef SVM_H
#define SVM_H

#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>
#include "svm_types.h"

/**
 * FFI representation for function result type
 */
typedef enum {
  SVM_SUCCESS = 0,
  SVM_FAILURE = 1,
} svm_result_t;

/**
 * Validates syntactically a raw `deploy template` transaction.
 *
 * Should be called while the transaction is in the `mempool` of the Host.
 * In case the transaction isn't valid - the transaction should be discarded.
 *
 *
 * # Example
 *
 * ```rust, no_run
 * use svm_runtime_c_api::*;
 *
 * use svm_ffi::svm_byte_array;
 * use svm_types::Address;
 *
 * // allocate imports
 * let mut imports = testing::imports_alloc(0);
 *
 * // create runtime
 *
 * let mut kv = std::ptr::null_mut();
 * let res = unsafe { svm_memory_state_kv_create(&mut kv) };
 * assert!(res.is_ok());
 *
 * let mut runtime = std::ptr::null_mut();
 * let mut error = svm_byte_array::default();
 *
 * let res = unsafe { svm_memory_runtime_create(&mut runtime, kv, imports, &mut error) };
 * assert!(res.is_ok());
 *
 * let bytes = svm_byte_array::default();
 * let _res = unsafe { svm_validate_template(runtime, bytes, &mut error) };
 * ```
 *
 */
svm_result_t svm_validate_template(void *runtime, svm_byte_array bytes, svm_byte_array *error);

/**
 * Validates syntactically a raw `spawn app` transaction.
 *
 * Should be called while the transaction is in the `mempool` of the Host.
 * In case the transaction isn't valid - the transaction should be discarded.
 *
 *
 * # Example
 *
 * ```rust, no_run
 * use svm_runtime_c_api::*;
 *
 * use svm_ffi::svm_byte_array;
 * use svm_types::Address;
 *
 * // allocate imports
 * let mut imports = testing::imports_alloc(0);
 *
 * // create runtime
 * let mut kv = std::ptr::null_mut();
 * let res = unsafe { svm_memory_state_kv_create(&mut kv) };
 * assert!(res.is_ok());
 *
 * let mut runtime = std::ptr::null_mut();
 * let mut error = svm_byte_array::default();
 *
 * let res = unsafe { svm_memory_runtime_create(&mut runtime, kv, imports, &mut error) };
 * assert!(res.is_ok());
 *
 * let bytes = svm_byte_array::default();
 * let _res = unsafe { svm_validate_app(runtime, bytes, &mut error) };
 * ```
 *
 */
svm_result_t svm_validate_app(void *runtime, svm_byte_array bytes, svm_byte_array *error);

/**
 * Validates syntactically a raw `execute app` transaction.
 * Returns the `App` address that appears in the transaction.
 *
 * # Example
 *
 * ```rust, no_run
 * use svm_runtime_c_api::*;
 *
 * use svm_ffi::svm_byte_array;
 * use svm_types::Address;
 *
 * // allocate imports
 * let mut imports = testing::imports_alloc(0);
 *
 * // create runtime
 *
 * let mut kv = std::ptr::null_mut();
 * let res = unsafe { svm_memory_state_kv_create(&mut kv) };
 * assert!(res.is_ok());
 *
 * let mut runtime = std::ptr::null_mut();
 * let mut error = svm_byte_array::default();
 *
 * let res = unsafe { svm_memory_runtime_create(&mut runtime, kv, imports, &mut error) };
 * assert!(res.is_ok());
 *
 * let mut app_addr = svm_byte_array::default();
 * let bytes = svm_byte_array::default();
 * let _res = unsafe { svm_validate_tx(&mut app_addr, runtime, bytes, &mut error) };
 * ```
 *
 */
svm_result_t svm_validate_tx(svm_byte_array *app_addr,
                             void *runtime,
                             svm_byte_array bytes,
                             svm_byte_array *error);

/**
 * Allocates space for the host imports.
 * See `svm_imports_destroy` for freeing the imports.
 *
 * # Example
 *
 * ```rust
 * use svm_runtime_c_api::svm_imports_alloc;
 *
 * let count = 2;
 * let mut imports = std::ptr::null_mut();
 *
 * let res = unsafe { svm_imports_alloc(&mut imports, count) };
 * assert!(res.is_ok());
 * ```
 *
 */
svm_result_t svm_imports_alloc(void **imports, uint32_t count);

/**
 * Builds a new `svm_import` (returned via `import` function parameter).
 * New built `svm_import_t` is pushed into `imports`
 *
 * # Example
 *
 * ```rust
 * use svm_runtime_c_api::*;
 *
 * use svm_ffi::{svm_env_t, svm_func_callback_t, svm_byte_array};
 * use svm_types::{WasmType, Type};
 *
 * unsafe extern "C" fn host_func(
 *   env:     *mut svm_env_t,
 *   args:    *const svm_byte_array,
 *   results: *mut svm_byte_array
 * ) -> *mut svm_byte_array {
 *   // ...
 *   return std::ptr::null_mut()
 * }
 *
 * #[repr(C)]
 * struct function_id(u32);
 *
 * // allocate one import
 * let mut imports = testing::imports_alloc(1);
 *
 * let namespace_ty = Type::Str("import ns");
 * let name_ty = Type::Str("import name");
 * let params_ty = Type::Str("import params");
 * let returns_ty = Type::Str("import returns");
 * let host_env_ty = Type::Str("host env");
 *
 * let namespace: svm_byte_array = (namespace_ty, String::from("env")).into();
 * let import_name: svm_byte_array = (name_ty, String::from("foo")).into();
 * let params: svm_byte_array = (params_ty, Vec::<WasmType>::new()).into();
 * let returns: svm_byte_array = (returns_ty,Vec::<WasmType>::new()).into();
 * let mut error = svm_byte_array::default();
 *
 * let host_env = svm_ffi::into_raw(host_env_ty, function_id(0));
 *
 * let res = unsafe {
 *   svm_import_func_new(
 *     imports,
 *     namespace.clone(),
 *     import_name.clone(),
 *     host_func,
 *     host_env,
 *     params.clone(),
 *     returns.clone(),
 *     &mut error)
 * };
 * assert!(res.is_ok());
 * ```
 *
 */
svm_result_t svm_import_func_new(void *imports,
                                 svm_byte_array namespace_,
                                 svm_byte_array import_name,
                                 svm_func_callback_t func,
                                 const void *host_env,
                                 svm_byte_array params,
                                 svm_byte_array returns,
                                 svm_byte_array *error);

/**
 * Creates a new in-memory key-value client.
 * Returns a raw pointer to allocated kv-store via input parameter `kv`.
 *
 * # Example
 *
 * ```rust
 * use svm_runtime_c_api::*;
 *
 * let mut kv = std::ptr::null_mut();
 * let res = unsafe { svm_memory_state_kv_create(&mut kv) };
 * assert!(res.is_ok());
 * ```
 *
 */
svm_result_t svm_memory_state_kv_create(void **kv);

/**
 * Creates a new FFI key-value client.
 * Returns a raw pointer to allocated kv-store via input parameter `kv`.
 *
 * # Example
 *
 * ```rust
 * use svm_runtime_c_api::*;
 *
 * unsafe extern "C" fn get(key_ptr: *const u8, key_len: u32, value_ptr: *mut u8, value_len: *mut u32) {}
 * unsafe extern "C" fn set(key_ptr: *const u8, key_len: u32, value_ptr: *const u8, value_len: u32) {}
 * unsafe extern "C" fn discard() {}
 * unsafe extern "C" fn checkpoint(state: *mut u8) {}
 * unsafe extern "C" fn head(state: *mut u8) {}
 *
 * let mut kv = std::ptr::null_mut();
 * let res = unsafe {
 *   svm_ffi_state_kv_create(
 *     &mut kv,
 *     get,
 *     set,
 *     discard,
 *     checkpoint,
 *     head)
 * };
 * assert!(res.is_ok());
 * ```
 *
 */
svm_result_t svm_ffi_state_kv_create(void **state_kv,
                                     void (*get_fn)(const uint8_t*, uint32_t, uint8_t*, uint32_t*),
                                     void (*set_fn)(const uint8_t*, uint32_t, const uint8_t*, uint32_t),
                                     void (*discard_fn)(void),
                                     void (*checkpoint_fn)(uint8_t*),
                                     void (*head_fn)(uint8_t*));

/**
 * Frees an in-memory key-value.
 *
 * # Example
 *
 * ```rust
 * use svm_runtime_c_api::*;
 *
 * let mut kv = std::ptr::null_mut();
 * let res = unsafe { svm_memory_state_kv_create(&mut kv) };
 * assert!(res.is_ok());
 *
 * let res = unsafe { svm_state_kv_destroy(kv) };
 * assert!(res.is_ok());
 * ```
 *
 */
svm_result_t svm_state_kv_destroy(void *kv);

/**
 * Creates a new SVM Runtime instance baced-by an in-memory KV.
 * Returns it via the `runtime` parameter.
 *
 * # Example
 *
 * ```rust
 * use svm_runtime_c_api::*;
 *
 * use svm_ffi::svm_byte_array;
 *
 * let mut runtime = std::ptr::null_mut();
 * let mut imports = testing::imports_alloc(0);
 *
 * let mut kv = std::ptr::null_mut();
 * let res = unsafe { svm_memory_state_kv_create(&mut kv) };
 * assert!(res.is_ok());
 *
 * let mut error = svm_byte_array::default();
 * let res = unsafe { svm_memory_runtime_create(&mut runtime, kv, imports, &mut error) };
 * assert!(res.is_ok());
 * ```
 *
 */
svm_result_t svm_memory_runtime_create(void **runtime,
                                       void *state_kv,
                                       void *imports,
                                       svm_byte_array *_error);

/**
 * Creates a new SVM Runtime instance.
 * Returns it via the `runtime` parameter.
 *
 * # Example
 *
 * ```rust, no_run
 * use svm_runtime_c_api::*;
 *
 * use svm_types::Type;
 * use svm_ffi::svm_byte_array;
 *
 * let mut runtime = std::ptr::null_mut();
 *
 * let ty = Type::Str("path");
 * let path = String::from("path goes here");
 * let path: svm_byte_array = (ty, path).into();
 * let mut imports = testing::imports_alloc(0);
 * let mut error = svm_byte_array::default();
 *
 * let res = unsafe { svm_runtime_create(&mut runtime, path, imports, &mut error) };
 * assert!(res.is_ok());
 * ```
 *
 */
svm_result_t svm_runtime_create(void **runtime,
                                svm_byte_array kv_path,
                                void *imports,
                                svm_byte_array *error);

/**
 * Deploys a new app-template
 *
 * # Example
 *
 * ```rust, no_run
 * use svm_runtime_c_api::*;
 *
 * use svm_ffi::svm_byte_array;
 * use svm_types::{Address, Type};
 *
 * // allocate imports
 * let mut imports = testing::imports_alloc(0);
 *
 * // create runtime
 * let mut state_kv = std::ptr::null_mut();
 * let res = unsafe { svm_memory_state_kv_create(&mut state_kv) };
 * assert!(res.is_ok());
 *
 * let mut runtime = std::ptr::null_mut();
 * let mut error = svm_byte_array::default();
 * let res = unsafe { svm_memory_runtime_create(&mut runtime, state_kv, imports, &mut error) };
 * assert!(res.is_ok());
 *
 * // deploy template
 * let mut receipt = svm_byte_array::default();
 * let ty = Type::Str("author");
 * let author: svm_byte_array = (ty, Address::of("@author")).into();
 * let template_bytes = svm_byte_array::default();
 * let gas_metering = false;
 * let gas_limit = 0;
 *
 * let res = unsafe {
 *   svm_deploy_template(
 *     &mut receipt,
 *     runtime,
 *     template_bytes,
 *     author,
 *     gas_metering,
 *     gas_limit,
 *     &mut error)
 * };
 *
 * assert!(res.is_ok());
 * ```
 *
 */
svm_result_t svm_deploy_template(svm_byte_array *receipt,
                                 void *runtime,
                                 svm_byte_array bytes,
                                 svm_byte_array author,
                                 bool gas_metering,
                                 uint64_t gas_limit,
                                 svm_byte_array *error);

/**
 * Spawns a new App.
 *
 * # Example
 *
 * ```rust, no_run
 * use svm_runtime_c_api::*;
 *
 * use svm_ffi::svm_byte_array;
 * use svm_types::{Address, Type};
 *
 * // allocate imports
 * let mut imports = testing::imports_alloc(0);
 *
 * // create runtime
 *
 * let mut state_kv = std::ptr::null_mut();
 * let res = unsafe { svm_memory_state_kv_create(&mut state_kv) };
 * assert!(res.is_ok());
 *
 * let mut runtime = std::ptr::null_mut();
 * let mut error = svm_byte_array::default();
 *
 * let res = unsafe { svm_memory_runtime_create(&mut runtime, state_kv, imports, &mut error) };
 * assert!(res.is_ok());
 *
 * let mut app_receipt = svm_byte_array::default();
 * let mut init_state = svm_byte_array::default();
 *
 * let spawner_ty = Type::Str("spawner");
 * let spawner: svm_byte_array = (spawner_ty, Address::of("@spawner")).into();
 * let app_bytes = svm_byte_array::default();
 * let gas_metering = false;
 * let gas_limit = 0;
 *
 * let _res = unsafe {
 *   svm_spawn_app(
 *     &mut app_receipt,
 *     runtime,
 *     app_bytes,
 *     spawner,
 *     gas_metering,
 *     gas_limit,
 *     &mut error)
 * };
 * ```
 *
 */
svm_result_t svm_spawn_app(svm_byte_array *receipt,
                           void *runtime,
                           svm_byte_array bytes,
                           svm_byte_array spawner,
                           bool gas_metering,
                           uint64_t gas_limit,
                           svm_byte_array *error);

/**
 * Triggers an app-transaction execution of an already deployed app.
 * Returns the receipt of the execution via the `receipt` parameter.
 *
 * # Example
 *
 * ```rust, no_run
 * use std::ffi::c_void;
 *
 * use svm_runtime_c_api::*;
 *
 * use svm_types::{State, Address, Type};
 * use svm_ffi::svm_byte_array;
 *
 * // allocate imports
 * let mut imports = testing::imports_alloc(0);
 *
 * // create runtime
 *
 * let mut state_kv = std::ptr::null_mut();
 * let res = unsafe { svm_memory_state_kv_create(&mut state_kv) };
 * assert!(res.is_ok());
 *
 * let mut runtime = std::ptr::null_mut();
 * let mut error = svm_byte_array::default();
 *
 * let res = unsafe { svm_memory_runtime_create(&mut runtime, state_kv, imports, &mut error) };
 * assert!(res.is_ok());
 *
 * let mut exec_receipt = svm_byte_array::default();
 * let bytes = svm_byte_array::default();
 * let ty = Type::of::<State>();
 * let state = (ty, State::empty()).into();
 * let gas_metering = false;
 * let gas_limit = 0;
 *
 * let _res = unsafe {
 *   svm_exec_app(
 *     &mut exec_receipt,
 *     runtime,
 *     bytes,
 *     state,
 *     gas_metering,
 *     gas_limit,
 *     &mut error)
 * };
 * ```
 *
 */
svm_result_t svm_exec_app(svm_byte_array *receipt,
                          void *runtime,
                          svm_byte_array bytes,
                          svm_byte_array state,
                          bool gas_metering,
                          uint64_t gas_limit,
                          svm_byte_array *error);

int32_t svm_total_live_resources(void);

void *svm_resource_iter_new(void);

void svm_resource_iter_destroy(void *iter);

svm_resource_t *svm_resource_iter_next(void *iter);

void svm_resource_destroy(svm_resource_t *resource);

svm_byte_array *svm_resource_type_name_resolve(uintptr_t ty);

void svm_resource_type_name_destroy(svm_byte_array *ptr);

/**
 * Destroys the Runtime and its associated resources.
 *
 * # Example
 *
 * ```rust, no_run
 * use svm_runtime_c_api::*;
 *
 * use svm_types::Address;
 * use svm_ffi::svm_byte_array;
 *
 * // allocate imports
 * let mut imports = testing::imports_alloc(0);
 *
 * // create runtime
 *
 * let mut state_kv = std::ptr::null_mut();
 * let res = unsafe { svm_memory_state_kv_create(&mut state_kv) };
 * assert!(res.is_ok());
 *
 * let mut runtime = std::ptr::null_mut();
 * let mut error = svm_byte_array::default();
 * let res = unsafe { svm_memory_runtime_create(&mut runtime, state_kv, imports, &mut error) };
 * assert!(res.is_ok());
 *
 * // destroy runtime
 * unsafe { svm_runtime_destroy(runtime); }
 * ```
 *
 */
void svm_runtime_destroy(void *runtime);

/**
 * Frees allocated imports resources.
 *
 * # Example
 *
 * ```rust
 * use svm_runtime_c_api::*;
 *
 * // allocate imports
 * let count = 0;
 * let mut imports = std::ptr::null_mut();
 * let _res = unsafe { svm_imports_alloc(&mut imports, count) };
 *
 * // destroy imports
 * unsafe { svm_imports_destroy(imports); }
 * ```
 *
 */
void svm_imports_destroy(void *imports);

/**
 * Frees `svm_byte_array`
 *
 * # Example
 *
 * ```rust
 * use svm_runtime_c_api::*;
 *
 * use svm_ffi::svm_byte_array;
 *
 * let bytes = svm_byte_array::default();
 * unsafe { svm_byte_array_destroy(bytes); }
 * ```
 *
 */
void svm_byte_array_destroy(svm_byte_array bytes);

svm_byte_array *svm_wasm_error_create(svm_byte_array msg);

/**
 * Given a raw `deploy-template` transaction (the `bytes` parameter),
 * if it's valid (i.e: passes the `svm_validate_template`), returns `SVM_SUCCESS` and the estimated gas that will be required
 * in order to execute the transaction (via the `estimate` parameter).
 * # Panics
 *
 * Panics when `bytes` input is not a valid `deploy-template` raw transaction.
 * Having `bytes` a valid raw input doesn't necessarily imply that `svm_validate_template` passes.
 *
 */
svm_result_t svm_estimate_deploy_template(uint64_t *estimation,
                                          void *runtime,
                                          svm_byte_array bytes,
                                          svm_byte_array *error);

/**
 * Given a raw `spawn-app` transaction (the `bytes` parameter),
 * if it's valid (i.e: passes the `svm_validate_app`), returns `SVM_SUCCESS` and the estimated gas that will be required
 * in order to execute the transaction (via the `estimate` parameter).
 *
 * # Panics
 *
 * Panics when `bytes` input is not a valid `spawn-app` raw transaction.
 * Having `bytes` a valid raw input doesn't necessarily imply that `svm_validate_app` passes.
 *
 */
svm_result_t svm_estimate_spawn_app(uint64_t *estimation,
                                    void *runtime,
                                    svm_byte_array bytes,
                                    svm_byte_array *error);

/**
 * Given a raw `exec-app` transaction (the `bytes` parameter),
 * if it's valid (i.e: passes the `svm_validate_tx`), returns `SVM_SUCCESS` and the estimated gas that will be required
 * in order to execute the transaction (via the `estimate` parameter).
 *
 * # Panics
 *
 * Panics when `bytes` input is not a valid `exec-app` raw transaction.
 * Having `bytes` a valid raw input doesn't necessarily imply that `svm_validate_tx` passes.
 *
 */
svm_result_t svm_estimate_exec_app(uint64_t *estimation,
                                   void *runtime,
                                   svm_byte_array bytes,
                                   svm_byte_array *error);

/**
 * Constructs a new raw `app_template` transaction.
 *
 */
svm_result_t svm_encode_app_template(svm_byte_array *app_template,
                                     uint32_t version,
                                     svm_byte_array name,
                                     svm_byte_array code,
                                     svm_byte_array data,
                                     svm_byte_array *error);

/**
 * Constructs a new raw `spawn_app` transaction.
 *
 */
svm_result_t svm_encode_spawn_app(svm_byte_array *spawn_app,
                                  uint32_t version,
                                  svm_byte_array template_addr,
                                  svm_byte_array name,
                                  svm_byte_array ctor_name,
                                  svm_byte_array calldata,
                                  svm_byte_array *error);

/**
 * Constructs a new raw `app_tx` transaction.
 *
 */
svm_result_t svm_encode_app_tx(svm_byte_array *app_tx,
                               uint32_t version,
                               svm_byte_array app_addr,
                               svm_byte_array func_name,
                               svm_byte_array calldata,
                               svm_byte_array *error);

#endif /* SVM_H */
