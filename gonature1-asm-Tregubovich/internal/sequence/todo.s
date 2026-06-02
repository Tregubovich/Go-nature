#include "textflag.h"

// func Fibonacci(n uint64) uint64
TEXT ·Fibonacci(SB), NOSPLIT, $0
    MOVQ n+0(FP), AX

    MOVQ $0, R9
    MOVQ $1, R10

loop:
    CMPQ AX, $0
    JE end

    MOVQ R10, R11
    ADDQ R9, R10
    MOVQ R11, R9

    SUBQ $1, AX
    JMP loop

end:
    MOVQ R9, ret+8(FP)
    RET
