# GeeGit - 从零实现 Git

一个用 Go 语言实现的 Git，提供两个版本：

## 📦 两个版本

### 🎓 教程版（1678 行）- 推荐学习使用

7 天渐进式教程，每天一个模块，详细注释，完整文档。

**特点**：
- ✅ 增量开发（每天基于前一天）
- ✅ 详细注释和文档
- ✅ 独立的测试示例
- ✅ 适合学习 Git 原理

**目录**：`day1-basic-objects/` ~ `day7-push/`

### ⚡ 极简版（653 行）- 推荐实际使用

所有功能合并到一个包，删除注释，极致精简。

**特点**：
- ✅ 仅 653 行代码
- ✅ 功能完全相同
- ✅ 代码极致精简
- ✅ 适合快速集成

**目录**：`minimal/`

## 代码量对比

| 版本 | 代码量 | 用途 | 特点 |
|------|--------|------|------|
| **教程版** | 1678 行 | 学习 | 详细文档、增量开发 |
| **极简版** | 653 行 | 实用 | 极致精简、快速集成 |
| go-git 7d6c5a56 | 2605 行 | 参考 | 仅读取功能 |

## 🎯 项目目标

从 [go-git](https://github.com/go-git/go-git) 项目的早期提交（commit `7d6c5a56`）出发，通过激进简化和功能扩展：

✅ **目标**：代码量 < 1000 行
✅ **实际**：极简版 **653 行**（教程版 1678 行）

实现的核心功能：

- ✅ 读取本地 Git 仓库（解析 .git 目录）
- ✅ 从远程仓库克隆/拉取（HTTP）
- ✅ 解析 Git 对象（Commit/Tree/Blob）
- ✅ Packfile 格式解析
- ✅ 创建 commit
- ✅ Push 推送到远端

## 设计原则

### 保留的核心功能
- Delta 解压算法（读取远程 packfile 必需）
- Packfile 格式解析
- Pkt-line 协议
- HTTP 客户端
- 对象存储读写

### 删除的内容（简化）
- Delta 压缩（写入时）- 总是发送完整对象
- 智能协商 - 总是获取/发送所有对象
- SSH 协议 - 仅 HTTP
- 认证 - 无凭据处理
- CRC/校验和验证 - 信任数据完整性
- 索引/暂存区 - 直接提交
- 工作树操作 - 仅 .git 目录
- 子模块、标签、多分支、合并提交

## 7 天学习路径

### ✅ Day 1: 基础对象模型（~150 行）
**已完成** - 理解 Git 的三种基础对象：Blob, Tree, Commit

- 文件: `day1-basic-objects/`
  - `hash.go` - SHA-1 哈希计算
  - `object.go` - 对象类型定义

### ✅ Day 2: 读取本地仓库（~300 行）
**已完成** - 从 .git 目录读取对象

- 文件: `day2-repository-reading/`
  - `repository.go` - 仓库操作、对象读取和解析
- **增量开发**: 复用 Day 1 的对象类型定义

### ✅ Day 3: Packfile 格式解析（~380 行）
**已完成** - Packfile 格式和 Delta 解压

- 文件: `day3-packfile-reading/`
  - `delta.go` - Delta 解压算法
  - `packfile.go` - Packfile 读取器
- **增量开发**: 复用 Day 1-2 的对象类型
- **核心**: Delta 复制和插入指令解析

### ✅ Day 4: 远程克隆（~380 行）
**已完成** - HTTP 协议和远程克隆

- 文件: `day4-remote-clone/`
  - `protocol.go` - Pkt-line 协议编解码
  - `client.go` - Git HTTP 客户端
- **增量开发**: 整合 Day 1-3 的所有代码
- **核心**: 可以从 GitHub 克隆仓库！

### ✅ Day 5: 共享工具（~80 行）
**已完成** - 格式化工具函数

- 文件: `common/`
  - `format.go` - Commit/Tree 格式化
- **增量开发**: 为 Day 6-7 提供工具函数

### ✅ Day 6: 写入对象和创建 Commit（~150 行）
**已完成** - 本地写入功能

- 文件: `day6-write-objects/`
  - `writer.go` - 对象写入器
- **增量开发**: 复用 Day 1, 5 的代码
- **核心**: 创建的仓库可被真实 git 识别！

### ✅ Day 7: Push 到远程（~200 行）
**已完成** - Push 协议实现

- 文件: `day7-push/`
  - `pusher.go` - Receive-pack 协议和 packfile 创建
- **增量开发**: 整合所有 7 天的代码
- **核心**: 完整的读写循环！

## 快速开始

### 环境要求
- Go 1.16+
- Git（用于验证结果）

### 运行单元测试

项目为每个 day 都编写了完整的单元测试，你可以独立运行每个 day 的测试：

```bash
# 运行所有测试
./run_tests.sh

# 运行单个 day 的测试
cd day1-basic-objects && go test -v
cd day2-repository-reading && go test -v
cd day3-packfile-reading && go test -v
cd day4-remote-clone && go test -v
cd day6-write-objects && go test -v
cd day7-push && go test -v

# 运行特定测试
cd day1-basic-objects && go test -v -run TestComputeHash
```

### 运行示例程序

```bash
# Day 1: 基础对象模型
cd day1-basic-objects/example
go run main.go

# Day 2: 读取本地仓库
cd day2-repository-reading/example
go run test_loose_object.go

# Day 3: Delta 算法
cd day3-packfile-reading/example
go run test_delta.go

# Day 4: Pkt-line 协议
cd day4-remote-clone/example
go run test_protocol.go

# Day 6: 创建 Commit
cd day6-write-objects/example
go run main.go

# Day 7: Push 演示
cd day7-push/example
go run main.go
```

## 项目结构

```
GeeGit/
├── day1-basic-objects/          # ✅ Git 对象模型
├── day2-repository-reading/     # ✅ 读取本地仓库
├── day3-packfile-reading/       # ✅ Packfile 解析
├── day4-remote-clone/           # ✅ 远程克隆
├── day5-shared-utils/           # ✅ 共享工具（文档）
├── day6-write-objects/          # ✅ 写入对象
├── day7-push/                   # ✅ Push 推送
├── common/                      # ✅ 共享代码
├── go.mod
└── README.md
```

## 代码统计

### 教程版（7天渐进）

| 模块 | 行数 | 状态 |
|------|------|------|
| Day 1: 基础对象 | 142 | ✅ 完成 |
| Day 2: 仓库读取 | 300 | ✅ 完成 |
| Day 3: Packfile | 380 | ✅ 完成 |
| Day 4: 远程克隆 | 380 | ✅ 完成 |
| Day 5: 共享工具 | 80 | ✅ 完成 |
| Day 6: 写入对象 | 150 | ✅ 完成 |
| Day 7: Push | 200 | ✅ 完成 |
| Common | 46 | ✅ 完成 |
| **教程版总计** | **1678** | **100%** |

### 极简版（单包）

| 文件 | 行数 | 功能 |
|------|------|------|
| objects.go | 82 | Hash + 对象定义 |
| repository.go | 190 | 仓库读写 |
| delta.go | 74 | Delta 算法 |
| packfile.go | 104 | Packfile 解析 |
| protocol.go | 44 | Pkt-line 协议 |
| client.go | 79 | HTTP 克隆 |
| push.go | 80 | Push 推送 |
| **极简版总计** | **653** | **100%** |

## 🎉 项目完成！

### 成就解锁

✅ **极简版 653 行**，低于 1000 行目标！
✅ **教程版 1678 行**，7 天增量开发，零修改旧代码
✅ **读写完整循环**：克隆 → 修改 → 提交 → 推送
✅ **功能完整**：比原始 go-git（2605行）代码少 75%，功能更多

### 可以做什么

- ✅ 读取真实的 Git 仓库
- ✅ 从 GitHub 克隆仓库
- ✅ 解析 packfile 和 delta 压缩
- ✅ 创建被真实 git 识别的 commit
- ✅ 理解 Git 的网络协议（pkt-line, upload-pack, receive-pack）

### 学到了什么

1. **Git 对象模型**：Blob, Tree, Commit 如何组织
2. **存储格式**：松散对象 vs Packfile
3. **压缩算法**：Delta 编码（复制+插入指令）
4. **网络协议**：Pkt-line, Smart HTTP
5. **增量开发**：如何构建复杂系统

### 架构图

```
┌─────────────────────────────────────────┐
│           GeeGit 架构                   │
└─────────────────────────────────────────┘

读取路径：
GitHub → HTTP → Pkt-line → Packfile → Delta → Objects

写入路径：
Objects → Zlib → .git/objects → Packfile → HTTP → GitHub

核心模块：
├── Day 1: 对象模型（Hash, Type, Compute）
├── Day 2: 本地读取（Repository, Storage）
├── Day 3: 压缩算法（Delta, Packfile）
├── Day 4: 网络通信（Protocol, Client）
├── Day 5: 工具函数（Format）
├── Day 6: 本地写入（Writer）
└── Day 7: 网络推送（Pusher）
```

### 下一步

从学习项目到生产级别，还需要：

1. **认证**：SSH keys, Personal Access Token
2. **性能**：流式处理，增量传输，delta 压缩写入
3. **功能**：工作区管理（checkout, diff, merge），分支管理，标签
4. **健壮性**：错误处理，重试，断点续传
5. **用户体验**：CLI 命令，进度条，颜色输出

但核心原理都在这 1632 行代码里了！

## 参考资料

- [Git 内部原理](https://git-scm.com/book/zh/v2/Git-%E5%86%85%E9%83%A8%E5%8E%9F%E7%90%86-Git-%E5%AF%B9%E8%B1%A1)
- [go-git 项目](https://github.com/go-git/go-git)
- [Git Packfile 格式](https://git-scm.com/docs/pack-format)
- [Git 传输协议](https://git-scm.com/book/en/v2/Git-Internals-Transfer-Protocols)

## 许可证

MIT License

## 致谢

本项目基于 [go-git](https://github.com/go-git/go-git) 的早期实现（commit `7d6c5a56`），通过大量简化和教学改造而成。感谢 go-git 项目提供的优秀参考实现！

---

**GeeGit** - 从零实现 Git，7 天掌握 Git 内部原理 🚀
