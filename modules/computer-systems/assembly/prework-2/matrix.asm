section .text
global index

; rdi: matrix
; rsi: rows
; rdx: cols
; rcx: rindex
; r8: cindex
;;
;; element    = matrix + (int * cols * rindex) + (cindex * int)
;;              where int = 4 

index:  
  imul    rcx, INT_SIZE
  imul    rdx, rcx

  imul    r8, INT_SIZE
  add     rdx, r8

  add     rdi, rdx
  mov     rax, [rdi]
  ret

section .data
INT_SIZE: equ 4
