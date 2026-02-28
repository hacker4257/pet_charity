# 宠物公益平台

一个连接流浪动物救助组织与爱心人士的公益平台，提供宠物领养、救助上报、在线捐赠、实时聊天等功能。

## 功能概览

**宠物与领养**
- 浏览待领养宠物，查看详情、图片
- 提交领养申请，组织审核通过后完成领养
- 收藏感兴趣的宠物

**救助上报**
- 发现流浪动物可上报救助信息，附带地理位置
- 地图可视化查看周边救助事件
- 关注救助进展，认领救助任务，康复后转为可领养宠物

**捐赠**
- 支持微信支付、支付宝向组织或特定宠物捐赠
- 公开捐赠流水，透明可查

**组织管理**
- 救助组织注册、资质审核
- 管理旗下宠物信息、处理领养申请

**社区互动**
- 实时聊天（WebSocket）
- 活动积分与排行榜
- 站内通知系统
- 动态 Feed 流

**管理后台**
- 用户管理、角色分配
- 组织审核
- 数据统计面板

## 技术栈

| 层 | 技术 |
|----|------|
| 后端框架 | Go + Gin |
| 数据库 | MySQL 8.0 + GORM |
| 缓存 | Redis 7 |
| 消息队列 | Kafka |
| 前端 | React + TypeScript + Ant Design |
| 状态管理 | Zustand |
| 实时通信 | WebSocket（Gorilla） |
| 认证 | JWT + Token 版本号机制 |
| 支付 | 微信支付、支付宝 |
| 短信 | 阿里云 SMS |
| 地图 | 高德地图 |
| 部署 | Docker Compose + Nginx |

## 快速启动

```bash
# 克隆项目
git clone https://github.com/hacker4257/pet_charity.git
cd pet_charity

# Docker 一键启动
docker-compose up -d

# 或本地开发
# 后端
go run ./cmd/server/main.go

# 前端
cd web
npm install
npm run dev
```

配置文件位于 `configs/config.yaml`，支持通过环境变量覆盖（前缀 `PET_`）。

## 许可证

MIT
