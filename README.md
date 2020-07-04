# jutf
This library implements support for the [modified UTF-8 scheme][1] used in Java.
In particular, this is the format used by the `DataInputStream#readUTF` and
`DataOutputStream#writeUTF` methods.

The library exports two functions:
````go
func Decode(d []byte) (string, error)
func Encode(s string) []byte
````

## License
MIT. See [LICENSE][2].

[1]: https://docs.oracle.com/javase/7/docs/api/java/io/DataInput.html#modified-utf-8 
[2]: ./LICENSE
