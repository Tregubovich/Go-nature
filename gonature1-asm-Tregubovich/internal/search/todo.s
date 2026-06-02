#include "textflag.h"

// func LowerBound(slice []int64, value int64) int64
TEXT ·LowerBound(SB), NOSPLIT, $0
    MOVQ slice_base+0(FP), AX
    MOVQ slice_len+8(FP), CX
    MOVQ value+24(FP), DX

    XORQ R8, R8 // l
    MOVQ CX, R9 // r

loop:
    CMPQ R8, R9 // (l < r)
    JGE done

    MOVQ R9, R10 // m = r
    ADDQ R8, R10 // m = r + l
    SHRQ $1, R10 // m = (r + l) / 2

    MOVQ (AX)(R10*8), R11

    // (slice[m] < value)
    CMPQ R11, DX
    JGE else

    // l = m + 1
    MOVQ R10, R8
    ADDQ $1, R8
    JMP loop

else:
    // r = m
    MOVQ R10, R9
    JMP loop

done:
    MOVQ R9, ret+32(FP)
    RET
