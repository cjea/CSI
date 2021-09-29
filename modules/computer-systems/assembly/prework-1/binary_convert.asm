section .text
global binary_convert

binary_convert:
	xor rax, rax
	jmp loop

loop:
	movzx ecx, byte [rdi] ; Load the next input bit char into a holding register.
	jmp check_done		; Check if it's the null byte (end of string) before continuing.

check_done:	; check_done returns from the function if ecx contains a null byte.
	cmp ecx, 0
	jnz body
	ret

body:						; body puts the next input bit into the output's lowest position.
								; If the bit char is '1', then output is `(output << 1) + 1`.
								; If the bit char is '0', then output is simply `(output << 1)`.
	sal rax, 1
	cmp ecx, '0'
	jz next_iter
	add rax, 1
	jmp next_iter

next_iter:
	add rdi, 1
	jmp loop
