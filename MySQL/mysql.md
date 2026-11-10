可以，这很适合你现在的学习场景：**MySQL 放在 Docker 中运行，Go 程序直接在 Windows 上运行**，通过端口连接。

```text
Windows 上的 Go 程序
        ↓ 127.0.0.1:3306
Docker 中的 MySQL
```

按下面的步骤操作即可。

**1. 启动 MySQL 容器**

先启动 Docker Desktop，确保使用的是 Linux 容器，然后在 PowerShell 中执行：

```powershell
docker run -d --name go-mysql `
  -p 127.0.0.1:3306:3306 `
  -e MYSQL_ROOT_PASSWORD=RootPass123 `
  -e MYSQL_DATABASE=go_learning `
  -e MYSQL_USER=gouser `
  -e MYSQL_PASSWORD=GoPass123 `
  -v go-mysql-data:/var/lib/mysql `
  mysql:8.4
```

注意：PowerShell 的续行符是反引号 `` ` ``，后面不能有空格。

这条命令会自动下载 MySQL 镜像，并创建：

| 配置 | 值 |
|---|---|
| Windows 连接地址 | `127.0.0.1:3306` |
| 数据库 | `go_learning` |
| Go 使用的账号 | `gouser` |
| Go 使用的密码 | `GoPass123` |
| 数据存储卷 | `go-mysql-data` |

`-v` 会将数据库保存在 Docker 数据卷中，删除容器后数据仍可保留。这些初始化变量只在数据目录为空时生效，之后修改启动命令里的密码不会修改已有账号密码。[MySQL 官方镜像说明](https://hub.docker.com/_/mysql)

**2. 确认 MySQL 已经启动**

```powershell
docker logs -f go-mysql
```

首次启动需要初始化。看到类似 `ready for connections` 的日志后，按 `Ctrl+C` 退出日志查看，容器会继续运行。

接着进入容器中的 MySQL 客户端：

```powershell
docker exec -it go-mysql mysql -u gouser -p go_learning
```

输入密码 `GoPass123`，然后执行：

```sql
SELECT VERSION();
SHOW DATABASES;
exit;
```

这里使用的是容器内自带的客户端，你的 Windows 不需要安装 MySQL。

**3. 用 Go 连接**

在你的 Go 项目目录中安装驱动：

```powershell
go get github.com/go-sql-driver/mysql
```

如果当前目录还没有 `go.mod`，需要先运行 `go mod init example.com/mysql-demo`；已有项目就跳过。

新建一个 `main.go`，写入：

```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "gouser:GoPass123@tcp(127.0.0.1:3306)/go_learning?parseTime=true&charset=utf8mb4"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// sql.Open 创建连接池，PingContext 才会实际检查连接。
	if err := db.PingContext(ctx); err != nil {
		log.Fatal("连接 MySQL 失败：", err)
	}

	fmt.Println("连接 MySQL 成功！")

	_, err = db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(100) NOT NULL
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	// 用 ? 传入参数，避免通过拼接字符串构造 SQL。
	result, err := db.ExecContext(ctx,
		"INSERT INTO users (name) VALUES (?)", "小明")
	if err != nil {
		log.Fatal(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}

	var name string
	err = db.QueryRowContext(ctx,
		"SELECT name FROM users WHERE id = ?", id,
	).Scan(&name)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("查询成功：id=%d，name=%s\n", id, name)
}
```

运行：

```powershell
go run main.go
```

正常输出类似：

```text
连接 MySQL 成功！
查询成功：id=1，name=小明
```

Go 的 `database/sql` 提供数据库操作接口，MySQL 驱动负责具体连接；`sql.Open` 后用 `PingContext` 验证连接是官方文档中的用法。[Go 数据库连接说明](https://go.dev/doc/database/open-handle)、[MySQL 驱动文档](https://github.com/go-sql-driver/mysql)

**4. 平时怎么启动和停止**

```powershell
# 暂时不用时停止
docker stop go-mysql

# 下次学习时启动，不用再次执行 docker run
docker start go-mysql

# 查看容器状态，包括已经停止的容器
docker ps -a
```

如果提示 **3306 端口被占用**，启动容器时把映射改为 `127.0.0.1:3307:3306`，Go 连接字符串也改为 `127.0.0.1:3307` 即可。