# Copyright 2022 The go-python Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import binascii

print("globals:")
for name in ("Error", "Incomplete"):
    v = getattr(binascii, name)
    print("\nbinascii.%s:\n%s" % (name,repr(v)))

def assertEqual(x, y):
    assert x == y, "got: %s, want: %s" % (repr(x), repr(y))

rawdata = b"The quick brown fox jumps over the lazy dog.\r\n"
rawdata += bytes(range(256))
rawdata += b"\r\nHello world.\n"

## base64
assertEqual(binascii.b2a_base64(b'hello world!'), b'aGVsbG8gd29ybGQh\n')
assertEqual(binascii.b2a_base64(b'hello\nworld!'), b'aGVsbG8Kd29ybGQh\n')
assertEqual(binascii.b2a_base64(b'hello world!', newline=False), b'aGVsbG8gd29ybGQh')
assertEqual(binascii.b2a_base64(b'hello\nworld!', newline=False), b'aGVsbG8Kd29ybGQh')
assertEqual(binascii.a2b_base64("aGVsbG8gd29ybGQh\n"), b'hello world!')

assertEqual(binascii.b2a_base64(rawdata), b"VGhlIHF1aWNrIGJyb3duIGZveCBqdW1wcyBvdmVyIHRoZSBsYXp5IGRvZy4NCgABAgMEBQYHCAkKCwwNDg8QERITFBUWFxgZGhscHR4fICEiIyQlJicoKSorLC0uLzAxMjM0NTY3ODk6Ozw9Pj9AQUJDREVGR0hJSktMTU5PUFFSU1RVVldYWVpbXF1eX2BhYmNkZWZnaGlqa2xtbm9wcXJzdHV2d3h5ent8fX5/gIGCg4SFhoeIiYqLjI2Oj5CRkpOUlZaXmJmam5ydnp+goaKjpKWmp6ipqqusra6vsLGys7S1tre4ubq7vL2+v8DBwsPExcbHyMnKy8zNzs/Q0dLT1NXW19jZ2tvc3d7f4OHi4+Tl5ufo6err7O3u7/Dx8vP09fb3+Pn6+/z9/v8NCkhlbGxvIHdvcmxkLgo=\n")

try:
    binascii.b2a_base64("string")
    print("expected an exception")
except TypeError as e:
    print("expected an exception:", e)
    pass

## crc32
assertEqual(binascii.crc32(b'hello world!'), 62177901)
assertEqual(binascii.crc32(b'hello world!', 0), 62177901)
assertEqual(binascii.crc32(b'hello world!', 42), 4055036404)

## hex
assertEqual(binascii.b2a_hex(b'hello world!'), b'68656c6c6f20776f726c6421')
assertEqual(binascii.a2b_hex(b'68656c6c6f20776f726c6421'), b'hello world!')
assertEqual(binascii.hexlify(b'hello world!'), b'68656c6c6f20776f726c6421')
assertEqual(binascii.unhexlify(b'68656c6c6f20776f726c6421'), b'hello world!')

assertEqual(binascii.b2a_hex(rawdata), b'54686520717569636b2062726f776e20666f78206a756d7073206f76657220746865206c617a7920646f672e0d0a000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebfc0c1c2c3c4c5c6c7c8c9cacbcccdcecfd0d1d2d3d4d5d6d7d8d9dadbdcdddedfe0e1e2e3e4e5e6e7e8e9eaebecedeeeff0f1f2f3f4f5f6f7f8f9fafbfcfdfeff0d0a48656c6c6f20776f726c642e0a')

try:
    binascii.a2b_hex(b'123')
    print("expected an exception")
except binascii.Error as e:
    print("expected an exception:",e)
    pass

try:
    binascii.a2b_hex(b'hell')
    print("expected an exception")
except binascii.Error as e:
    print("expected an exception:",e)
    pass

