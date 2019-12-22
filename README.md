# jutf
This library implements support for the modified UTF-8 scheme [used in Java][1].
In particular, this is the format used by the `DataInputStream.readUTF` and
`DataOutputStream.writeUTF` methods.

The library exports two functions:
````
func Decode(d []byte) (string, error)
func Encode(s string) []byte
````

[1]: https://docs.oracle.com/javase/7/docs/api/java/io/DataInput.html#modified-utf-8 