// Server component layout for [id] dynamic segment
// ID 来自 CMS API，构建时无法预知，返回空数组允许运行时按需生成
export const dynamicParams = true;

export function generateStaticParams() {
    return [];
}

export default function CategoryDetailLayout({children}: { children: React.ReactNode }) {
    return children;
}
