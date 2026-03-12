# useI18nRouter - 多语言路由 Hook

自动处理语言前缀的路由导航，无需手动拼接 locale。

## 📖 简介

`useI18nRouter` 是一个封装了 next-intl 路由功能的 Hook，自动为所有路径添加当前语言前缀，让你专注于业务逻辑而无需关心多语言路由细节。

## ✨ 特性

- ✅ 自动添加语言前缀
- ✅ 支持所有 Next.js Router API
- ✅ 类型安全
- ✅ 与 next-intl 完美集成
- ✅ 简化国际化路由代码

## 🚀 快速开始

### 基础用法

```tsx
'use client';

import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

export function MyComponent() {
    const router = useI18nRouter();
    
    // 自动添加语言前缀（如 /zh-CN/category/1）
    return (
        <button onClick={() => router.push('/category/1')}>
            查看分类
        </button>
    );
}
```

## 📋 API

### 方法列表

| 方法 | 说明 | 返回值 |
|------|------|--------|
| `push(path, options?)` | 导航到指定路径（自动添加语言前缀） | `void` |
| `replace(path, options?)` | 替换当前历史记录（不保留历史） | `void` |
| `back()` | 返回上一页 | `void` |
| `forward()` | 前进 | `void` |
| `refresh()` | 刷新当前页面 | `void` |
| `localizedPath(path)` | 获取本地化路径（不执行导航） | `string` |

### 参数类型

```typescript
interface RouterOptions {
    scroll?: boolean; // 是否滚动到顶部，默认为 true
}
```

## 💡 使用示例

### 1. 基础导航

```tsx
'use client';

import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

export function NavigationExample() {
    const router = useI18nRouter();
    
    return (
        <div>
            {/* 自动添加语言前缀的导航 */}
            <button onClick={() => router.push('/category/1')}>
                查看分类
            </button>
            
            <button onClick={() => router.push('/post/123')}>
                查看文章
            </button>
            
            <button onClick={() => router.push('/settings/profile')}>
                个人设置
            </button>
            
            {/* 获取本地化路径 */}
            <a href={router.localizedPath('/about')}>
                关于我们
            </a>
            
            {/* 返回上一页 */}
            <button onClick={() => router.back()}>
                返回
            </button>
        </div>
    );
}
```

### 2. 在分类列表中使用

```tsx
'use client';

import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

export function CategoryListExample() {
    const router = useI18nRouter();
    
    const categories = [
        {id: 1, name: '技术'},
        {id: 2, name: '生活'},
        {id: 3, name: '随笔'},
    ];
    
    return (
        <ul>
            {categories.map(category => (
                <li key={category.id}>
                    <button 
                        onClick={() => router.push(`/category/${category.id}`)}
                    >
                        {category.name}
                    </button>
                </li>
            ))}
        </ul>
    );
}
```

### 3. 结合 useI18n 使用

```tsx
'use client';

import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';
import {useI18n} from '@/i18n/helpers/useI18n';

export function ArticleCardExample() {
    const t = useI18n('page.home');
    const router = useI18nRouter();
    
    const article = {
        id: 123,
        title: 'React Hooks 最佳实践',
        categoryId: 1,
    };
    
    return (
        <div>
            <h3>{article.title}</h3>
            
            <button onClick={() => router.push(`/post/${article.id}`)}>
                {t('read_more')}
            </button>
            
            <button onClick={() => router.push(`/category/${article.categoryId}`)}>
                {t('view_category')}
            </button>
        </div>
    );
}
```

### 4. 带参数的导航

```tsx
'use client';

import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

export function NavigationWithParamsExample() {
    const router = useI18nRouter();
    
    const handleViewPost = (postId: number) => {
        // 路径参数和查询参数都会自动保留
        router.push(`/post/${postId}?tab=comments`);
    };
    
    const handleSearch = (query: string) => {
        // 查询参数也会自动保留
        router.push(`/search?q=${encodeURIComponent(query)}`);
    };
    
    return (
        <div>
            <button onClick={() => handleViewPost(123)}>
                查看文章
            </button>
            <button onClick={() => handleSearch('React')}>
                搜索
            </button>
        </div>
    );
}
```

