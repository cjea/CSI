section .text
global sum_to_n

sum_to_n:
	xor rax, rax
	jmp _sum_to_n

_sum_to_n:
	add rax, rdi
	cmp rdi, 0
	jnz _next
	ret

_next:
	sub rdi, 1
	jmp _sum_to_n
