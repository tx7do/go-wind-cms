import {useState, useEffect, useMemo} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text, Image} from '@tarojs/components';

import XIcon from '@/plugins/xicon';
import {formatDateTime} from "@/utils";

import {
  type contentservicev1_ListPostResponse,
  type contentservicev1_Post,
  identityservicev1_User
} from "@/api/generated/app/service/v1";
import {fetchListPosts, getPostTitle, getPostSummary} from "@/api/hooks/post";
import {fetchUserProfile} from "@/api/hooks/user-profile";

import './user.scss';

export default function UserProfilePage() {
  const {t} = useTranslation();

  const [loading, setLoading] = useState(false);
  const [postsLoading, setPostsLoading] = useState(false);

  const [activeTab, setActiveTab] = useState<'posts' | 'activities' | 'collections'>('posts');

  const [user, setUser] = useState<identityservicev1_User | null>(null);
  const [posts, setPosts] = useState<contentservicev1_Post[]>([]);

  // 统计数据
  const stats = useMemo(() => ({
    followers: user?.followers ?? 0,
    following: user?.following ?? 0,
    posts: user?.postCount ?? 0,
    likes: user?.likeCount ?? 0,
  }), [user]);

  // 格式化性别
  const formatGender = (gender?: string) => {
    if (!gender) return t('page.user.gender_secret');
    const genderMap: Record<string, string> = {
      'SECRET': t('page.user.gender_secret'),
      'MALE': t('page.user.gender_male'),
      'FEMALE': t('page.user.gender_female'),
    };
    return genderMap[gender] || gender;
  };

  // 格式化状态
  const formatStatus = (status?: string) => {
    if (!status) return '-';
    const statusMap: Record<string, string> = {
      'NORMAL': t('page.status_normal'),
      'DISABLED': t('page.status_disabled'),
      'PENDING': t('page.status_pending'),
      'LOCKED': t('page.status_locked'),
      'EXPIRED': t('page.status_expired'),
      'CLOSED': t('page.status_closed'),
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
      }) as contentservicev1_ListPostResponse;

      setPosts(result.items || []);
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
      <View className='user-profile-page'>
        <View className='profile-header'>
          <View className='header-skeleton'></View>
        </View>
        <View className='profile-main'>
          <View className='profile-sidebar'>
            <View className='info-card'>
              <View className='skeleton-text'></View>
              <View className='skeleton-text'></View>
              <View className='skeleton-text'></View>
            </View>
          </View>
          <View className='profile-content-area'>
            <View className='content-card'>
              <View className='skeleton-text'></View>
              <View className='skeleton-text'></View>
            </View>
          </View>
        </View>
      </View>
    );
  }

  // 用户未登录
  if (!user) {
    return (
      <View className='empty-state'>
        <Text className='empty-text'>{t('page.user.please_login')}</Text>
      </View>
    );
  }

  return (
    <View className='user-profile-page'>
      {/* 顶部背景和基本信息 */}
      <View className='profile-header'>
        <View className='header-bg'></View>
        <View className='profile-content'>
          {/* 用户基本信息 */}
          <View className='user-basic-info'>
            <View className='avatar-section'>
              <View className='avatar-wrapper'>
                <View className='user-avatar'>
                  {user?.avatar ? (
                    <Image src={user.avatar} />
                  ) : (
                    <Text>{user?.nickname?.charAt(0) || user?.username?.charAt(0) || 'U'}</Text>
                  )}
                </View>
              </View>
            </View>

            <View className='info-section'>
              <View className='user-header'>
                <View className='user-title'>
                  <Text className='user-name'>{user?.nickname || user?.username}</Text>
                  {user?.gender && (
                    <Text className={`gender-tag ${user.gender.toLowerCase()}`}>
                      {formatGender(user.gender)}
                    </Text>
                  )}
                </View>
                <View className='edit-btn'>{t('page.user.edit_profile')}</View>
              </View>

              {user?.description && (
                <Text className='user-bio'>{user.description}</Text>
              )}

              <View className='user-meta'>
                {user?.region && (
                  <Text className='meta-item'>
                    {user.region.replace(/\d+/g, '').trim()}
                  </Text>
                )}
                {user?.tenantName && <Text className='meta-item'>{user.tenantName}</Text>}
              </View>
            </View>
          </View>
        </View>

        {/* 统计数据 */}
        <View className='user-stats'>
          <View className='stat-item'>
            <Text className='stat-value'>{stats.following}</Text>
            <Text className='stat-label'>{t('page.user.following')}</Text>
          </View>
          <View className='stat-item'>
            <Text className='stat-value'>{stats.followers}</Text>
            <Text className='stat-label'>{t('page.user.followers')}</Text>
          </View>
          <View className='stat-item'>
            <Text className='stat-value'>{stats.posts}</Text>
            <Text className='stat-label'>{t('page.user.posts')}</Text>
          </View>
          <View className='stat-item'>
            <Text className='stat-value'>{stats.likes}</Text>
            <Text className='stat-label'>{t('page.user.likes_received')}</Text>
          </View>
        </View>
      </View>

      {/* 内容区域 */}
      <View className='profile-main'>
        {/* 左侧面板 */}
        <View className='profile-sidebar'>
          <View className='info-card'>
            <View className='card-section'>
              <Text className='section-title'>{t('page.user.basic_info')}</Text>
              <View className='info-list'>
                {user?.username && (
                  <View className='info-row'>
                    <Text className='info-key'>{t('page.user.username')}:</Text>
                    <Text className='info-val'>{user.username}</Text>
                  </View>
                )}
                {user?.realname && (
                  <View className='info-row'>
                    <Text className='info-key'>{t('page.user.realname')}:</Text>
                    <Text className='info-val'>{user.realname}</Text>
                  </View>
                )}
                {user?.email && (
                  <View className='info-row'>
                    <Text className='info-key'>{t('page.user.email')}:</Text>
                    <Text className='info-val'>{user.email}</Text>
                  </View>
                )}
                {user?.mobile && (
                  <View className='info-row'>
                    <Text className='info-key'>{t('page.user.mobile')}:</Text>
                    <Text className='info-val'>{user.mobile}</Text>
                  </View>
                )}
              </View>
            </View>

            <View className='divider'></View>

            <View className='card-section'>
              <Text className='section-title'>{t('page.user.account_info')}</Text>
              <View className='info-list'>
                <View className='info-row'>
                  <Text className='info-key'>{t('page.user.status')}:</Text>
                  <Text className={`status-tag status${user?.status}`}>
                    {formatStatus(user?.status)}
                  </Text>
                </View>
                {user?.roleNames?.length && (
                  <View className='info-row'>
                    <Text className='info-key'>{t('page.user.roles')}:</Text>
                    <Text className='info-val'>{user.roleNames.join(', ')}</Text>
                  </View>
                )}
                {user?.createdAt && (
                  <View className='info-row'>
                    <Text className='info-key'>{t('page.user.created_at')}:</Text>
                    <Text className='info-val'>{formatDateTime(user.createdAt)}</Text>
                  </View>
                )}
                {user?.lastLoginAt && (
                  <View className='info-row'>
                    <Text className='info-key'>{t('page.user.last_login_at')}:</Text>
                    <Text className='info-val'>{formatDateTime(user.lastLoginAt)}</Text>
                  </View>
                )}
              </View>
            </View>
          </View>
        </View>

        {/* 主内容区 */}
        <View className='profile-content-area'>
          <View className='content-card'>
            <View className='tabs'>
              <View
                className={`tab ${activeTab === 'posts' ? 'active' : ''}`}
                onClick={() => setActiveTab('posts')}
              >
                <Text>{t('page.user.tab_posts')}</Text>
              </View>
              <View
                className={`tab ${activeTab === 'activities' ? 'active' : ''}`}
                onClick={() => setActiveTab('activities')}
              >
                <Text>{t('page.user.tab_activities')}</Text>
              </View>
              <View
                className={`tab ${activeTab === 'collections' ? 'active' : ''}`}
                onClick={() => setActiveTab('collections')}
              >
                <Text>{t('page.user.tab_collections')}</Text>
              </View>
            </View>

            <View className='tab-content'>
              {activeTab === 'posts' && (
                <>
                  {postsLoading ? (
                    <View className='posts-list-loading'>
                      <View className='skeleton-text'></View>
                      <View className='skeleton-text'></View>
                    </View>
                  ) : posts.length > 0 ? (
                    <View className='posts-list'>
                      {posts.map((post) => (
                        <View key={post.id} className='post-item'>
                          <View className='post-content'>
                            <Text className='post-title'>{getPostTitle(post)}</Text>
                            <Text className='post-summary'>{getPostSummary(post)}</Text>
                            <View className='post-meta'>
                              <View className='meta-info'>
                                <XIcon name='carbon:view' size={16} />
                                <Text>{post.visits || 0} {t('page.views')}</Text>
                              </View>
                              <View className='meta-info'>
                                <XIcon name='carbon:thumb-up' size={16} />
                                <Text>{post.likes || 0} {t('page.likes')}</Text>
                              </View>
                              <View className='meta-info'>
                                <XIcon name='carbon:chat' size={16} />
                                <Text>{post.commentCount || 0} {t('page.comments')}</Text>
                              </View>
                              <View className='meta-info'>
                                <XIcon name='carbon:time' size={16} />
                                <Text>{formatDateTime(post.createdAt)}</Text>
                              </View>
                            </View>
                          </View>
                          <View className='view-post-btn'>
                            <Text>{t('page.view_post')} →</Text>
                          </View>
                        </View>
                      ))}
                    </View>
                  ) : (
                    <View className='empty'>
                      <Text className='empty-icon'>
                        <XIcon name='carbon:document' size={48} />
                      </Text>
                      <Text className='empty-text'>{t('page.user.no_posts')}</Text>
                    </View>
                  )}
                </>
              )}

              {activeTab === 'activities' && (
                <View className='empty'>
                  <Text className='empty-icon'>
                    <XIcon name='carbon:target' size={48} />
                  </Text>
                  <Text className='empty-text'>{t('page.user.no_activities')}</Text>
                </View>
              )}

              {activeTab === 'collections' && (
                <View className='empty'>
                  <Text className='empty-icon'>
                    <XIcon name='carbon:bookmark' size={48} />
                  </Text>
                  <Text className='empty-text'>{t('page.user.no_collections')}</Text>
                </View>
              )}
            </View>
          </View>
        </View>
      </View>
    </View>
  );
}
