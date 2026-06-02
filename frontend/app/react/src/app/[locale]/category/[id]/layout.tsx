// Passthrough layout — [id] 动态路由由 page.tsx 的 generateStaticParams 处理
export default function CategoryDetailLayout({children}: { children: React.ReactNode }) {
    return children;
}
