package wasm

import (
	"bytes"
	"fmt"
	"io"
)

const (
	opEnd      = 0x0B // end
	opI32Const = 0x41 // i32.const
	opI64Const = 0x42 // i64.const
)

// Eval tries to evaluate the (constant) expression in the buffer. It returns
// the evaluation stack at the end of the sequence or an error, if there was
// one.
func Eval(r *bytes.Buffer) ([]interface{}, error) {
	// Right now, this is only meant to parse very simple expressions like the
	// one used for the start offset of a data segment (which is an expression).
	var stack []interface{}
	for {
		b, err := r.ReadByte()
		if err == io.EOF {
			return nil, io.ErrUnexpectedEOF
		}
		if b == opEnd {
			// End of expression.
			break
		}
		switch b {
		case opI32Const:
			var n int32
			err := readVarInt32(r, &n)
			if err != nil {
				return nil, err
			}
			stack = append(stack, n)
		default:
			return nil, fmt.Errorf("unknown opcode: 0x%02X", b)
		}
	}
	return stack, nil
}

// Read a sequence of instructions used to initialize something.
func readInitExpr(r io.Reader, expr *[]byte) error {
	for {
		// Read a single byte.
		b, err := readByte(r)
		if err != nil {
			return err
		}
		*expr = append(*expr, b)
		switch b {
		case opEnd:
			return nil // finished the expression
		case opI32Const, opI64Const:
			// Read a single LEB128 encoded number.
			for {
				b, err := readByte(r)
				if err != nil {
					return err
				}
				*expr = append(*expr, b)
				if (b & 0x80) == 0 {
					break
				}
			}
		default:
			return fmt.Errorf("unknown opcode in init expression: 0x%02X", b)
		}
	}
}
