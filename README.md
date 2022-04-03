# crc-sha

计算文件的hash值，支持md5,sha256等

```bash
用法：sha [-s hash] [-c value] file
  -c string
        该值不为空时，会将校验码和该值进行比较。
  -s string
        指定hash算法。 (default "sha256")
支持的hash算法：
  BLAKE2b_256
  BLAKE2b_384
  BLAKE2b_512
  BLAKE2s_256
  MD4
  MD5
  RIPEMD160
  SHA1
  SHA224
  SHA256
  SHA384
  SHA3_224
  SHA3_256
  SHA3_384
  SHA3_512
  SHA512
  SHA512_224
  SHA512_256
```
