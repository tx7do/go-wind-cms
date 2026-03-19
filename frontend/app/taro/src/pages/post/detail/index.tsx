import {useEffect, useState, useRef, useMemo, useCallback} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text, Image} from '@tarojs/components';
import Taro from '@tarojs/taro';

import CommentSection from '@/components/comment/CommentSection';
import ContentViewer from '@/components/content/ContentViewer';
import PostList from '@/components/post/PostList';
import BackToTop from '@/components/layout/BackToTop';
import XIcon from '@/plugins/xicon';

import {usePostStore} from '@/store/slices/post/hooks';
import {formatDate} from "@/utils";
import {contentservicev1_Post} from "@/api/generated/app/service/v1";

import './post-detail.scss';

// 常量定义
const HEADING_OFFSET = 150;
const THROTTLE_DELAY = 200;

interface TocItem {
  id: string;
  level: number;
  text: string;
  element: Element;
}

export default function PostDetailPage() {
  const {t} = useTranslation();
  const postStore = usePostStore();

  const post = postStore.detail as contentservicev1_Post | null;

  const [localLoading, setLocalLoading] = useState(true);
  const isLoading = localLoading || (postStore.loading && !post);
  const [tableOfContents, setTableOfContents] = useState<TocItem[]>([]);
  const [activeHeading, setActiveHeading] = useState('');
  const [isLiked, setIsLiked] = useState(false);
  const [isBookmarked, setIsBookmarked] = useState(false);
  const [isTocExpanded, setIsTocExpanded] = useState(true);
  const [isNarrowScreen, setIsNarrowScreen] = useState(false);

  const contentRef = useRef<HTMLDivElement>(null);

  const postId = useMemo(() => {
    let id: string | null | undefined = null;
    if (typeof Taro.getCurrentInstance === 'function') {
      const instance = Taro.getCurrentInstance();
      const routeId = instance?.router?.params?.id;
      id = typeof routeId === 'string' ? routeId : null;
    } else {
      // 兼容 getCurrentPages
      const pages = Taro.getCurrentPages();
      const currentPage = pages[pages.length - 1];
      const pageId = (currentPage as any).options?.id;
      id = typeof pageId === 'string' ? pageId : null;
    }
    return id ? parseInt(id) : null;
  }, []);

  const displayTitle = useMemo(() => {
    if (!post) return '';
    return postStore.getPostTitle(post);
  }, [post, postStore]);

  const displayContent = useMemo(() => {
    if (!post) return '';
    return postStore.getPostContent(post);
  }, [post, postStore]);

  const displayThumbnail = useMemo(() => {
    if (!post) return '';
    return postStore.getPostThumbnail(post);
  }, [post, postStore]);

  // 相关文章数据
  const relatedPostsQuery = useMemo(() => {
    if (!post?.categoryIds) return null;
    return {
      status: 'POST_STATUS_PUBLISHED',
      id__not: postId,
      category_ids__in: post.categoryIds
    };
  }, [post?.categoryIds, postId]);

  // Load post data
  useEffect(() => {
    async function loadPost() {
      if (!postId) return;

      setLocalLoading(true);
      try {
        console.log('[PostDetail] Before getPost call:', {
          postId,
          currentDetail: postStore.detail
        });

        const fetchedPost = (await postStore.getPost({
          id: postId
        })) as contentservicev1_Post;

        console.log('[PostDetail] API returned:', {
          hasFetchedPost: !!fetchedPost,
          hasTranslations: !!fetchedPost?.translations?.length,
          firstTranslationTitle: fetchedPost?.translations?.[0]?.title,
          firstTranslationContentLength: fetchedPost?.translations?.[0]?.content?.length
        });

        if (fetchedPost) {
          // TODO: Taro 中设置页面
          // document.title = `${postStore.getPostTitle(fetchedPost)} - GoWind Content Hub`;
        }
      } catch (error) {
        console.error('Load post failed:', error);
      } finally {
        setLocalLoading(false);
      }
    }

    loadPost();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [postId]);

  // 检测屏幕宽度，决定是否显示 TOC
  useEffect(() => {
    const checkScreenWidth = () => {
      // 小于 768px 认为是窄屏幕（平板和手机）
      const narrow = window.innerWidth < 768;
      setIsNarrowScreen(narrow);
      // 窄屏幕时自动收起 TOC
      if (narrow) {
        setIsTocExpanded(false);
      }
    };

    // 初始检査
    checkScreenWidth();

    // 监听窗口大小变化
    window.addEventListener('resize', checkScreenWidth);
    return () => window.removeEventListener('resize', checkScreenWidth);
  }, []);

  // Generate table of contents - 在内容渲染后生成目录
  useEffect(() => {
    if (!displayContent || !contentRef.current) {
      console.log('[TOC] Skip generation:', {hasContent: !!displayContent, hasRef: !!contentRef.current});
      return;
    }

    const timeoutId = setTimeout(() => {
      const contentEl = contentRef.current;
      if (!contentEl) {
        console.log('[TOC] No content element');
        return;
      }

      // 等待 RichText 渲染完成
      setTimeout(() => {
        const headings = contentEl.querySelectorAll('h2, h3');
        const toc: TocItem[] = [];

        console.log('[TOC] Found headings:', headings.length);

        headings.forEach((heading, index) => {
          const level = heading.tagName === 'H2' ? 2 : 3;
          const id = `${heading.tagName.toLowerCase()}-${index}`;
          if (!heading.id) heading.setAttribute('id', id);

          toc.push({
            id,
            level,
            text: heading.textContent || '',
            element: heading
          });
        });

        setTableOfContents(toc);
        console.log('[TOC] Generated:', toc.length, 'items');
      }, 100);
    }, 200); // 等待内容完全渲染

    return () => clearTimeout(timeoutId);
  }, [displayContent]); // 依赖计算后的内容

  useEffect(() => {
    if (displayContent && tableOfContents.length === 0) {
      generateTableOfContents();
    }
  }, [displayContent, tableOfContents.length]);

  // 生成目录函数
  const generateTableOfContents = useCallback(() => {
    setTimeout(() => {
      if (!contentRef.current) return;

      // 等待 RichText 渲染
      setTimeout(() => {
        const contentEl = contentRef.current;
        const headings = contentEl ? contentEl.querySelectorAll('h2, h3') : null;
        const toc: TocItem[] = [];

        console.log('[GenerateTOC] Found headings:', headings?.length);

        headings?.forEach((heading, index) => {
          const level = heading.tagName === 'H2' ? 2 : 3;
          // 保证 id 唯一，使用tagName + index
          const id = `${heading.tagName.toLowerCase()}-${index}`;

          // 确保 ID 存在
          if (!heading.id) {
            heading.setAttribute('id', id);
            console.log(`[GenerateTOC] Set ID "${id}" on ${heading.tagName}:`, heading.textContent);
          }

          toc.push({
            id,
            level,
            text: heading.textContent || '',
            element: heading
          });
        });

        setTableOfContents(toc);
        console.log('[GenerateTOC] Generated TOC:', toc);
      }, 100);
    }, 300);
  }, []);

  // 防抖函数
  const throttle = useCallback((fn: (...args: unknown[]) => void, delay: number) => {
    let timer: ReturnType<typeof setTimeout> | null = null;
    return function (this: unknown, ...args: unknown[]) {
      if (!timer) {
        timer = setTimeout(() => {
          fn.apply(this as (...args: unknown[]) => void, args as unknown[]);
          timer = null;
        }, delay);
      }
    };
  }, []);

  // 滚动处理函数
  const handleScroll = useCallback(() => {
    if (tableOfContents.length > 0) {
      let currentActive = activeHeading;
      for (const heading of tableOfContents) {
        if (heading.element && heading.element.getBoundingClientRect().top < HEADING_OFFSET) {
          currentActive = heading.id;
        }
      }
      setActiveHeading(currentActive);
    }
  }, [tableOfContents, activeHeading]);

  const throttledScroll = useMemo(() => {
    return throttle(handleScroll, THROTTLE_DELAY);
  }, [throttle, handleScroll]);

  // Scroll handler
  useEffect(() => {
    window.addEventListener('scroll', throttledScroll);
    return () => window.removeEventListener('scroll', throttledScroll);
  }, [throttledScroll]);

  // Handlers
  const handleBack = () => {
    // TODO: Taro.getCurrentPages() 获取参数
    Taro.navigateBack();
  };

  const handleLike = () => {
    setIsLiked(!isLiked);
    // TODO: Add toast message
  };

  const handleBookmark = () => {
    setIsBookmarked(!isBookmarked);
    // TODO: Add toast message
  };

  const handleShare = () => {
    // TODO: Taro 中使用onShareAppMessage 配置分享
    console.log('Share:', displayTitle);
  };

  const scrollToHeading = (id: string) => {
    const tocItem = tableOfContents.find(item => item.id === id);
    const element = tocItem?.element || document.getElementById(id);

    if (element) {
      console.log('[TOC Scroll] Start',
        'ID:', id,
        '| Element:', element.textContent?.trim()
      );

      // Taro 中使用pageScrollTo
      Taro.pageScrollTo({
        selector: `#${id}`,
        duration: 300
      }).then(() => {
        setActiveHeading(id);
      }).catch((err) => {
        console.error('Scroll failed:', err);
      });
    } else {
      console.error('[ScrollToHeading] Element not found:', id);
    }
  };

  if (isLoading) {
    return (
      <View className="post-detail-page">
        {/* Loading skeleton */}
        <View className="back-navigation">
          <View className="skeleton-btn"></View>
        </View>
        <View className="post-article">
          <View className="post-banner">
            <View className="skeleton-image"></View>
          </View>
          <View className="post-wrapper">
            <View className="toc-sidebar">
              <View className="toc-container">
                <View className="skeleton-line" style={{width: '200px', height: '24px'}}></View>
                <View className="skeleton-line" style={{width: '180px', height: '20px', marginTop: '16px'}}></View>
                <View className="skeleton-line" style={{width: '160px', height: '20px', marginTop: '8px'}}></View>
                <View className="skeleton-line" style={{width: '140px', height: '20px', marginTop: '8px'}}></View>
              </View>
            </View>
            <View className="article-content">
              <View className="post-header">
                <View className="skeleton-title" style={{width: '80%', height: '48px'}}></View>
                <View className="skeleton-subtitle" style={{width: '60%', height: '32px', marginTop: '16px'}}></View>
                <View className="post-meta">
                  <View className="skeleton-meta" style={{width: '100px', height: '20px'}}></View>
                  <View className="skeleton-meta" style={{width: '100px', height: '20px'}}></View>
                  <View className="skeleton-meta" style={{width: '100px', height: '20px'}}></View>
                  <View className="skeleton-meta" style={{width: '100px', height: '20px'}}></View>
                </View>
              </View>
              <View className="post-content">
                <View className="skeleton-paragraph"></View>
              </View>
            </View>
          </View>
        </View>
      </View>
    );
  }

  if (!post) {
    return (
      <View className="post-detail-page">
        <View className="empty-state">
          <XIcon name='carbon:warning' size={28} className='empty-icon'
                 style={{display: 'inline-block', width: '28px', height: '28px', fontSize: '28px'}}/>
          <Text>{t('page.post_detail.not_found') || 'Post not found'}</Text>
        </View>
      </View>
    );
  }

  return (
    <View className="post-detail-page">
      {/* Back Navigation */}
      <View className="back-navigation">
        <View onClick={handleBack} className="back-btn" aria-label={t('page.post_detail.back')}>
          <XIcon name='carbon:arrow-left' size={16}
                 style={{display: 'inline-block', width: '16px', height: '16px', fontSize: '16px'}}/>
          <Text>{t('page.post_detail.back')}</Text>
        </View>
      </View>

      {/* Post Article */}
      <View className="post-article">
        {/* Article Content Wrapper */}
        <View className="article-wrapper">
          {/* Post Thumbnail Banner */}
          {displayThumbnail && (
            <View className="post-banner">
              <Image src={displayThumbnail}/>
              <View className="banner-overlay"/>
            </View>
          )}

          {/* Content with TOC */}
          <View className="content-with-toc">
            {/* Left: Table of Contents - 窄屏幕不显示 */}
            {!isNarrowScreen && tableOfContents.length > 0 && isTocExpanded && (
              <View className="toc-sidebar">
                <View className="toc-container">
                  <View className="toc-header">
                    <Text className="toc-title">
                      <XIcon name='carbon:list' size={18}
                             style={{display: 'inline-block', width: '18px', height: '18px', fontSize: '18px'}}/>
                      <Text>{t('page.post_detail.table_of_contents')}</Text>
                    </Text>
                    <View
                      onClick={() => setIsTocExpanded(false)}
                      className="toc-collapse-btn"
                    >
                      <XIcon name='carbon:chevron-left' size={18}
                             style={{display: 'inline-block', width: '18px', height: '18px', fontSize: '18px'}}/>
                    </View>
                  </View>
                  <View className="toc-list">
                    {tableOfContents.map(item => (
                      <View
                        key={item.id}
                        className={`toc-item level-${item.level}${activeHeading === item.id ? ' active' : ''}`}
                        onClick={() => {
                          setTimeout(() => scrollToHeading(item.id), 10);
                        }}
                      >
                        <Text>{item.text}</Text>
                      </View>
                    ))}
                  </View>
                </View>
              </View>
            )}

            {/* Right: Article Content */}
            <View className="article-content-inner">
              {/* TOC Expand Button - Only show when collapsed */}
              {tableOfContents.length > 0 && !isTocExpanded && (
                <View className="toc-expand-trigger">
                  <View onClick={() => setIsTocExpanded(true)}>
                    <XIcon name='carbon:list' size={18}
                           style={{display: 'inline-block', width: '18px', height: '18px', fontSize: '18px'}}/>
                    <Text>{t('page.post_detail.table_of_contents')}</Text>
                    <XIcon name='carbon:chevron-right' size={18}
                           style={{display: 'inline-block', width: '18px', height: '18px', fontSize: '18px'}}/>
                  </View>
                </View>
              )}

              {/* Post Header */}
              <View className="post-header">
                <Text className="post-title">{displayTitle}</Text>
                <View className="post-meta">
                  <View className="meta-item">
                    <View className="meta-icon">
                      <XIcon name='carbon:user' size={14}
                             style={{display: 'inline-block', width: '14px', height: '14px', fontSize: '14px'}}/>
                    </View>
                    <Text>{post.authorName}</Text>
                  </View>
                  <View className="meta-item">
                    <View className="meta-icon">
                      <XIcon name='carbon:calendar' size={14}
                             style={{display: 'inline-block', width: '14px', height: '14px', fontSize: '14px'}}/>
                    </View>
                    <Text>{formatDate(post.createdAt)}</Text>
                  </View>
                  <View className="meta-item">
                    <View className="meta-icon">
                      <XIcon name='carbon:view' size={14}
                             style={{display: 'inline-block', width: '14px', height: '14px', fontSize: '14px'}}/>
                    </View>
                    <Text>{post.visits || 0}</Text>
                  </View>
                  <View className="meta-item">
                    <View className="meta-icon">
                      <XIcon name='carbon:thumb-up' size={14}
                             style={{display: 'inline-block', width: '14px', height: '14px', fontSize: '14px'}}/>
                    </View>
                    <Text>{post.likes || 0}</Text>
                  </View>
                </View>
              </View>

              {/* Post Content */}
              <View className="post-content" ref={contentRef}>
                <ContentViewer content={displayContent} type="markdown"/>
              </View>

              {/* Post Actions */}
              <View className="post-actions">
                <View className="action-buttons">
                  <View
                    onClick={handleLike}
                    className={`action-btn${isLiked ? ' liked' : ''}`}
                    aria-label={t('page.post_detail.likes')}
                  >
                    <XIcon name={isLiked ? 'carbon:thumbs-up-filled' : 'carbon:thumbs-up'} size={20}
                           style={{display: 'inline-block', width: '20px', height: '20px', fontSize: '20px'}}/>
                  </View>
                  <View
                    onClick={handleBookmark}
                    className={`action-btn${isBookmarked ? ' bookmarked' : ''}`}
                    aria-label={t('page.post_detail.bookmark')}
                  >
                    <XIcon name={isBookmarked ? 'carbon:bookmark-filled' : 'carbon:bookmark'} size={20}
                           style={{display: 'inline-block', width: '20px', height: '20px', fontSize: '20px'}}/>
                  </View>
                  <View
                    onClick={handleShare}
                    className="action-btn"
                    aria-label={t('page.post_detail.share')}
                  >
                    <XIcon name='carbon:share' size={20}
                           style={{display: 'inline-block', width: '20px', height: '20px', fontSize: '20px'}}/>
                  </View>
                </View>
              </View>
            </View>
          </View>
        </View>
      </View>

      {/* Comments Section */}
      <CommentSection
        objectId={postId}
        contentType="CONTENT_TYPE_POST"
        onUpdateComments={() => {
        }}
      />

      {/* Related Posts */}
      <View className="related-section">
        <View className="section-header">
          <Text className="section-title">
            <XIcon name='carbon:document' size={20}
                   style={{display: 'inline-block', width: '20px', height: '20px', fontSize: '20px'}}/>
            <Text>{t('page.post_detail.related_posts')}</Text>
          </Text>
        </View>
        {relatedPostsQuery && (
          <PostList
            queryParams={relatedPostsQuery}
            fieldMask="id,status,sort_order,is_featured,visits,likes,comment_count,author_name,available_languages,created_at,translations.id,translations.post_id,translations.language_code,translations.title,translations.summary,translations.thumbnail"
            orderBy={['-sortOrder']}
            page={1}
            pageSize={3}
            showSkeleton={false}
          />
        )}
      </View>

      {/* Back to Top Button */}
      <BackToTop scrollThreshold={500}/>
    </View>
  );
}
