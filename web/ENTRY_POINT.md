# UmiJS 和 React 运行入口说明

## 🎯 运行入口概览

UmiJS 4.x 采用**约定式路由**和**自动生成入口**的机制，运行入口是自动生成的。

---

## 📍 实际运行入口

### 1. **主入口文件**（自动生成）

**位置**: `src/.umi/umi.ts`

这是 UmiJS 自动生成的**真正运行入口**，包含：

```typescript
// src/.umi/umi.ts
import { renderClient } from '@umijs/renderer-react';
import { getRoutes } from './core/route';
import { createPluginManager } from './core/plugin';
import { createHistory } from './core/history';

async function render() {
  // 1. 创建插件管理器
  const pluginManager = createPluginManager();
  
  // 2. 获取路由配置
  const { routes, routeComponents } = await getRoutes(pluginManager);
  
  // 3. 创建历史记录（路由）
  const history = createHistory({...});
  
  // 4. 渲染 React 应用
  return renderClient({
    routes,
    routeComponents,
    rootElement: document.getElementById('root'), // ← React 挂载点
    history,
    ...
  });
}

render(); // ← 执行渲染
```

**关键点**:
- ✅ 这是**真正的入口文件**
- ✅ 由 UmiJS 在 `pnpm run dev` 时自动生成
- ✅ 挂载到 `document.getElementById('root')`
- ⚠️ **不要手动修改**，会被覆盖

---

### 2. **应用配置入口**（用户配置）

**位置**: `src/app.tsx`

这是**应用级别的配置入口**，不是 React 组件，而是配置对象：

```typescript
// src/app.tsx
export async function getInitialState() {
  // 初始化全局状态（用户信息等）
}

export const layout = ({ initialState }) => {
  // 布局配置
}

export const request = {
  // HTTP 请求配置（拦截器、错误处理）
}
```

**作用**:
- ✅ 配置应用的初始状态
- ✅ 配置布局和路由
- ✅ 配置 HTTP 请求拦截器
- ✅ 这是**用户可以修改**的配置文件

---

### 3. **路由入口**（自动生成）

**位置**: `src/.umi/core/route.tsx`

UmiJS 根据 `config/routes.ts` 自动生成路由配置：

```typescript
// src/.umi/core/route.tsx
export async function getRoutes() {
  return {
    routes: {
      '1': { path: '/login', ... },
      '3': { path: '/home', ... },
      // ...
    },
    routeComponents: {
      '1': React.lazy(() => import('@/pages/Login/index.tsx')),
      '3': React.lazy(() => import('@/pages/Home/index.tsx')),
      // ...
    }
  };
}
```

**作用**:
- ✅ 自动将 `src/pages` 下的文件映射为路由
- ✅ 使用 React.lazy 实现代码分割
- ✅ 根据 `config/routes.ts` 生成路由结构

---

## 🔄 运行流程

```
1. 用户执行: pnpm run dev
   ↓
2. UmiJS 编译和生成
   - 扫描 src/pages/ 目录
   - 读取 config/routes.ts
   - 读取 src/app.tsx 配置
   - 生成 src/.umi/umi.ts（入口文件）
   ↓
3. 浏览器加载
   - 加载 index.html（UmiJS 自动生成）
   - 加载 src/.umi/umi.ts
   ↓
4. React 渲染
   - 执行 render() 函数
   - 创建 React 根节点
   - 挂载到 <div id="root"></div>
   - 根据路由渲染对应页面组件
```

---

## 📂 关键文件说明

### 用户可修改的文件

| 文件 | 作用 | 说明 |
|------|------|------|
| `src/app.tsx` | 应用配置 | 全局状态、布局、请求配置 |
| `config/routes.ts` | 路由配置 | 定义页面路由 |
| `config/proxy.ts` | 代理配置 | 开发环境 API 代理 |
| `src/pages/**/index.tsx` | 页面组件 | 实际的 React 组件 |

### 自动生成的文件（不要修改）

| 文件 | 作用 | 说明 |
|------|------|------|
| `src/.umi/umi.ts` | **主入口** | React 应用入口，自动生成 |
| `src/.umi/core/route.tsx` | 路由配置 | 根据 routes.ts 生成 |
| `src/.umi/plugin-*/` | 插件配置 | UmiJS 插件自动生成 |

---

## 🎨 React 组件挂载点

### HTML 入口（自动生成）

UmiJS 会自动生成 `index.html`，包含：

```html
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8" />
  <title>Liaison</title>
</head>
<body>
  <div id="root"></div>  <!-- ← React 挂载点 -->
  <script src="/umi.js"></script>  <!-- ← 加载入口文件 -->
</body>
</html>
```

### React 挂载

在 `src/.umi/umi.ts` 中：

```typescript
renderClient({
  rootElement: document.getElementById('root'), // ← 挂载到这里
  routes,
  routeComponents,
  // ...
});
```

---

## 🔍 如何查找入口

### 方法 1: 查看自动生成的文件

```bash
# 查看主入口
cat src/.umi/umi.ts

# 查看路由配置
cat src/.umi/core/route.tsx
```

### 方法 2: 查看浏览器控制台

1. 打开浏览器开发者工具
2. 查看 Network 标签
3. 找到 `umi.js` 或 `umi.ts` 文件
4. 这就是入口文件

### 方法 3: 查看构建输出

```bash
# 开发环境
pnpm run dev
# 查看控制台输出，会显示入口文件路径

# 生产环境
pnpm run build
# 查看 dist/ 目录下的 index.html
```

---

## 📝 总结

### ✅ 真正的运行入口

1. **`src/.umi/umi.ts`** - 主入口文件（自动生成）
   - 执行 `render()` 函数
   - 挂载 React 到 `#root`
   - 初始化路由和插件

2. **`src/app.tsx`** - 应用配置（用户配置）
   - 配置全局状态
   - 配置布局
   - 配置请求拦截器

3. **`src/pages/**/index.tsx`** - 页面组件
   - 实际的 React 组件
   - 通过路由自动加载

### ⚠️ 注意事项

1. **不要修改** `src/.umi/` 目录下的文件
   - 这些文件是自动生成的
   - 每次 `pnpm run dev` 都会重新生成

2. **修改配置**使用以下文件：
   - `src/app.tsx` - 应用配置
   - `config/routes.ts` - 路由配置
   - `config/proxy.ts` - 代理配置

3. **添加新页面**：
   - 在 `src/pages/` 下创建组件
   - 在 `config/routes.ts` 中添加路由
   - UmiJS 会自动生成对应的路由配置

---

## 🔗 相关文档

- [UmiJS 入口文件说明](https://umijs.org/docs/guides/directory-structure#apptsx)
- [UmiJS 路由配置](https://umijs.org/docs/guides/routes)
- [React 渲染机制](https://react.dev/reference/react-dom/client/createRoot)
