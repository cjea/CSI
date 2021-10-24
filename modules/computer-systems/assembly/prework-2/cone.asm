default rel

section .text
global volume
volume: ; (PI * (r^2) * h) / 3
  mulss xmm0, xmm0
  mulss xmm0, [pi]
  mulss xmm0, xmm1
  divss xmm0, [f_three]
 	ret
  
section .rodata
pi:       dd 3.14159
f_three:  dd 3.0
