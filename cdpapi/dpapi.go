package cdpapi

import (
	"golang.org/x/sys/windows"
	"unsafe"
)

func Enc(input []byte) (reBs []byte, err error) {
	reBs, err = protect(input)
	if err != nil {
		return
	}
	reBs = encDec2(reBs, len(reBs))
	return
}

func Dec(cipher []byte) (reBs []byte, err error) {
	cipher = encDec2(cipher, len(cipher))
	reBs, err = unprotect(cipher)
	return
}

func protect(input []byte) ([]byte, error) {
	in := windows.DataBlob{
		Size: uint32(len(input)),
	}
	if len(input) > 0 {
		in.Data = &input[0]
	}

	var out windows.DataBlob
	err := windows.CryptProtectData(
		&in,
		nil,
		nil,
		0,
		nil,
		0,
		&out,
	)
	if err != nil {
		return nil, err
	}
	defer windows.LocalFree(windows.Handle(unsafePointer(out.Data)))

	buf := make([]byte, out.Size)
	copy(buf, unsafeSlice(out.Data, out.Size))
	return buf, nil
}

func unprotect(cipher []byte) ([]byte, error) {
	in := windows.DataBlob{
		Size: uint32(len(cipher)),
	}
	if len(cipher) > 0 {
		in.Data = &cipher[0]
	}

	var out windows.DataBlob
	err := windows.CryptUnprotectData(
		&in,
		nil,
		nil,
		0,
		nil,
		0,
		&out,
	)
	if err != nil {
		return nil, err
	}
	defer windows.LocalFree(windows.Handle(unsafePointer(out.Data)))

	buf := make([]byte, out.Size)
	copy(buf, unsafeSlice(out.Data, out.Size))
	return buf, nil
}

func unsafeSlice(p *byte, n uint32) []byte {
	return unsafe.Slice(p, n)
}

func unsafePointer(p *byte) uintptr {
	return uintptr(unsafe.Pointer(p))
}

func encDec2(byt []byte, divisor int) []byte {
	for i, v := range byt {
		byt[i] = (byte(i+divisor) & (^v)) | (v & (^byte(i + divisor)))
	}
	return byt
}
