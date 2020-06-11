#ifndef SVM_H
#define SVM_H

#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

/**
 * FFI representation for function result type
 */
typedef enum {
  SVM_SUCCESS = 0,
  SVM_FAILURE = 1,
} svm_result_t;

/**
 * FFI representation for a byte-array
 *
 * # Example
 *
 * ```rust
 * use std::{convert::TryFrom, string::FromUtf8Error};
 * use svm_runtime_c_api::svm_byte_array;
 *
 * let s1 = "Hello World!".to_string();
 * let ptr = s1.as_ptr();
 * let length = s1.len() as u32;
 * let bytes = svm_byte_array { bytes: ptr, length };
 *
 * let s2 = String::try_from(bytes);
 * assert_eq!(s1, s2.unwrap());
 * ```
 *
 */
typedef struct {
  /**
   * Raw pointer to the beginning of array.
   */
  const uint8_t *bytes;
  /**
   * Number of bytes,
   */
  uint32_t length;
} svm_byte_array;

/**
 * Extracts the spawned-app `Address`.
 * When spawning succeeds returns `SVM_SUCCESS` and the `Address` via `app_addr` parameter.
 * Otherise, returns `SVM_FAILURE` and the error message via the `error` parameter.
 *
 * # Panics
 *
 * Panics when `receipt` input is invalid.
 *
 */
svm_result_t svm_app_receipt_addr(svm_byte_array *app_addr,
                                  svm_byte_array receipt,
                                  svm_byte_array *error);

/**
 * Extracts the `gas_used` for spawned-app (including running its constructor).
 * When spawn succeeded returns `SVM_SUCCESS`, returns the amount of gas used via `gas_used` parameter.
 * Othewrise, returns `SVM_FAILURE` and the error message via the `error` parameter.
 *
 * It's up for the Host to decide the gas fee to for failed spawning.
 * (usually the strategy will be to fine with the `gas_limit` of the failed transaction).
 *
 * # Panics
 *
 * Panics when `receipt` input is invalid.
 *
 */
svm_result_t svm_app_receipt_gas(uint64_t *gas_used,
                                 svm_byte_array receipt,
                                 svm_byte_array *error);

/**
 * Extracts the spawned-app constructor `returns`.
 * The `returns` are encoded as `svm_byte_array`.
 * More info regarding the encoding in `byte_array.rs`.
 *
 * If it succeeded, returns `SVM_SUCCESS`,
 * Otherwise returns `SVM_FAILURE` and the error message via `error` parameter.
 *
 * # Panics
 *
 * Panics when `receipt` input is invalid.
 *
 */
svm_result_t svm_app_receipt_returns(svm_byte_array *returns,
                                     svm_byte_array receipt,
                                     svm_byte_array *error);

/**
 * `Exec-App` Receipt helpers
 *  -------------------------------------------------------
 * Extracts the spawned-app initial `State`.
 * When spawning succeeds returns `SVM_SUCCESS` and the initial `State` via `state` parameter.
 * Otherise, returns `SVM_FAILURE` and the error message via the `error` parameter.
 *
 * # Panics
 *
 * Panics when `receipt` input is invalid.
 *
 */
svm_result_t svm_app_receipt_state(svm_byte_array *state,
                                   svm_byte_array receipt,
                                   svm_byte_array *error);

/**
 * `Spawn-App` Receipt helpers
 *  -------------------------------------------------------
 * Extracts whether the `spawn-app` transaction succeeded.
 * If it succeeded, returns `SVM_SUCCESS`,
 * Otherwise returns `SVM_FAILURE` and the error message via `error` parameter.
 *
 * # Panics
 *
 * Panics when `receipt` input is invalid.
 *
 */
svm_result_t svm_app_receipt_status(svm_byte_array receipt, svm_byte_array *error);

/**
 * Frees `svm_byte_array`
 *
 * # Example
 *
 * ```rust
 * use svm_runtime_c_api::*;
 *
 * let bytes = svm_byte_array::default();
 * unsafe { svm_byte_array_destroy(bytes); }
 * ```
 *
 */
void svm_byte_array_destroy(svm_byte_array bytes);

