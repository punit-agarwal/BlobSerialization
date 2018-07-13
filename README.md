# Blob Serialization

---

_For reference, look at the following post:_ [_https://ethresear.ch/t/blob-serialisation/1705_](https://ethresear.ch/t/blob-serialisation/1705)

This will deal with the implementation of Blob serialization in XVM. A few conditions that will have to hold to build a blob parser:

``` 
    COLLATION_BODY :=   2 ^20
    CHUNK_SIZE := 32
    INDICATOR_SIZE := 1
    CHUNK_DATA_SIZE := CHUNK_SIZE - INDICATOR_SIZE
    COLLATION_BODY % CHUNK_SIZE == 0
```
How the Blob parser would work would be as follows:

- The parser would check if the collation body was of the required size where                          0 &lt; COLLATION\_BODY &lt;= 2^20
- Then it would split the body into chunks of 31 bytes. This part of the chunk would contain the data from collation.
- The other 1 byte would represent the indicator byte. The 5 least significant bits signify index of the chunk delimiter. If the length is zero, it specifies a non-terminal chunk.

Below is how a Parser would check for validity of the body:

```golang
func(collationbody []bytes) checkcollationbody() {
  return len(collationbody) <= 2^20 && len(collationbody) > 0
}
func createChunks(collationbody []bytes) {
     validCollation,err := collationbody.checkcollationbody()
     if err !=nil {
         fmt.Errorf(....)
     }
     if(!validCollation) {
         fmt.Errorf(...)
     }
}
```

- After the COLLATION\_BODY has been determined as valid we can start packaging the chunks. We do this by finding the total size of the collation body. And find what is the closest multiple of 31 to that figure.
- Then when we look at len(COLLATION\_BODY) - x\*31 . This signifies the remaining bytes that will be allocated in the terminal chunk. X\*31 represents the multiple of 31 that is closest to the size of the collation body.
- Then once we have separated out 31 bytes for each chunk we can append each non-terminal chunk with a indicator byte with a value of zero. (We ignore the 3 flags which represent the most significant bits of the indicator byte)
- In the terminal chunk data bytes are padded to the left. And the bytes not holding blob data are ignored.



Below is how we would end up concatenating the indicator byte with the data bytes into a 32 byte chunk:
```golang
index := int32(len(collationbody) / chunkdatasize)
var serialisedblob byte[]
for i = 0; i < ; index; i++ {
serialisedblob[i\*chunksize] = 0
for f = 0; f <chunkdatasize ; f++ {
serialisedblob[(f+1) + i*chunksize] = collationbody[f + i*chunkdatasize]
}
serialisedblob[index*chunksize] = len(collationbody) - index*chunkdatasize
} 

```


After this the terminal chunk will add the remaining data bytes and append an indicator byte specifying the size of the data bytes.
