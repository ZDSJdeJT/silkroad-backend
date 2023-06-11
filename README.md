# 丝绸之路

[English](https://github.com/ZDSJdeJT/silkroad-backend/blob/main/README_en.md)

这是一个用于快速传输文件的项目。它支持国际化，提供用户友好的界面，易于使用。

## 功能特性

- 快速传输：快传可在数秒内安全、稳定地将文件从一台计算机传输到另一台
- 用户界面友好：快传提供简单直观的用户界面，使用户无需专业知识即可完成文件传输

## 安装和使用

1. 克隆该项目至本地计算机
2. 打开项目根目录，运行 `go mod download` 命令以安装所需的依赖项
3. 如果您是 Windows 用户，请运行 `.\start.development.bat` 启动应用程序，如果您是 Linux 或 Mac 用户，请运行 `bash start.development.sh` 启动应用程序。注意，请不要将这种运行方式投入生产

## 参与贡献

如果您对该项目有任何建议或发现了问题，请在后端项目中创建一个 issue。我们非常感谢您的帮助和支持！

## 部署

1. 仿照 .env.sample 创建 .env.production
2. `docker build -t silkroad .`
3. `docker run -d --name silkroad -p <port>:4000 silkroad`
4. 访问 http://<ip>:<port>/admin/login 地址修改管理员密码（管理员名称和密码初始均为`admin`）
5. 访问 http://<ip>:<port> 进行使用

## 许可证

本项目是基于 LGPL-3.0 许可证开发的。有关详细信息，请参阅 [LICENSE](https://github.com/ZDSJdeJT/silkroad-backend/blob/main/LICENSE) 文件。
