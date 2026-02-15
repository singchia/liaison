# Liaison 官网（静态站点）

该目录提供一个可直接部署到公网的产品官网静态页，面向“产品介绍 / 优势 / 架构 / 快速开始 / 文档入口”场景。

## 本地预览

```bash
cd website
python3 -m http.server 5173
```

然后访问 `http://localhost:5173/`。

## 生产部署

把 `website/` 目录作为静态站点发布即可（任意对象存储 + CDN、Nginx、Caddy、GitHub Pages 等均可）。

### Nginx 示例

```nginx
server {
  listen 80;
  server_name example.com;

  root /var/www/liaison-website;
  index index.html;

  location / {
    try_files $uri $uri/ /index.html;
  }
}
```

## 关联文档与仓库

官网中的“文档/Swagger/API/安装指南”等链接默认指向 GitHub 仓库页面，适合将 `website/` 作为独立静态站点单独部署。若你希望跳转到自建文档站/控制台域名，可按需调整 `index.html` 中的链接目标。
