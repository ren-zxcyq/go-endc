# go-endc
Encrypt / Decrypt a file via the command line
AES - GCM

## Example:

### Let's create a test file.
```
root@kali:~/go/src/github.com/ren-zxcyq/go-endc# go installroot@kali:~/go/src/github.com/ren-zxcyq/go-endc# echo "testing" > test
```

### Encrypting test files contents.
```
root@kali:~/go/src/github.com/ren-zxcyq/go-endc# ~/go/bin/go-endc -s "/root/go/src/github.com/ren-zxcyq/go-endc/test" -d "/root/go/src/github.com/ren-zxcyq/go-endc/encryptedthefile" -p "01234567890123456789012345678911" -encrypt
[*]     Src:     /root/go/src/github.com/ren-zxcyq/go-endc/test
[*]     Dst:     /root/go/src/github.com/ren-zxcyq/go-endc/encryptedthefile
[*]     Encrypting file contents
```

### Decrypting encrypted files contents.
```
root@kali:~/go/src/github.com/ren-zxcyq/go-endc# ~/go/bin/go-endc -s "/root/go/src/github.com/ren-zxcyq/go-endc/encryptedthefile" -d "/root/go/src/github.com/ren-zxcyq/go-endc/testdecryption" -p "01234567890123456789012345678911" -decrypt
[*]     Src:     /root/go/src/github.com/ren-zxcyq/go-endc/encryptedthefile
[*]     Dst:     /root/go/src/github.com/ren-zxcyq/go-endc/testdecryption
[*]     Decrypting file contents
```

### Current directory file contents.
```
root@kali:~/go/src/github.com/ren-zxcyq/go-endc# ls
encryptedthefile  go-endc.go  test  testdecryption
root@kali:~/go/src/github.com/ren-zxcyq/go-endc# cat test
testing
root@kali:~/go/src/github.com/ren-zxcyq/go-endc# cat encryptedthefile 
>�Ԋ)���"��yP���nDT�|��D��ކ=}��<B 
root@kali:~/go/src/github.com/ren-zxcyq/go-endc# cat testdecryption 
testing
```