/**
 * Deploys a new app-template
 *
 * # Example
 *
 * ```rust, no_run
 * use svm_runtime_c_api::*;
 * use svm_common::Address;
 *
 * let mut host = std::ptr::null_mut();
 *
 * // allocate imports
 * let mut imports = testing::imports_alloc(0);
 *
 * // create runtime
 * let mut kv = std::ptr::null_mut();
 * let res = unsafe { svm_memory_kv_create(&mut kv) };
 * assert!(res.is_ok());
 *
 * let mut raw_kv = std::ptr::null_mut();
 * let res = unsafe { svm_memory_kv_create2(&mut raw_kv) };
 * assert!(res.is_ok());
 *
 * let mut runtime = std::ptr::null_mut();
 * let mut error = svm_byte_array::default();
 * let res = unsafe { svm_memory_runtime_create(&mut runtime, kv, raw_kv, host, imports, &mut error) };
 * assert!(res.is_ok());
 *
 * // deploy template
 * let mut receipt = svm_byte_array::default();
 * let author: svm_byte_array = Address::of("@author").into();
 * let host_ctx = svm_byte_array::default();
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
 *     host_ctx,
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
                                 svm_byte_array host_ctx,
                                 bool gas_metering,
                                 uint64_t gas_limit,
                                 svm_byte_array *error);

/**
 * Constructs a new raw `app_template` transaction.
 *
 */
svm_result_t svm_encode_app_template(svm_byte_array *app_template,
                                     uint32_t version,
                                     svm_byte_array name,
                                     uint16_t page_count,
                                     svm_byte_array code,
                                     svm_byte_array *error);

/**
 * Constructs a new raw `app_tx` transaction.
 *
 * The `func_args` is `svm_byte_array` representing a slice of `WasmValue`.
 * More info regarding the encoding in `byte_array.rs`.
 *
 */
svm_result_t svm_encode_app_tx(svm_byte_array *app_tx,
                               uint32_t version,
                               svm_byte_array app_addr,
                               uint16_t func_idx,
                               svm_byte_array func_buf,
                               svm_byte_array func_args,
                               svm_byte_array *error);

/**
 * Constructs a new raw `spawn_app` transaction.
 *
 * The `ctor_args` is `svm_byte_array` representing a slice of `WasmValue`.
 * More info regarding the encoding in `byte_array.rs`.
 *
 */
