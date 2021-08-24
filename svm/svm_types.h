#ifndef SVM_TYPES_H
#define SVM_TYPES_H

#include <stdarg.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdlib.h>

/**
 * FFI representation for a byte-array
 *
 * # Example
 *
 * ```rust
 * use std::convert::TryFrom;
 * use std::string::FromUtf8Error;
 *
 * use svm_types::Type;
 * use svm_ffi::svm_byte_array;
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
typedef struct {
  /**
   * Raw pointer to the beginning of array.
   */
  const uint8_t *bytes;
  /**
   * Number of bytes of the data view.
   */
  uint32_t length;
  /**
   * Total number of allocated bytes.
   * It may be unequal and bigger than `length` if the `svm_byte_array` instance is an alias to
   * an instance of a data structure such as `Vec` (which in order to properly get deallocated
   * needs first to be re-constructed using the proper allocated capacity).
   */
  uint32_t capacity;
  /**
   * The `svm_types::Type` associated with the data represented by `bytes`.
   * It's the interned value of the type. (For more info see `tracking::interning.rs`)
   */
  uintptr_t type_id;
} svm_byte_array;

/**
 * The function environment of host functions.
 */
typedef struct {
  /**
   * The SVM's inner environment.
   * (see `Context` at `svm-runtime` crate).
   */
  const void *inner_env;
  /**
   * The host environment.
   */
  const void *host_env;
} svm_env_t;

/**
 * Import function signature.
 */
typedef svm_byte_array *(*svm_func_callback_t)(svm_env_t *env, const svm_byte_array *args, svm_byte_array *results);

/**
 * Represents a manually-allocated resource.
 */
typedef struct {
  /**
   * Type interned value
   */
  uintptr_t type_id;
  /**
   * `#resources` of that type
   */
  int32_t count;
} svm_resource_t;

#endif /* SVM_TYPES_H */
