#ifndef SVM_H
#define SVM_H

#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

/**
 * FFI representation for function result type
 */
typedef enum svm_result_t {
  SVM_SUCCESS = 0,
  SVM_FAILURE = 1,
} svm_result_t;

/**
 * FFI representation for a byte-array
 *
 * # Examples
 *
 * ```rust
 * use svm_runtime_ffi::svm_byte_array;
 * use svm_types::Type;
 *
 * use std::convert::TryFrom;
 * use std::string::FromUtf8Error;
 *
 * let ty = Type::Str("test string");
 *
 * let s1 = "Hello World!".to_string();
 * let bytes: svm_byte_array = (ty, s1).into();
 *
 * let s2 = String::try_from(bytes).unwrap();
 * assert_eq!(s2, "Hello World!".to_string());
 * ```
 *
 */
typedef struct svm_byte_array {
  const uint8_t *bytes;
  uint32_t length;
  uint32_t capacity;
  uintptr_t type_id;
} svm_byte_array;

/**
 * Represents a manually-allocated resource.
 */
typedef struct svm_resource_t {
  /**
   * Type interned value
   */
  uintptr_t type_id;
  /**
   * `#resources` of that type
   */
  int32_t count;
} svm_resource_t;

/**
 *
 * Start of the Public C-API
 *
 * * Each method is annotated with `#[no_mangle]`
 * * Each method has `unsafe extern "C"` before `fn`
 *
 * See `build.rs` for using `cbindgen` to generate `svm.h`
 *
 *
 * Creates a new SVM Runtime instance backed-by an in-memory KV.
 *
 * Returns it the created Runtime via the `runtime` parameter.
 *
 * # Examples
 *
 * ```rust
 * use svm_runtime_ffi::*;
 *
 * let mut runtime = std::ptr::null_mut();
 *
 * let mut error = svm_byte_array::default();
 * let res = unsafe { svm_memory_runtime_create(&mut runtime, &mut error) };
 * assert!(res.is_ok());
 * ```
 *
 */
enum svm_result_t svm_memory_runtime_create(void **runtime, struct svm_byte_array *error);

/**
 * Destroys the Runtime and its associated resources.
 *
 * # Examples
 *
 * ```rust, no_run
 * use svm_runtime_ffi::*;
 *
 * let mut runtime = std::ptr::null_mut();
 * let mut error = svm_byte_array::default();
 * let res = unsafe { svm_memory_runtime_create(&mut runtime, &mut error) };
 * assert!(res.is_ok());
 *
 * // Destroys the Runtime
 * unsafe { svm_runtime_destroy(runtime); }
 * ```
 *
 */
void svm_runtime_destroy(void *runtime);

/**
 * Allocates `svm_byte_array` to be used later for passing a binary [`Envelope`].
 *
 * The number of allocated bytes is equal to [`Envelope`]'s
 * [`Codec::fixed_size()`].
 */
struct svm_byte_array svm_envelope_alloc(void);

/**
 * Allocates `svm_byte_array` of `size` bytes, meant to be used for passing a
 * binary message.
 */
struct svm_byte_array svm_message_alloc(uint32_t size);

/**
 * Allocates `svm_byte_array` to be used later for passing a binary [`Context`].
 *
 * The number of allocated bytes is equal to [`Context`]'s
 * [`Codec::fixed_size()`].
 */
struct svm_byte_array svm_context_alloc(void);

/**
 * Validates syntactically a binary `Deploy Template` transaction.
 *
 * Should be called while the transaction is in the `mempool` of the Host.
 * In case the transaction isn't valid - the transaction should be discarded.
 *
 *
 * # Examples
 *
 * ```rust, no_run
 * use svm_runtime_ffi::*;
 *
 * let mut runtime = std::ptr::null_mut();
 * let mut error = svm_byte_array::default();
 *
 * let res = unsafe { svm_memory_runtime_create(&mut runtime, &mut error) };
 * assert!(res.is_ok());
 *
 * let message = svm_byte_array::default();
 *
 * let _res = unsafe { svm_validate_deploy(runtime, message, &mut error) };
 * ```
 *
 */
