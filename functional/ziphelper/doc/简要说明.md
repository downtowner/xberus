# 简要说明

提供对zip文件的一个简单操作

# 详细说明

对zip文件常用的操作进行封装，主要的接口包括：

```
func NewZipHelper() *ZipHelper
```

创建一个zip助手对象

```
func (z *ZipHelper) Add(fName string, data []byte) error
```

向对象中添加一个文件，`fName` 文件名称，`data` 文件的数据，可以无限添加

```
func (z *ZipHelper) Compress() ([]byte, error)
```

压缩成zip包并返回数据

```
func (z *ZipHelper) Uncompress(data []byte) (map[string][]byte, error)
```

解压zip压缩包，`data` zip文件数据, 以文件名-数据的方式返回包含的文件信息

```
func (z *ZipHelper) AddDir(dir string) error
```

添加文件夹为zip文件源



# 缺陷

如果把文件夹作为源进行压缩，会导致文件夹内的层级改变，比如文件夹A包含文件夹B，压缩后全部文件都在文件夹A中，这个问题后续会修复。



欢迎大家交流使用,如有问题，请大家及时指正，联系邮箱 [xiaobing@novastar.tech](mailto:moubo@novastar.tech)