svm_result_t svm_encode_spawn_app(svm_byte_array *spawn_app,
                                  uint32_t version,
                                  svm_byte_array template_addr,
                                  uint16_t ctor_idx,
                                  svm_byte_array ctor_buf,
                                  svm_byte_array ctor_args,
                                  svm_byte_array *error);

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
 * Triggers an app-transaction execution of an already deployed app.
 * Returns the receipt of the execution via the `receipt` parameter.
 *
 * # Example
 *
 * ```rust, no_run
 * use std::ffi::c_void;
 *
 * use svm_runtime_c_api::*;
 * use svm_common::{State, Address};
 *
 * let mut host = std::ptr::null_mut();
 *
 * // allocate imports
 * let mut imports = testing::imports_alloc(0);
 *
 * // create runtime
 * let mut kv = std::ptr::null_mut();
 * let res = unsafe { svm_memory_kv_create(&mut kv) };
 * assert!(res.is_ok());
 *
 * let mut raw_kv = std::ptr::null_mut();
 * let res = unsafe { svm_memory_kv_create2(&mut raw_kv) };
 * assert!(res.is_ok());
 *
 * let mut runtime = std::ptr::null_mut();
 * let mut error = svm_byte_array::default();
 *
 * let res = unsafe { svm_memory_runtime_create(&mut runtime, kv, raw_kv, host, imports, &mut error) };
 * assert!(res.is_ok());
 *
 * let mut exec_receipt = svm_byte_array::default();
 * let tx_bytes = svm_byte_array::default();
 * let state = State::empty().into();
 * let host_ctx = svm_byte_array::default();
 * let gas_metering = false;
 * let gas_limit = 0;
 *
 * let _res = unsafe {
 *   svm_exec_app(
 *     &mut exec_receipt,
 *     runtime,
 *     tx_bytes,
 *     state,
 *     host_ctx,
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
                          svm_byte_array host_ctx,
                          bool gas_metering,
                          uint64_t gas_limit,
                          svm_byte_array *error);

/**
 * Extracts the executed transaction `gas_used`.
 * When transaction succeeded returns `SVM_SUCCESS`, returns the amount of gas used via `gas_used` parameter.
 * Othewrise, returns `SVM_FAILURE` and the error message via the `error` parameter.
 *
 * It's up for the Host to decide the gas fee to for failed transactions.
 * (usually the strategy will be to fine with the `gas_limit` of the failed transaction).
 *
 * # Panics
 *
 * Panics when `receipt` input is invalid.
 *
 */
svm_result_t svm_exec_receipt_gas(uint64_t *gas_used,
                                  svm_byte_array receipt,
                                  svm_byte_array *error);

/**
 * Extracts the `Exec App` `returns`.
 * The `returns` are encoded as `svm_byte_array`.
 * More info regarding the encoding in `byte_array.rs`.
 *
 * If it succeeded, returns `SVM_SUCCESS`,
 * Otherwise returns `SVM_FAILURE` and the error message via `error` parameter.
 *
 * # Panics
 *
 * Panics when `receipt` input is invalid.
 *
 */
svm_result_t svm_exec_receipt_returns(svm_byte_array *returns,
                                      svm_byte_array receipt,
                                      svm_byte_array *error);

/**
 * Extracts the executed transaction new `State`.
 * When transaction succeeded returns `SVM_SUCCESS` and the new `State` via `state` parameter.
 * Othewrise, returns `SVM_FAILURE` and the error message via the `error` parameter.
 *
 * # Panics
 *
 * Panics when `receipt` input is invalid.
 *
 */
svm_result_t svm_exec_receipt_state(svm_byte_array *state,
                                    svm_byte_array receipt,
                                    svm_byte_array *error);

/**
 * Extracts whether the `exec-app` transaction succeeded.
 * If it succeeded, returns `SVM_SUCCESS`,
 * Otherwise returns `SVM_FAILURE` and the error message via `error` parameter.
 *
 * # Panics
 *
 * Panics when `receipt` input is invalid.
 *
 */
svm_result_t svm_exec_receipt_status(svm_byte_array receipt, svm_byte_array *error);

/**
 * Builds a new `svm_import` (returned via `import` function parameter).
 * New built `svm_import_t` is pushed into `imports`
 *
 * # Example
 *
 * ```rust
 * use svm_app::types::WasmType;
 * use svm_runtime_c_api::*;
 *
 * fn foo() {
 *   // ...
 * }
 *
 * // allocate one imports
 * let mut imports = testing::imports_alloc(1);
 *
 * let module_name = "env".into();
 * let import_name = "foo".into();
 * let params = Vec::<WasmType>::new();
 * let returns = Vec::<WasmType>::new();
 * let func = foo as *const std::ffi::c_void;
 * let mut error = svm_byte_array::default();
 *
 * let res = unsafe {
 *   svm_import_func_build(
 *     imports,
 *     module_name,
 *     import_name,
 *     func,
 *     params.into(),
 *     returns.into(),
 *     &mut error)
 * };
 * assert!(res.is_ok());
 * ```
 *
 */
svm_result_t svm_import_func_build(void *imports,
                                   svm_byte_array module_name,
                                   svm_byte_array import_name,
                                   const void *func,
                                   svm_byte_array params,
                                   svm_byte_array returns,
                                   svm_byte_array *error);

/**
 * Allocates space for the host imports.
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
void svm_imports_destroy(const void *imports);

/**
 * Returns a raw pointer to `the host` extracted from a raw pointer to `wasmer` context.
 */
void *svm_instance_context_host_get(void *ctx);

/**
 * Creates a new in-memory `MemKVStore`.
 * Returns a raw pointer to allocated kv-store via input parameter `raw_kv`.
 */
svm_result_t svm_memory_kv_create(void **kv);

/**
 * Creates a new in-memory key-value client.
 * Returns a raw pointer to allocated kv-store via input parameter `raw_kv`.
 */
svm_result_t svm_memory_kv_create2(void **kv);

/**
 * Creates a new SVM Runtime instance baced-by an in-memory KV.
 * Returns it via the `runtime` parameter.
 *
 * # Example
 *
 * ```rust
 * use svm_runtime_c_api::*;
 *
 * let mut runtime = std::ptr::null_mut();
 * let host = std::ptr::null_mut();
 * let mut imports = testing::imports_alloc(0);
 *
 * let mut kv = std::ptr::null_mut();
 * let res = unsafe { svm_memory_kv_create(&mut kv) };
 * assert!(res.is_ok());
 *
 * let mut raw_kv = std::ptr::null_mut();
 * let res = unsafe { svm_memory_kv_create2(&mut raw_kv) };
 * assert!(res.is_ok());
 *
 * let mut error = svm_byte_array::default();
 * let res = unsafe { svm_memory_runtime_create(&mut runtime, kv, raw_kv, host, imports, &mut error) };
 * assert!(res.is_ok());
 * ```
 *
 */
svm_result_t svm_memory_runtime_create(void **runtime,
                                       void *kv,
                                       void *raw_kv,
                                       void *host,
                                       const void *imports,
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
 * let mut runtime = std::ptr::null_mut();
 * let path = "path goes here".into();
 * let host = std::ptr::null_mut();
 * let mut imports = testing::imports_alloc(0);
 * let mut error = svm_byte_array::default();
 *
 * let res = unsafe { svm_runtime_create(&mut runtime, path, host, imports, &mut error) };
 * assert!(res.is_ok());
 * ```
 *
 */