enum svm_result_t svm_validate_deploy(void *runtime,
                                      struct svm_byte_array message,
                                      struct svm_byte_array *error);

/**
 * Validates syntactically a binary `Spawn Account` transaction.
 *
 * Should be called while the transaction is in the `mempool` of the Host.
 * In case the transaction isn't valid - the transaction should be discarded.
 *
 *
 * # Examples
 *
 * ```rust, no_run
 * use svm_runtime_ffi::*;
 *
 * let mut runtime = std::ptr::null_mut();
 * let mut error = svm_byte_array::default();
 *
 * let res = unsafe { svm_memory_runtime_create(&mut runtime, &mut error) };
 * assert!(res.is_ok());
 *
 * let message = svm_byte_array::default();
 * let _res = unsafe { svm_validate_spawn(runtime, message, &mut error) };
 * ```
 *
 */
enum svm_result_t svm_validate_spawn(void *runtime,
                                     struct svm_byte_array message,
                                     struct svm_byte_array *error);

/**
 * Validates syntactically a binary `Call Account` transaction.
 *
 * # Examples
 *
 * ```rust, no_run
 * use svm_runtime_ffi::*;
 *
 * let mut runtime = std::ptr::null_mut();
 * let mut error = svm_byte_array::default();
 *
 * let res = unsafe { svm_memory_runtime_create(&mut runtime, &mut error) };
 * assert!(res.is_ok());
 *
 * let message = svm_byte_array::default();
 * let _res = unsafe { svm_validate_call(runtime, message, &mut error) };
 * ```
 *
 */
enum svm_result_t svm_validate_call(void *runtime,
                                    struct svm_byte_array message,
                                    struct svm_byte_array *error);

/**
 * Deploys a `Template`
 *
 * # Examples
 *
 * ```rust, no_run
 * use svm_runtime_ffi::*;
 *
 * let mut runtime = std::ptr::null_mut();
 * let mut error = svm_byte_array::default();
 * let res = unsafe { svm_memory_runtime_create(&mut runtime, &mut error) };
 * assert!(res.is_ok());
 *
 * let mut receipt = svm_byte_array::default();
 * let envelope = svm_byte_array::default();
 * let message = svm_byte_array::default();
 * let context = svm_byte_array::default();
 *
 * let res = unsafe {
 *   svm_deploy(
 *     &mut receipt,
 *     runtime,
 *     envelope,
 *     message,
 *     context,
 *     &mut error)
 * };
 *
 * assert!(res.is_ok());
 * ```
 *
 */
enum svm_result_t svm_deploy(struct svm_byte_array *receipt,
                             void *runtime,
                             struct svm_byte_array envelope,
                             struct svm_byte_array message,
                             struct svm_byte_array context,
                             struct svm_byte_array *error);

/**
 * Spawns a new `Account`.
 *
 * # Examples
 *
 * ```rust, no_run
 * use svm_runtime_ffi::*;
 *
 * let mut runtime = std::ptr::null_mut();
 * let mut error = svm_byte_array::default();
 *
 * let res = unsafe { svm_memory_runtime_create(&mut runtime, &mut error) };
 * assert!(res.is_ok());
 *
 * let mut receipt = svm_byte_array::default();
 * let mut init_state = svm_byte_array::default();
 *
 * let envelope = svm_byte_array::default();
 * let message = svm_byte_array::default();
 * let context = svm_byte_array::default();
 *
 * let _res = unsafe {
 *   svm_spawn(
 *     &mut receipt,
 *     runtime,
 *     envelope,
 *     message,
 *     context,
 *     &mut error)
 * };
 * ```
 *
 */
