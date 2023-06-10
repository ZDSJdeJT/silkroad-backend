# 丝绸之路

[English](https://github.com/ZDSJdeJT/silkroad-backend/blob/main/README_en.md)

## 部署
仿照 .env.sample 创建 .env.production
docker build -t silkroad .
docker run -d --name silkroad -p <port>:<port> silkroad
访问 http://<ip>:<port>