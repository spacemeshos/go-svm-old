cargo +nightly build --release --target wasm32-unknown-unknown

mv ./target/wasm32-unknown-unknown/release/go_svm_examples_counter.wasm ./counter.wasm
