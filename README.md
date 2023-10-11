# Go语言学习记录
* [基础学习](#basic)
* [作业提交](#work)
* [知识点整理](#point)



<span id="basic"></span>
## 基础学习

<span id="work"></span>
## 作业提交
## 第一周 实现切片的删除操作
> 实现删除切片特定下标元素的方法。 
> * 要求一：能够实现删除操作就可以。
> * 要求二：考虑使用比较高性能的实现。
> * 要求三：改造为泛型方法
> * 要求四：支持缩容，并旦设计缩容机制。

详情点击：[sliceDeleteIdx.go](homework/week1/sliceDeleteIdx.go)

## 第二周 完善用户更新查看个人信息接口
需要完善 /users/edit 对应的接口。要求：

* 允许用户补充基本个人信息，包括：
    * 昵称：字符串，你需要考虑允许的长度。
    * 生日：前端输入为 1992-01-01 这种字符串。
    * 个人简介：一段文本，你需要考虑允许的长度。
* 尝试校验这些输入，并且返回准确的信息。
* 修改 /users/profile 接口，确保这些信息也能输出到前端。

详情点击：[第二周作业明细](homework/week2/README.md)

<span id="point"></span>
## 知识点整理
* [defer实现机制](KnowledgeBase/defer实现机制.md)
* [切片:扩容](KnowledgeBase/切片扩容.md)