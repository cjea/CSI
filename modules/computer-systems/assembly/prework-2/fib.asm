; Implement the following C program:
;;
;;     #include <stdio.h>
;; 
;;     int fib(int n) {
;;       if (n <= 1) return n;
;; 
;;       return fib(n-1) + fib(n-2);
;;     }
;; 
;; 
;;     int main(int argc, char const *argv[]) {
;;       int r = fib(10);
;;       printf("%d", r); // 55
;;     }
;

section .text
global fib

fib:  xor   rax, rax
      call  _fib
      ret

_fib: cmp       rdi, 1
      jle       accumulate
      push      rdi
      sub       rdi, 1
      call      _fib
      sub       rdi, 1
      call      _fib
      pop       rdi
      ret

accumulate:   add   rax, rdi
              ret
