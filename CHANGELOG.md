# 更新日志

## [Unreleased] - 2026-02-07

### 新增 🎉

- ⭐ **Stream V2 模式支持** - 集成官方 `dingtalk-stream-sdk-go`
  - 使用 Builder 模式创建客户端
  - 支持机器人消息监听
  - 支持事件监听
  - 无需公网 IP，自动重连
  - 新增 `examples/stream_v2/main.go` 示例

### 文档 📚

- 新增 `docs/STREAM_V2_GUIDE.md` - Stream V2 完整使用指南
- 更新 `docs/QUICK_START.md` - 添加 Stream 模式快速入门
- 更新 `README.md` - 添加核心模式对比表格

### 优化 ✨

- 更新依赖: 添加 `github.com/open-dingtalk/dingtalk-stream-sdk-go v0.9.1`
- 优化项目结构说明，更清晰的示例分类

---

# DingTalk SDK

[![Go Version](https://img.shields.io/badge/Go-1.22+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## 项目概述

DingTalk SDK 是一个基于 Go 语言开发的钉钉机器人 SDK，提供完整的钉钉消息发送、接收、流式卡片等功能。本项目从 [chatgpt-dingtalk](https://github.com/eryajf/chatgpt-dingtalk) 项目中提取钉钉相关模块，并重构为独立的 SDK。

## 模块说明

### 1. Client 模块 (`client/`)
提供钉钉客户端核心功能：
- OAuth 2.0 认证
- AccessToken 自动获取和缓存
- 媒体文件上传（图片、视频、文件等）
- 客户端管理器，支持多应用管理

### 2. Message 模块 (`message/`)
提供消息接收和发送功能：
- 接收钉钉回调消息
- 发送文本消息
- 发送 Markdown 格式消息
- @用户功能
- 消息格式化

### 3. Stream 模块 (`stream/`)
提供流式卡片功能：
- 创建并投放流式卡片
- 流式更新卡片内容
- 支持私聊和群聊场景
- 自动管理卡片生命周期

## 核心特性

✅ **简单易用** - 提供清晰的 API 接口，快速集成  
✅ **功能完整** - 涵盖消息、认证、流式卡片等核心功能  
✅ **自动管理** - AccessToken 自动缓存和刷新  
✅ **类型安全** - 完整的类型定义和接口设计  
✅ **可扩展** - 模块化设计，便于扩展新功能  

## 与原项目的关系

本 SDK 提取自 [eryajf/chatgpt-dingtalk](https://github.com/eryajf/chatgpt-dingtalk) 项目的以下模块：
- `pkg/dingbot/client.go` → `client/client.go`
- `pkg/dingbot/dingbot.go` → `message/message.go`
- `pkg/dingbot/stream.go` → `stream/stream.go`

## 下一步计划

- [ ] 添加更多测试用例
- [ ] 支持更多消息类型（卡片消息、互动卡片等）
- [ ] 添加 Webhook 签名验证
- [ ] 支持企业内部应用
- [ ] 提供更多使用示例

## 参考文档

- [钉钉开放平台文档](https://open.dingtalk.com/)
- [钉钉机器人开发指南](https://open.dingtalk.com/document/robots/robot-overview)
- [流式消息更新 API](https://open.dingtalk.com/document/development/api-streamingupdate)

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件
