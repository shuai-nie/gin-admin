# gin-admin 管理后台

````
.
├── cmd/                    # 应用程序入口点
│   ├── api/               # REST API 服务
│   │   └── main.go
│   ├── cli/               # 命令行工具
│   │   └── main.go
│   └── worker/            # 后台工作进程
│       └── main.go
│
├── internal/              # 私有应用程序代码
│   ├── config/           # 配置结构
│   ├── handler/          # HTTP 处理器
│   ├── service/          # 业务逻辑
│   ├── repository/       # 数据访问层
│   └── models/           # 领域模型
│
├── pkg/                  # 公共库代码
│   ├── logger/          # 日志库
│   ├── database/        # 数据库包装
│   └── util/            # 通用工具
│
├── api/                  # API 定义
│   ├── openapi/         # OpenAPI 3.0 定义
│   └── proto/           # gRPC Proto 文件
│
├── configs/              # 配置文件
│   ├── default.yaml
│   ├── development.yaml
│   └── production.yaml
│
├── deployments/          # 部署配置
│   ├── docker/          # Docker 相关
│   ├── kubernetes/      # K8s 配置
│   └── terraform/       # 基础设施代码
│
├── scripts/              # 脚本文件
│   ├── build/           # 构建脚本
│   ├── migrate/         # 数据库迁移
│   └── deploy/          # 部署脚本
│
├── test/                 # 测试
│   ├── integration/     # 集成测试
│   ├── e2e/            # 端到端测试
│   └── fixtures/       # 测试数据
│
├── web/                  # Web 前端资源
│   ├── static/          # 静态文件
│   └── templates/       # Go 模板
│
├── docs/                 # 文档
│   ├── api.md
│   ├── architecture.md
│   └── development.md
│
├── third_party/          # 外部工具
├── vendor/               # 依赖项（Go Modules 前）
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile
├── .gitignore
├── README.md
└── LICENSE
````


https://github.com/LyricTian/gin-admin

https://github.com/go-admin-team/go-admin