## quotedprintable
assertEqual(binascii.b2a_qp(b'hello world! = \t'), b'hello world! =3D =09')
assertEqual(binascii.a2b_qp(b'hello world! =3D =09'), b'hello world! = \t')
## ## TODO
## #assertEqual(binascii.a2b_qp(b'hello world!', header=True), b'hello world!')
## assertEqual(binascii.a2b_qp(rawdata, header=False), b'The quick brown fox jumps over the lazy dog.\r\n\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\x0b\x0c\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f !"#$%&\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\x7f\x80\x81\x82\x83\x84\x85\x86\x87\x88\x89\x8a\x8b\x8c\x8d\x8e\x8f\x90\x91\x92\x93\x94\x95\x96\x97\x98\x99\x9a\x9b\x9c\x9d\x9e\x9f\xa0\xa1\xa2\xa3\xa4\xa5\xa6\xa7\xa8\xa9\xaa\xab\xac\xad\xae\xaf\xb0\xb1\xb2\xb3\xb4\xb5\xb6\xb7\xb8\xb9\xba\xbb\xbc\xbd\xbe\xbf\xc0\xc1\xc2\xc3\xc4\xc5\xc6\xc7\xc8\xc9\xca\xcb\xcc\xcd\xce\xcf\xd0\xd1\xd2\xd3\xd4\xd5\xd6\xd7\xd8\xd9\xda\xdb\xdc\xdd\xde\xdf\xe0\xe1\xe2\xe3\xe4\xe5\xe6\xe7\xe8\xe9\xea\xeb\xec\xed\xee\xef\xf0\xf1\xf2\xf3\xf4\xf5\xf6\xf7\xf8\xf9\xfa\xfb\xfc\xfd\xfe\xff\r\nHello world.\n')
## assertEqual(binascii.b2a_qp(binascii.a2b_qp(rawdata)), b'The quick brown fox jumps over the lazy dog.\r\n=00=01=02=03=04=05=06=07=08=09\r\n=0B=0C\r=0E=0F=10=11=12=13=14=15=16=17=18=19=1A=1B=1C=1D=1E=1F !"#$%&\'()*+,-=\r\n./0123456789:;<=3D>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuv=\r\nwxyz{|}~=7F=80=81=82=83=84=85=86=87=88=89=8A=8B=8C=8D=8E=8F=90=91=92=93=94=\r\n=95=96=97=98=99=9A=9B=9C=9D=9E=9F=A0=A1=A2=A3=A4=A5=A6=A7=A8=A9=AA=AB=AC=AD=\r\n=AE=AF=B0=B1=B2=B3=B4=B5=B6=B7=B8=B9=BA=BB=BC=BD=BE=BF=C0=C1=C2=C3=C4=C5=C6=\r\n=C7=C8=C9=CA=CB=CC=CD=CE=CF=D0=D1=D2=D3=D4=D5=D6=D7=D8=D9=DA=DB=DC=DD=DE=DF=\r\n=E0=E1=E2=E3=E4=E5=E6=E7=E8=E9=EA=EB=EC=ED=EE=EF=F0=F1=F2=F3=F4=F5=F6=F7=F8=\r\n=F9=FA=FB=FC=FD=FE=FF\r\nHello world.\r\n')
## ## TODO
## #assertEqual(binascii.a2b_qp(rawdata, header=True), b'The quick brown fox jumps over the lazy dog.\r\n\x00\x01\x02\x03\x04\x05\x06\x07\x08\t\n\x0b\x0c\r\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f !"#$%&\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^ `abcdefghijklmnopqrstuvwxyz{|}~\x7f\x80\x81\x82\x83\x84\x85\x86\x87\x88\x89\x8a\x8b\x8c\x8d\x8e\x8f\x90\x91\x92\x93\x94\x95\x96\x97\x98\x99\x9a\x9b\x9c\x9d\x9e\x9f\xa0\xa1\xa2\xa3\xa4\xa5\xa6\xa7\xa8\xa9\xaa\xab\xac\xad\xae\xaf\xb0\xb1\xb2\xb3\xb4\xb5\xb6\xb7\xb8\xb9\xba\xbb\xbc\xbd\xbe\xbf\xc0\xc1\xc2\xc3\xc4\xc5\xc6\xc7\xc8\xc9\xca\xcb\xcc\xcd\xce\xcf\xd0\xd1\xd2\xd3\xd4\xd5\xd6\xd7\xd8\xd9\xda\xdb\xdc\xdd\xde\xdf\xe0\xe1\xe2\xe3\xe4\xe5\xe6\xe7\xe8\xe9\xea\xeb\xec\xed\xee\xef\xf0\xf1\xf2\xf3\xf4\xf5\xf6\xf7\xf8\xf9\xfa\xfb\xfc\xfd\xfe\xff\r\nHello world.\n')

print("done.")
