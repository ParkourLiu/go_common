package localcipher

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

func ChromeVauleDecrypt(data []byte) {

}
func DpapiEncrypt(data []byte) ([]byte, error) {
	out := windows.DataBlob{}

	err := windows.CryptProtectData(bytesToBlob(data), nil, nil, 0, nil, windows.CRYPTPROTECT_UI_FORBIDDEN, &out)

	if err != nil {
		return nil, fmt.Errorf("unable to encrypt DPAPI protected data: %w", err)
	}
	ret := make([]byte, out.Size)
	copy(ret, unsafe.Slice(out.Data, out.Size))
	windows.LocalFree(windows.Handle(unsafe.Pointer(out.Data)))
	return ret, nil
}

func DpapiDecrypt(data []byte) ([]byte, error) {
	out := windows.DataBlob{}
	err := windows.CryptUnprotectData(bytesToBlob(data), nil, nil, 0, nil, 0, &out)
	if err != nil {
		return nil, fmt.Errorf("unable to decrypt DPAPI protected data: %w", err)
	}
	ret := make([]byte, out.Size)
	copy(ret, unsafe.Slice(out.Data, out.Size))
	windows.LocalFree(windows.Handle(unsafe.Pointer(out.Data)))

	return ret, nil
}
func bytesToBlob(bytes []byte) *windows.DataBlob {
	blob := &windows.DataBlob{Size: uint32(len(bytes))}
	if len(bytes) > 0 {
		blob.Data = &bytes[0]
	}
	return blob
}
