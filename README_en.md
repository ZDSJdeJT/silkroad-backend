# Silk Road

This is a project designed for fast file transfer. It supports internationalization, provides a user-friendly interface, and is easy to use.

## Features

- Fast Transfer: Silk Road can safely and stably transfer files from one computer to another in seconds
- User-Friendly Interface: Silk Road provides a simple and intuitive user interface that allows users to complete file transfers without specialized knowledge

## Installation and Usage

1. Clone the project to your local machine.
2. Open the root directory of the project and run the `go mod download` command to install the necessary dependencies.
3. If you are a Windows user, please run `.\start.development.bat` to start the application. If you are a Linux or Mac user, please run `bash start.development.sh` to start the application. Note: Do not use this method in production.

## Contributing

If you have any suggestions or find any issues with the project, please create an issue in the backend project. We appreciate your help and support!

## 部署

1. Create a .env.production file based on .env.sample.
2. Build the Silk Road Docker image with docker build -t silkroad .
3. Run the Silk Road Docker container with docker run -d --name silkroad -p \<port\>:\<port\> silkroad
4. Access http://\<ip\>:\<port\>/admin/login to change the administrator password (the initial administrator name and password are both `admin`)
5. Access http://\<ip\>:\<port\> to use Silk Road

## 许可证

This project is developed under the LGPL-3.0 license. For more information, see the [LICENSE](https://github.com/ZDSJdeJT/silkroad-backend/blob/main/LICENSE) file.