### 5. 替换历史记录（不保留当前历史）

```tsx
'use client';

import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

export function ReplaceExample() {
    const router = useI18nRouter();
    
    const handleLoginSuccess = () => {
        // 登录后替换到首页，用户无法通过后退回到登录页
        router.replace('/');
    };
    
    return (
        <button onClick={handleLoginSuccess}>
            登录成功，跳转
        </button>
    );
}
```

### 6. 实际应用场景

#### 文章详情页

```tsx
'use client';

import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';
import {useI18n} from '@/i18n/helpers/useI18n';

interface PostDetailProps {
    postId: number;
}

export default function PostDetail({postId}: PostDetailProps) {
    const t = useI18n('page.post');
    const router = useI18nRouter();
    
    return (
        <article>
            <header>
                <button onClick={() => router.back()}>
                    ← {t('back')}
                </button>
                <button onClick={() => router.push(`/category/1`)}>
                    {t('view_category')}
                </button>
            </header>
            
            <div className="content">
                {/* 文章内容 */}
            </div>
            
            <footer>
                <button onClick={() => router.push(`/post/${postId}?tab=comments`)}>
                    {t('view_comments')}
                </button>
            </footer>
        </article>
    );
}
```

#### 搜索结果页面

```tsx
'use client';

import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';
import {useState} from 'react';

export function SearchPage() {
    const router = useI18nRouter();
    const [query, setQuery] = useState('');
    
    const handleSearch = (e: React.FormEvent) => {
        e.preventDefault();
        if (query.trim()) {
            router.push(`/search?q=${encodeURIComponent(query.trim())}`);
        }
    };
    
    return (
        <div>
            <form onSubmit={handleSearch}>
                <input
                    type="text"
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                    placeholder="搜索..."
                />
                <button type="submit">搜索</button>
            </form>
        </div>
    );
}
```

## ⚠️ 注意事项

### 1. 路径格式

```typescript
// ✅ 正确：不需要语言前缀
router.push('/category/1');  // → /zh-CN/category/1

// ❌ 错误：不要手动添加语言前缀
router.push('/zh-CN/category/1');  // → /zh-CN/zh-CN/category/1
```

### 2. 外部链接

```typescript
// ✅ 外部链接使用原始 useRouter
import {useRouter} from 'next/navigation';

const router = useRouter();
router.push('https://example.com');
```

### 3. API 调用

```typescript
// ✅ API 调用不需要路由
const response = await fetch('/api/posts');
```

## 🔄 对比：新旧实现方式

### ❌ 旧方式（需要手动处理）

```typescript
import {useRouter} from 'next/navigation';

const router = useRouter();
const locale = 'zh-CN';

// 需要手动拼接语言前缀
router.push(`/${locale}/category/1`);
```

### ✅ 新方式（自动处理）

```typescript
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

const router = useI18nRouter();

// 自动添加语言前缀
router.push('/category/1');  // → /zh-CN/category/1
```

## 🎯 最佳实践

### 1. 统一使用 useI18nRouter

```typescript
// ✅ 推荐：在所有组件中使用
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

export function MyComponent() {
    const router = useI18nRouter();
    // ...
}
```

### 2. 避免混合使用

```typescript
// ❌ 不推荐：混合使用会导致不一致
import {useRouter} from 'next/navigation';
import {useI18nRouter} from '@/i18n/helpers/useI18nRouter';

export function MyComponent() {
    const nextRouter = useRouter();      // 不会自动添加 locale
    const i18nRouter = useI18nRouter();  // 会自动添加 locale
    
    // 容易导致混淆
}
```

### 3. 在 Server Component 中

```typescript
// Server Component 中使用 redirect
import {redirect} from 'next/navigation';

export default function Page() {
    // 直接重定向到带 locale 的路径
    redirect('/zh-CN/category/1');
}
```

## 📚 相关文件

- `useI18nRouter.ts` - Hook 实现
- `useI18n.ts` - 国际化 Hook
- `routing.ts` - 路由配置

## 🔗 参考资料

- [Next.js Router 文档](https://nextjs.org/docs/app/api-reference/functions/use-router)
- [next-intl 路由文档](https://next-intl-docs.vercel.app/docs/routing/navigation)