svm_result_t svm_runtime_create(void **runtime,
                                svm_byte_array kv_path,
                                void *host,
                                const void *imports,
                                svm_byte_array *error);

/**
 * Destroys the Runtime and its associated resources.
 *
 * # Example
 *
 * ```rust, no_run
 * use svm_runtime_c_api::*;
 * use svm_common::Address;
 *
 * let mut host = std::ptr::null_mut();
 *
 * // allocate imports
 * let mut imports = testing::imports_alloc(0);
 *
 * // create runtime
 * let mut kv = std::ptr::null_mut();
 * let res = unsafe { svm_memory_kv_create(&mut kv) };
 * assert!(res.is_ok());
 *
 * let mut raw_kv = std::ptr::null_mut();
 * let res = unsafe { svm_memory_kv_create2(&mut raw_kv) };
 * assert!(res.is_ok());
 *
 * let mut runtime = std::ptr::null_mut();
 * let mut error = svm_byte_array::default();
 * let res = unsafe { svm_memory_runtime_create(&mut runtime, kv, raw_kv, host, imports, &mut error) };
 * assert!(res.is_ok());
 *
 * // destroy runtime
 * unsafe { svm_runtime_destroy(runtime); }
 * ```
 *
 */
void svm_runtime_destroy(void *runtime);

