(module
  (func $get32 (import "svm" "get32") (param i32) (result i32))
  (func $set32 (import "svm" "set32") (param i32 i32))
  (func $env_inc (import "env" "inc") (param i32))
  (func $env_get (import "env" "get") (result i32))

  (memory 1) ;; memory `0` (default) is initialized with one page

  (func (export "storage_inc") (param $val i32)
      ;; push var_id = 0 for later `$set32` usage
      i32.const 0

      ;; read var #0
      i32.const 0  ;; var_id = 0
      call $get32

      ;; calculate var #0 new value
      get_local $val
      i32.add

      ;; store var #0 new value
      call $set32
  )

  (func $storage_get (export "storage_get") (result i32)
      ;; return var #0
      i32.const 0  ;; var_id = 0
      call $get32
  )

  (func (export "host_inc") (param $val i32)
      get_local $val
      call $env_inc
  )

  (func (export "host_get") (result i32)
      call $env_get
  )
)
