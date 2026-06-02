'use client';

import {useState, useEffect, useMemo} from 'react';
import {useTranslations} from 'next-intl';

import {
    type contentservicev1_ListPostResponse,
    type contentservicev1_Post,
    identityservicev1_User
} from "@/api/generated/app/service/v1";
import {fetchListPosts, getPostTitle, getPostSummary} from "@/api/hooks/post";
import {fetchUserProfile} from "@/api/hooks/user-profile";

import {formatDateTime} from "@/utils";
import {useI18nRouter} from "@/i18n/helpers";
import UserStatGrid from '@/components/user/UserStatGrid';
import UserInfoSidebar, {type InfoSection, type InfoRow} from '@/components/user/UserInfoSidebar';
import MetaDataItem from '@/components/ui/meta-data-item';

export default function UserProfilePage() {
    const t = useTranslations('page.user');
    const router = useI18nRouter();

    const [loading, setLoading] = useState(false);
    const [postsLoading, setPostsLoading] = useState(false);
    const [activeTab, setActiveTab] = useState<'posts' | 'activities' | 'collections'>('posts');
    const [user, setUser] = useState<identityservicev1_User | null>(null);
    const [posts, setPosts] = useState<contentservicev1_Post[]>([]);
    const [postsTotal, setPostsTotal] = useState(0);

    // 统计数据
    const stats = useMemo(() => [
        {value: user?.following ?? 0, label: t('following')},
        {value: user?.followers ?? 0, label: t('followers')},
        {value: user?.postCount ?? 0, label: t('posts')},
        {value: user?.likeCount ?? 0, label: t('likes_received')},
    ], [user]);

    // 格式化性别
    const formatGender = (gender?: string) => {
        if (!gender) return t('gender_secret');
        const genderMap: Record<string, string> = {
            'SECRET': t('gender_secret'),
            'MALE': t('gender_male'),
            'FEMALE': t('gender_female'),
        };
        return genderMap[gender] || gender;
    };

    // 格式化状态
    const formatStatus = (status?: string) => {
        if (!status) return '-';
        const statusMap: Record<string, string> = {
            'NORMAL': t('status_normal'),
            'DISABLED': t('status_disabled'),
            'PENDING': t('status_pending'),
            'LOCKED': t('status_locked'),
            'EXPIRED': t('status_expired'),
            'CLOSED': t('status_closed'),
        };
        return statusMap[status] || status;
    };

    const statusColor: Record<string, string> = {
        'NORMAL': 'bg-emerald-50 text-emerald-700 border border-emerald-200/60 dark:bg-emerald-500/10 dark:text-emerald-400 dark:border-emerald-500/20',
        'DISABLED': 'bg-red-50 text-red-700 border border-red-200/60 dark:bg-red-500/10 dark:text-red-400 dark:border-red-500/20',
        'PENDING': 'bg-yellow-50 text-yellow-700 border border-yellow-200/60 dark:bg-yellow-500/10 dark:text-yellow-400 dark:border-yellow-500/20',
        'LOCKED': 'bg-orange-50 text-orange-700 border border-orange-200/60 dark:bg-orange-500/10 dark:text-orange-400 dark:border-orange-500/20',
        'EXPIRED': 'bg-gray-100 text-gray-600 border border-gray-200/60 dark:bg-gray-500/10 dark:text-gray-400 dark:border-gray-500/20',
        'CLOSED': 'bg-red-50 text-red-700 border border-red-200/60 dark:bg-red-500/10 dark:text-red-400 dark:border-red-500/20',
    };

    // 侧栏信息分区
    const sidebarSections: InfoSection[] = useMemo(() => {
        const sections: InfoSection[] = [];

        // 基础信息
        const basicRows = [];
        if (user?.username) basicRows.push({label: t('username'), value: user.username});
        if (user?.realname) basicRows.push({label: t('realname'), value: user.realname});
        if (user?.email) basicRows.push({label: t('email'), value: user.email});
        if (user?.mobile) basicRows.push({label: t('mobile'), value: user.mobile});
        if (basicRows.length > 0) sections.push({title: t('basic_info'), rows: basicRows});

        // 账号信息
        const accountRows: InfoRow[] = [{
            label: t('status'),
            value: (
                <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${statusColor[user?.status || ''] || 'bg-muted text-muted-foreground'}`}>
                    {formatStatus(user?.status)}
                </span>
            )
        }];
        if (user?.roleNames?.length) {
            accountRows.push({label: t('roles'), value: user.roleNames.join(', ')});
        }
        if (user?.createdAt) {
            accountRows.push({label: t('created_at'), value: formatDateTime(user.createdAt)});
        }
        if (user?.lastLoginAt) {
            accountRows.push({label: t('last_login_at'), value: formatDateTime(user.lastLoginAt)});
        }
        sections.push({title: t('account_info'), rows: accountRows});

        return sections;
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [user, t]);

    async function loadUserProfile() {
        setLoading(true);
        try {
            const result = await fetchUserProfile() as identityservicev1_User;
            setUser(result || null);
            if (activeTab === 'posts') {
                await loadUserPosts(result?.id);
            }
        } catch (error) {
            console.error('Load user profile failed:', error);
            setUser(null);
        } finally {
            setLoading(false);
        }
    }

    async function loadUserPosts(userId?: number) {
        const targetId = userId || user?.id;
        if (!targetId) return;

        setPostsLoading(true);
        try {
            const result = await fetchListPosts({
                paging: {page: 1, pageSize: 10},
                formValues: {author_id: targetId},
                fieldMask: undefined,
                orderBy: ['-createdAt']
            }) as unknown as contentservicev1_ListPostResponse;

            setPosts(result.items || []);
            setPostsTotal(result.total || 0);
        } catch (error) {
            console.error('Load user posts failed:', error);
        } finally {
            setPostsLoading(false);
        }
    }

    useEffect(() => {
        if (activeTab === 'posts' && posts.length === 0 && user) {
            loadUserPosts();
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [activeTab, user]);

    useEffect(() => {
        loadUserProfile();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, []);

    if (loading) {
        return (
            <div className="w-full">
                <div className="h-[200px] animate-pulse bg-gradient-to-r from-primary/10 to-primary/5"/>
                <div className="w-full max-w-[1200px] mx-auto grid grid-cols-[300px_1fr] gap-6 px-8 py-8 max-md:grid-cols-1 max-md:px-4">
                    <aside>
                        <div className="space-y-3 rounded-xl border border-border bg-card p-6">
                            <div className="h-4 w-3/4 animate-pulse rounded bg-muted"/>
                            <div className="h-4 w-1/2 animate-pulse rounded bg-muted"/>
                            <div className="h-4 w-2/3 animate-pulse rounded bg-muted"/>
                        </div>
                    </aside>
                    <main>
                        <div className="space-y-3 rounded-xl border border-border bg-card p-6">
                            <div className="h-4 w-full animate-pulse rounded bg-muted"/>
                            <div className="h-4 w-3/4 animate-pulse rounded bg-muted"/>
                        </div>
                    </main>
                </div>
            </div>
        );
    }

    if (!user) {
        return (
            <div className="flex w-full items-center justify-center py-20">
                <div className="text-lg text-muted-foreground">{t('please_login')}</div>
            </div>
        );
    }

    const genderColor: Record<string, string> = {
        'MALE': 'bg-blue-50 text-blue-700 border border-blue-200/60 dark:bg-blue-500/10 dark:text-blue-400 dark:border-blue-500/20',
        'FEMALE': 'bg-pink-50 text-pink-700 border border-pink-200/60 dark:bg-pink-500/10 dark:text-pink-400 dark:border-pink-500/20',
        'SECRET': 'bg-muted text-muted-foreground border border-border',
    };

    const tabBase = 'cursor-pointer border-b-2 px-3 py-2 text-sm font-medium transition-colors border-transparent text-muted-foreground hover:text-foreground';
    const tabActive = 'border-primary text-primary';

    return (
        <div className="w-full">
            {/* 顶部背景和基本信息 */}
            <div className="relative">
                <div className="h-[200px] bg-gradient-to-r from-primary/10 via-primary/5 to-primary/10"/>
                <div className="w-full max-w-[1200px] mx-auto px-8 max-md:px-4">
                    <div className="-mt-16 flex gap-8 max-md:flex-col max-md:items-center max-md:gap-4">
                        {/* 用户基本信息 */}
                        <div className="flex gap-6 max-md:flex-col max-md:items-center">
                            <div className="h-20 w-20 shrink-0 overflow-hidden rounded-full border-4 border-background bg-muted text-2xl font-bold text-muted-foreground ring-4 ring-primary/10 max-md:h-24 max-md:w-24">
                                {user?.avatar ? (
                                    <img src={user.avatar} alt={user.nickname} className="h-full w-full object-cover"/>
                                ) : (
                                    <div className="flex h-full w-full items-center justify-center">
                                        {(user?.nickname?.charAt(0) || user?.username?.charAt(0) || 'U')}
                                    </div>
                                )}
                            </div>

                            <div className="flex-1 py-2 max-md:text-center">
                                <div className="flex items-center gap-3 max-md:justify-center">
                                    <h1 className="text-2xl font-bold text-foreground">
                                        {user?.nickname || user?.username}
                                    </h1>
                                    {user?.gender && (
                                        <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${genderColor[user.gender] || genderColor['SECRET']}`}>
                                            {formatGender(user.gender)}
                                        </span>
                                    )}
                                </div>
                                {user?.description && (
                                    <p className="mt-2 text-sm text-muted-foreground">{user.description}</p>
                                )}
                                <div className="mt-2 flex gap-3 text-sm text-muted-foreground max-md:justify-center">
                                    {user?.region && (
                                        <span>{user.region.replace(/\d+/g, '').trim()}</span>
                                    )}
                                    {user?.tenantName && <span>{user.tenantName}</span>}
                                </div>
                            </div>
                        </div>

                        <div className="ml-auto py-2 max-md:ml-0">
                            <button
                                className="cursor-pointer rounded-lg border border-border bg-card px-4 py-2 text-sm font-medium text-foreground transition-all hover:border-primary hover:bg-primary/5"
                                onClick={() => router.push('/settings')}
                            >
                                {t('edit_profile')}
                            </button>
                        </div>
                    </div>

                    {/* 统计数据 */}
                    <UserStatGrid stats={stats}/>
                </div>
            </div>

            {/* 内容区域 */}
            <div className="w-full max-w-[1200px] mx-auto grid grid-cols-[300px_1fr] gap-6 px-8 py-8 max-md:grid-cols-1 max-md:px-4">
                {/* 左侧面板 */}
                <UserInfoSidebar sections={sidebarSections}/>

                {/* 主内容区 */}
                <main className="min-w-0">
                    <div className="rounded-xl border border-border bg-card">
                        <div className="flex border-b border-border">
                            <button
                                className={`${tabBase} ${activeTab === 'posts' ? tabActive : ''}`}
                                onClick={() => setActiveTab('posts')}
                            >
                                {t('tab_posts')}
                            </button>
                            <button
                                className={`${tabBase} ${activeTab === 'activities' ? tabActive : ''}`}
                                onClick={() => setActiveTab('activities')}
                            >
                                {t('tab_activities')}
                            </button>
                            <button
                                className={`${tabBase} ${activeTab === 'collections' ? tabActive : ''}`}
                                onClick={() => setActiveTab('collections')}
                            >
                                {t('tab_collections')}
                            </button>
                        </div>

                        <div className="p-6">
                            {activeTab === 'posts' && (
                                <>
                                    {postsLoading ? (
                                        <div className="space-y-3">
                                            <div className="h-16 w-full animate-pulse rounded bg-muted"/>
                                            <div className="h-16 w-full animate-pulse rounded bg-muted"/>
                                        </div>
                                    ) : posts.length > 0 ? (
                                        <div className="space-y-4">
                                            {posts.map((post) => (
                                                <div key={post.id} className="flex items-start justify-between gap-4 rounded-lg border border-border bg-background p-4 transition-all hover:border-primary hover:shadow-sm">
                                                    <div className="flex-1 min-w-0">
                                                        <h3 className="mb-1 font-semibold text-foreground">{getPostTitle(post)}</h3>
                                                        <p className="mb-2 line-clamp-2 text-sm text-muted-foreground">{getPostSummary(post)}</p>
                                                        <div className="flex flex-wrap items-center gap-3 text-xs text-muted-foreground">
                                                            <MetaDataItem icon="carbon:view" text={`${post.visits || 0} ${t('views')}`}/>
                                                            <MetaDataItem icon="carbon:thumbs-up" text={`${post.likes || 0} ${t('likes')}`}/>
                                                            <MetaDataItem icon="carbon:chat" text={`${post.commentCount || 0} ${t('comments')}`}/>
                                                            <MetaDataItem icon="carbon:time" text={formatDateTime(post.createdAt)}/>
                                                        </div>
                                                    </div>
                                                    <button className="shrink-0 cursor-pointer text-sm text-primary transition-colors hover:text-primary/80 hover:underline">
                                                        {t('view_post')} →
                                                    </button>
                                                </div>
                                            ))}
                                        </div>
                                    ) : (
                                        <EmptyTab icon="📄" text={t('no_posts')}/>
                                    )}
                                </>
                            )}

                            {activeTab === 'activities' && (
                                <EmptyTab icon="🎯" text={t('no_activities')}/>
                            )}

                            {activeTab === 'collections' && (
                                <EmptyTab icon="🔖" text={t('no_collections')}/>
                            )}
                        </div>
                    </div>
                </main>
            </div>
        </div>
    );
}

/** 空状态占位 */
function EmptyTab({icon, text}: {icon: string; text: string}) {
    return (
        <div className="py-12 text-center">
            <div className="mb-2 text-4xl">{icon}</div>
            <div className="text-sm text-muted-foreground">{text}</div>
        </div>
    );
}