/**
 * Spawns a new App.
 *
 * # Example
 *
 * ```rust, no_run
 * use svm_runtime_c_api::*;
 * use svm_common::Address;
 *
 * let mut host = std::ptr::null_mut();
 *
 * // allocate imports
 * let mut imports = testing::imports_alloc(0);
 *
 * // create runtime
 * let mut kv = std::ptr::null_mut();
 * let res = unsafe { svm_memory_kv_create(&mut kv) };
 * assert!(res.is_ok());
 *
 * let mut raw_kv = std::ptr::null_mut();
 * let res = unsafe { svm_memory_kv_create2(&mut raw_kv) };
 * assert!(res.is_ok());
 *
 * let mut runtime = std::ptr::null_mut();
 * let mut error = svm_byte_array::default();
 *
 * let res = unsafe { svm_memory_runtime_create(&mut runtime, kv, raw_kv, host, imports, &mut error) };
 * assert!(res.is_ok());
 *
 * let mut app_receipt = svm_byte_array::default();
 * let mut init_state = svm_byte_array::default();
 * let creator = Address::of("@creator").into();
 * let host_ctx = svm_byte_array::default();
 * let app_bytes = svm_byte_array::default();
 * let gas_metering = false;
 * let gas_limit = 0;
 *
 * let _res = unsafe {
 *   svm_spawn_app(
 *     &mut app_receipt,
 *     runtime,
 *     app_bytes,
 *     creator,
 *     host_ctx,
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
                           svm_byte_array creator,
                           svm_byte_array host_ctx,
                           bool gas_metering,
                           uint64_t gas_limit,
                           svm_byte_array *error);

/**
 * Receipts helpers
 * In order to spare the SVM client the implementation of the `Receipt`(s) raw decoding the receipts helpers
 * can fetch one field each. This functionality should be useful for writing tests when using client code that interfaces with
 * SVM FFI interface.
 *
 * Each helper methods returns `svm_result_t`.
 * When `svm_result_t` equals `SVM_SUCCESS` is means that the field extraction succeeded.
 * Otherwise, it signals that the field can't be extracted out of the receipt.
 *
 * For example, if the `svm_deploy_template` failed to deploy the template (it may happen for many reason, one is having invalid wasm code),
 * then calling `svm_template_receipt_addr` should return `SVM_FAILURE` since there is no template `Address` to extract.
 * The error will be returned via the `error` parameter.
 * `Deploy-Template` Receipt helpers
 *  -------------------------------------------------------
 * Extracts the deploy-template `Address` into the `template_addr` parameter. (useful for tests).
 *
 * # Panics
 *
 * Panics when `receipt` input is invalid.
 *
 */
svm_result_t svm_template_receipt_addr(svm_byte_array *template_addr,
                                       svm_byte_array receipt,
                                       svm_byte_array *error);

/**
 * Extracts the `gas_used` for the deploy-template.
 * When deploying succeeded returns `SVM_SUCCESS`, returns the amount of gas used via `gas_used` parameter.
 * Othewrise, returns `SVM_FAILURE` and the error message via the `error` parameter.
 *
 * It's up for the Host to decide the gas fee to for a failed deploy.
 * (usually the strategy will be to fine with the `gas_limit` of the failed transaction).
 *
 * # Panics
 *
 * Panics when `receipt` input is invalid.
 *
 */
svm_result_t svm_template_receipt_gas(uint64_t *gas_used,
                                      svm_byte_array receipt,
                                      svm_byte_array *error);

svm_result_t svm_validate_app(const void *runtime, svm_byte_array bytes, svm_byte_array *error);

svm_result_t svm_validate_template(const void *runtime,
                                   svm_byte_array bytes,
                                   svm_byte_array *error);

/**
 * Parses `exec-app` raw transaction.
 * Returns the `App` address that appears in the transaction.
 *
 * # Example
 *
 * ```rust, no_run
 * use svm_runtime_c_api::*;
 * use svm_common::Address;
 *
 * let mut host = std::ptr::null_mut();
 *
 * // allocate imports
 * let mut imports = testing::imports_alloc(0);
 *
 * // create runtime
 * let mut kv = std::ptr::null_mut();
 * let res = unsafe { svm_memory_kv_create(&mut kv) };
 * assert!(res.is_ok());
 *
 * let mut raw_kv = std::ptr::null_mut();
 * let res = unsafe { svm_memory_kv_create(&mut raw_kv) };
 * assert!(res.is_ok());
 *
 * let mut runtime = std::ptr::null_mut();
 * let mut error = svm_byte_array::default();
 *
 * let res = unsafe { svm_memory_runtime_create(&mut runtime, kv, raw_kv, host, imports, &mut error) };
 * assert!(res.is_ok());
 *
 * let mut app_addr = svm_byte_array::default();
 * let tx_bytes = svm_byte_array::default();
 * let _res = unsafe { svm_validate_tx(&mut app_addr, runtime, tx_bytes, &mut error) };
 * ```
 *
 */
svm_result_t svm_validate_tx(svm_byte_array *app_addr,
                             const void *runtime,
                             svm_byte_array bytes,
                             svm_byte_array *error);

#endif /* SVM_H */
