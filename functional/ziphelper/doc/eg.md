# 使用实例

实例1：

```
zipHelper := NewZipHelper()
zipHelper.AddDir("D:\\test")
data, err := zipHelper.Compress()
```

实例2：

```
zipHelper := NewZipHelper()
zipHelper.Add("1.bin", []byte("1test2test3test4test"))
zipHelper.Add("2.json", []byte("I have nothing in the world"))
data, err := zipHelper.Compress()
```