enum svm_result_t svm_spawn(struct svm_byte_array *receipt,
                            void *runtime,
                            struct svm_byte_array envelope,
                            struct svm_byte_array message,
                            struct svm_byte_array context,
                            struct svm_byte_array *error);

/**
 * Calls `verify` on an Account.
 * The inputs `envelope`, `message` and `context` should be the same ones
 * passed later to `svm_call`.(in case the `verify` succeeds).
 *
 * Returns the Receipt of the execution via the `receipt` parameter.
 *
 * # Examples
 *
 * ```rust, no_run
 * use std::ffi::c_void;
 *
 * use svm_runtime_ffi::*;
 *
 * let mut runtime = std::ptr::null_mut();
 * let mut error = svm_byte_array::default();
 *
 * let res = unsafe { svm_memory_runtime_create(&mut runtime, &mut error) };
 * assert!(res.is_ok());
 *
 * let mut receipt = svm_byte_array::default();
 * let envelope = svm_byte_array::default();
 * let message = svm_byte_array::default();
 * let context = svm_byte_array::default();
 *
 * let _res = unsafe {
 *   svm_verify(
 *     &mut receipt,
 *     runtime,
 *     envelope,
 *     message,
 *     context,
 *     &mut error)
 * };
 * ```
 *
 */
enum svm_result_t svm_verify(struct svm_byte_array *receipt,
                             void *runtime,
                             struct svm_byte_array envelope,
                             struct svm_byte_array message,
                             struct svm_byte_array context,
                             struct svm_byte_array *error);

/**
 * `Call Account` transaction.
 * Returns the Receipt of the execution via the `receipt` parameter.
 *
 * # Examples
 *
 * ```rust, no_run
 * use std::ffi::c_void;
 *
 * use svm_runtime_ffi::*;
 *
 * let mut runtime = std::ptr::null_mut();
 * let mut error = svm_byte_array::default();
 *
 * let res = unsafe { svm_memory_runtime_create(&mut runtime, &mut error) };
 * assert!(res.is_ok());
 *
 * let mut receipt = svm_byte_array::default();
 * let envelope = svm_byte_array::default();
 * let message = svm_byte_array::default();
 * let context = svm_byte_array::default();
 *
 * let _res = unsafe {
 *   svm_call(
 *     &mut receipt,
 *     runtime,
 *     envelope,
 *     message,
 *     context,
 *     &mut error)
 * };
 * ```
 *
 */
enum svm_result_t svm_call(struct svm_byte_array *receipt,
                           void *runtime,
                           struct svm_byte_array envelope,
                           struct svm_byte_array message,
                           struct svm_byte_array context,
                           struct svm_byte_array *error);

/**
 * Returns the total live manually-managed resources.
 */
int32_t svm_total_live_resources(void);

/**
 * Initializes a new iterator over the manually-managed resources
 */
void *svm_resource_iter_new(void);

/**
 * Destroys the manually-managed resources iterator
 */
void svm_resource_iter_destroy(void *iter);

/**
 * Returns the next manually-managed resource.
 * If there is no resource to return, returns `NULL`
 */
struct svm_resource_t *svm_resource_iter_next(void *iter);

/**
 * Destroy the resource
 */
void svm_resource_destroy(struct svm_resource_t *resource);

/**
 * Given a type in an interned form, returns its textual name
 */
struct svm_byte_array *svm_resource_type_name_resolve(uintptr_t ty);

/**
 * Destroys a resource holding a type textual name
 */
void svm_resource_type_name_destroy(struct svm_byte_array *ptr);

/**
 * Frees `svm_byte_array`
 *
 * # Examples
 *
 * ```rust
 * use svm_runtime_ffi::*;
 *
 * let bytes = svm_byte_array::default();
 * unsafe { svm_byte_array_destroy(bytes); }
 * ```
 *
 */
void svm_byte_array_destroy(struct svm_byte_array bytes);

#endif /* SVM_H */
