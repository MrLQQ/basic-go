# BSON中的EDMA类型

> 对于BSON来说，第一部就是要理解:
> * E: <font color=Darkorange>一个普普通通的键值对结构体</font>，Value可以是其他三个
> * D: <font color=Darkorange>本质上是一个E的切片</font>
> * M: <font color=Darkorange>本质上是一个Map</font>,key必须是string,value可以是任何值，也可以是其他三种
> * A: <font color=Darkorange>切片</font>，元素可以是其他三种

