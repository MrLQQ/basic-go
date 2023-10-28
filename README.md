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

> * 允许用户补充基本个人信息，包括：
>    * 昵称：字符串，你需要考虑允许的长度。
>    * 生日：前端输入为 1992-01-01 这种字符串。
>    * 个人简介：一段文本，你需要考虑允许的长度。
>* 尝试校验这些输入，并且返回准确的信息。
>* 修改 /users/profile 接口，确保这些信息也能输出到前端。

详情点击：[第二周作业明细](homework/week2/README.md)

## 第三周 修改已有的部署方案

> * 将 webook 的启动端口修改 8081。
> * 将 webook 修改为部署 2 个 Pod。
> * 将 webook 访问 Redis 的端口修改为 6380。
> * 将 webook 访问 MySQL 的端口修改为 3308。

详情点击：[第三周作业明细](homework/week3/README.md)

## 第四周 用本地缓存来替换 Redis
> 不使用 Redis 作为缓存，提供一个基于本地缓存实现的 cache.CodeCache。
> * 定义一个 CodeCache 接口，将现在的 CodeCache 改名为 CodeRedisCache。
> * 提供一个基于本地缓存的 CodeCache 实现。你可以自主决定用什么本地缓存，在这个过程注意体会技术选型要考虑的点。
> * 保证单机并发安全，也就是你可以假定这个实现只用在开发环境，或者单机环境下。

详情点击：[第四周作业明细](homework/week4/README.md)

<span id="point"></span>
## 知识点整理
* [defer实现机制](KnowledgeBase/defer实现机制.md)
* [切片:扩容](KnowledgeBase/切片扩容.md)