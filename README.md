# UNCORRUPT
Encryptor/decryptor library written using custom algorithm.

## Usage

```
go get github.com/blackistheneworange/uncorrupt
```

```
For corrupting a set of bytes
uncorrupt.Corrupt([]byte, string) []byte

For uncorrupting a set of bytes
uncorrupt.Uncorrupt([]byte, string) []byte
```