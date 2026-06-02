'use client';

import {useState, useEffect, useMemo} from 'react';
import {useTranslations} from 'next-intl';

import XIcon from '@/plugins/xicon';

import {
    type contentservicev1_ListPostResponse,
    type contentservicev1_Post,
    identityservicev1_User
} from "@/api/generated/app/service/v1";
import {fetchListPosts, getPostTitle, getPostSummary} from "@/api/hooks/post";
import {fetchUserProfile} from "@/api/hooks/user-profile";

import '../../globals.css';
import {formatDateTime} from "@/utils";

export default function UserProfilePage() {
    const t = useTranslations('page.user');

    const [loading, setLoading] = useState(false);
    const [postsLoading, setPostsLoading] = useState(false);

    const [activeTab, setActiveTab] = useState<'posts' | 'activities' | 'collections'>('posts');

    const [user, setUser] = useState<identityservicev1_User | null>(null);
    const [posts, setPosts] = useState<contentservicev1_Post[]>([]);
    const [postsTotal, setPostsTotal] = useState(0);

    // 统计数据
    const stats = useMemo(() => ({
        followers: user?.followers ?? 0,
        following: user?.following ?? 0,
        posts: user?.postCount ?? 0,
        likes: user?.likeCount ?? 0,
    }), [user]);

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

    // 获取用户信息
    async function loadUserProfile() {
        setLoading(true);
        try {
            const result = await fetchUserProfile() as identityservicev1_User;
            setUser(result || null);

            if (activeTab === 'posts') {
                await loadUserPosts();
            }
        } catch (error) {
            console.error('Load user profile failed:', error);
            setUser(null);
        } finally {
            setLoading(false);
        }
    }

    // 加载用户的帖子列表
    async function loadUserPosts() {
        if (!user?.id) return;

        setPostsLoading(true);
        try {
            const result = await fetchListPosts({
                paging: {page: 1, pageSize: 10},
                formValues: {author_id: user.id},
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

    // 监听 activeTab 变化
    useEffect(() => {
        if (activeTab === 'posts' && posts.length === 0 && user) {
            loadUserPosts();
        }
    }, [activeTab, user]);

    // 初始加载
    useEffect(() => {
        loadUserProfile();
    }, []);

    // 渲染 Loading Skeleton
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

    // 用户未登录
    if (!user) {
        return (
            <div className="flex w-full items-center justify-center py-20">
                <div className="text-lg text-muted-foreground">{t('please_login')}</div>
            </div>
        );
    }

    const genderColor: Record<string, string> = {
        'MALE': 'bg-blue-500/10 text-blue-600',
        'FEMALE': 'bg-pink-500/10 text-pink-600',
        'SECRET': 'bg-muted text-muted-foreground',
    };

    const statusColor: Record<string, string> = {
        'NORMAL': 'bg-green-500/10 text-green-600',
        'DISABLED': 'bg-red-500/10 text-red-600',
        'PENDING': 'bg-yellow-500/10 text-yellow-600',
        'LOCKED': 'bg-orange-500/10 text-orange-600',
        'EXPIRED': 'bg-gray-500/10 text-gray-600',
        'CLOSED': 'bg-red-500/10 text-red-600',
    };

    const tabBase = 'cursor-pointer border-b-2 px-4 py-2 text-sm font-medium transition-colors border-transparent text-muted-foreground hover:text-foreground';
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
                            <button className="cursor-pointer rounded-lg border border-border bg-card px-4 py-2 text-sm font-medium text-foreground transition-all hover:border-primary hover:bg-primary/5">
                                {t('edit_profile')}
                            </button>
                        </div>
                    </div>

                    {/* 统计数据 */}
                    <div className="mt-6 grid grid-cols-4 gap-4 border-t border-border pt-6 max-md:grid-cols-2">
                        <div className="text-center">
                            <div className="text-2xl font-bold text-foreground">{stats.following}</div>
                            <div className="text-xs text-muted-foreground">{t('following')}</div>
                        </div>
                        <div className="text-center">
                            <div className="text-2xl font-bold text-foreground">{stats.followers}</div>
                            <div className="text-xs text-muted-foreground">{t('followers')}</div>
                        </div>
                        <div className="text-center">
                            <div className="text-2xl font-bold text-foreground">{stats.posts}</div>
                            <div className="text-xs text-muted-foreground">{t('posts')}</div>
                        </div>
                        <div className="text-center">
                            <div className="text-2xl font-bold text-foreground">{stats.likes}</div>
                            <div className="text-xs text-muted-foreground">{t('likes_received')}</div>
                        </div>
                    </div>
                </div>
            </div>

            {/* 内容区域 */}
            <div className="w-full max-w-[1200px] mx-auto grid grid-cols-[300px_1fr] gap-6 px-8 py-8 max-md:grid-cols-1 max-md:px-4">
                {/* 左侧面板 */}
                <aside>
                    <div className="rounded-xl border border-border bg-card p-6">
                        <div>
                            <h3 className="mb-3 text-sm font-semibold text-foreground">{t('basic_info')}</h3>
                            <div className="space-y-2">
                                {user?.username && (
                                    <div className="flex justify-between text-sm">
                                        <span className="text-muted-foreground">{t('username')}:</span>
                                        <span className="text-foreground">{user.username}</span>
                                    </div>
                                )}
                                {user?.realname && (
                                    <div className="flex justify-between text-sm">
                                        <span className="text-muted-foreground">{t('realname')}:</span>
                                        <span className="text-foreground">{user.realname}</span>
                                    </div>
                                )}
                                {user?.email && (
                                    <div className="flex justify-between text-sm">
                                        <span className="text-muted-foreground">{t('email')}:</span>
                                        <span className="text-foreground">{user.email}</span>
                                    </div>
                                )}
                                {user?.mobile && (
                                    <div className="flex justify-between text-sm">
                                        <span className="text-muted-foreground">{t('mobile')}:</span>
                                        <span className="text-foreground">{user.mobile}</span>
                                    </div>
                                )}
                            </div>
                        </div>

                        <div className="my-4 border-t border-border"/>

                        <div>
                            <h3 className="mb-3 text-sm font-semibold text-foreground">{t('account_info')}</h3>
                            <div className="space-y-2">
                                <div className="flex justify-between text-sm">
                                    <span className="text-muted-foreground">{t('status')}:</span>
                                    <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${statusColor[user?.status || ''] || 'bg-muted text-muted-foreground'}`}>
                                        {formatStatus(user?.status)}
                                    </span>
                                </div>
                                {user?.roleNames?.length && (
                                    <div className="flex justify-between text-sm">
                                        <span className="text-muted-foreground">{t('roles')}:</span>
                                        <span className="text-foreground">{user.roleNames.join(', ')}</span>
                                    </div>
                                )}
                                {user?.createdAt && (
                                    <div className="flex justify-between text-sm">
                                        <span className="text-muted-foreground">{t('created_at')}:</span>
                                        <span className="text-foreground">{formatDateTime(user.createdAt)}</span>
                                    </div>
                                )}
                                {user?.lastLoginAt && (
                                    <div className="flex justify-between text-sm">
                                        <span className="text-muted-foreground">{t('last_login_at')}:</span>
                                        <span className="text-foreground">{formatDateTime(user.lastLoginAt)}</span>
                                    </div>
                                )}
                            </div>
                        </div>
                    </div>
                </aside>

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
                                                            <span className="flex items-center gap-1">
                                                                <XIcon name="carbon:view" size={16}/>
                                                                {post.visits || 0} {t('views')}
                                                            </span>
                                                            <span className="flex items-center gap-1">
                                                                <XIcon name="carbon:thumbs-up" size={16}/>
                                                                {post.likes || 0} {t('likes')}
                                                            </span>
                                                            <span className="flex items-center gap-1">
                                                                <XIcon name="carbon:chat" size={16}/>
                                                                {post.commentCount || 0} {t('comments')}
                                                            </span>
                                                            <span className="flex items-center gap-1">
                                                                <XIcon name="carbon:time" size={16}/>
                                                                {formatDateTime(post.createdAt)}
                                                            </span>
                                                        </div>
                                                    </div>
                                                    <button className="shrink-0 cursor-pointer text-sm text-primary transition-colors hover:text-primary/80 hover:underline">
                                                        {t('view_post')} →
                                                    </button>
                                                </div>
                                            ))}
                                        </div>
                                    ) : (
                                        <div className="py-12 text-center">
                                            <div className="mb-2 text-4xl">📄</div>
                                            <div className="text-sm text-muted-foreground">{t('no_posts')}</div>
                                        </div>
                                    )}
                                </>
                            )}

                            {activeTab === 'activities' && (
                                <div className="py-12 text-center">
                                    <div className="mb-2 text-4xl">🎯</div>
                                    <div className="text-sm text-muted-foreground">{t('no_activities')}</div>
                                </div>
                            )}

                            {activeTab === 'collections' && (
                                <div className="py-12 text-center">
                                    <div className="mb-2 text-4xl">🔖</div>
                                    <div className="text-sm text-muted-foreground">{t('no_collections')}</div>
                                </div>
                            )}
                        </div>
                    </div>
                </main>
            </div>
        </div>
    );
}
