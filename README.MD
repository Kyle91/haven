### 基本说明
1. 仓库属于基建项目，包括一些通用功能的封装
2. `log`、`goroutine pool`
3. 目前是私有仓库，以后看情况会不会改成public

### 导入说明
1. 目前是私有仓库，首先在github生成token，设置地址[github](https://github.com/settings/tokens)
2. 然后在本地配置好环境，
3. 在你的主目录中创建或打开 .netrc 文件。
- 对于 Windows 用户，路径是 %USERPROFILE%\.netrc。
- 对于 Linux 和 macOS 用户，路径是 ~/.netrc。
4. 在 .netrc 文件中添加以下内容：
```
machine github.com
login YOUR_GITHUB_USERNAME
password YOUR_PERSONAL_ACCESS_TOKEN
```
5. 使用`go get github.com/Kyle91/haven@latest`获取依赖