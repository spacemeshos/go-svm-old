## go-svm

Go bindings for [SVM](https://github.com/spacemeshos/svm)

---

### Project structure

This repository contains both Rust and Go library packages.

The Rust package (`/svm-dep`) is defined as an empty package, with a dependency for the `svm-runtime-c-api` Rust package. Once compiled via Cargo, the `.dylib`/`.so`/`.dll` artifacts (on MacOS, Linux and Windows, respectively), in addition to the header file, can be copied to to the Go package (`/svm`), to be linked via `cgo`. 

To allow direct and seamless import of the Go package, it includes the pre-compiled binaries mentioned above, which will be continuously updated. It *currently* supports only macOS x86_64.
