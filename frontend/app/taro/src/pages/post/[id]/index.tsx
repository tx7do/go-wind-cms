import {useEffect, useState, useRef, useMemo, useCallback} from 'react';
import {useTranslation} from 'react-i18next';
import {View, Text, Image} from '@tarojs/components';
import Taro from '@tarojs/taro';

import CommentSection from '@/components/comment/CommentSection';
import ContentViewer from '@/components/content/ContentViewer';
import PostList from '@/components/post/PostList';
import BackToTop from '@/components/layout/BackToTop';

import {usePostStore} from '@/store/slices/post/hooks';
import {formatDate} from "@/utils";
import {contentservicev1_Post} from "@/api/generated/app/service/v1";

import '../../../globals.css'; // 导入全局 CSS，确保 CSS 变量可用
import styles from './post-detail.scss';

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
  const {t} = useTranslation('page');
  const postStore = usePostStore();

  // 直接使用 store 中的数据，而不是本地 state
  const post = postStore.detail as contentservicev1_Post | null;

  const [localLoading, setLocalLoading] = useState(true);
  const isLoading = localLoading || (postStore.loading && !post);
  const [tableOfContents, setTableOfContents] = useState<TocItem[]>([]);
  const [activeHeading, setActiveHeading] = useState('');
  const [isLiked, setIsLiked] = useState(false);
  const [isBookmarked, setIsBookmarked] = useState(false);
  const [isTocExpanded, setIsTocExpanded] = useState(true);

  const contentRef = useRef<HTMLDivElement>(null);

  const postId = useMemo(() => {
    const pages = Taro.getCurrentPages();
    const currentPage = pages[pages.length - 1];
    const id = (currentPage as any).options?.id;
    return id ? parseInt(id) : null;
  }, []);

  // 计算属性 - 使用 postStore 提供的工具函数
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
          // TODO: Taro 中设置页面标题
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
  }, [postId]); // 只依赖 postId

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

      const headings = contentEl.querySelectorAll('h2, h3');
      const toc: TocItem[] = [];

      console.log('[TOC] Found headings:', headings.length);

      headings.forEach((heading, index) => {
        const level = heading.tagName === 'H2' ? 2 : 3;
        const id = `heading-${index}`;
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

      // 页面加载时检查 URL hash，自动滚动到对应位置
      // TODO: Taro 中使用 pageScrollTo API
    }, 100); // 等待内容完全渲染

    return () => clearTimeout(timeoutId);
  }, [displayContent]); // 依赖计算后的内容

  // 监听内容变化，重新生成目录
  useEffect(() => {
    if (displayContent && tableOfContents.length === 0) {
      // 如果还没有生成目录，立即生成一次
      generateTableOfContents();
    }
  }, [displayContent, tableOfContents.length]);

  // 生成目录函数
  const generateTableOfContents = useCallback(() => {
    setTimeout(() => {
      if (!contentRef.current) return;

      const contentEl = contentRef.current;
      const headings = contentEl.querySelectorAll('h2, h3');
      const toc: TocItem[] = [];

      console.log('[GenerateTOC] Found headings:', headings.length);

      headings.forEach((heading, index) => {
        const level = heading.tagName === 'H2' ? 2 : 3;
        const id = `heading-${index}`;

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
    }, 500);
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
    // TODO: 从 Taro.getCurrentPages() 获取参数
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
    // TODO: Taro 中使用 onShareAppMessage 配置分享
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

      // Taro 中使用 pageScrollTo
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
      <View className={styles['post-detail-page']}>
        {/* Loading skeleton */}
        <View className={styles['back-navigation']}>
          <View className={styles['skeleton-btn']}></View>
        </View>
        <View className={styles['post-article']}>
          <View className={styles['post-banner']}>
            <View className={styles['skeleton-image']}></View>
          </View>
          <View className={styles['post-wrapper']}>
            <View className={styles['toc-sidebar']}>
              <View className={styles['toc-container']}>
                <View className={styles['skeleton-line']} style={{width: '200px', height: '24px'}}></View>
                <View className={styles['skeleton-line']}
                      style={{width: '180px', height: '20px', marginTop: '16px'}}></View>
                <View className={styles['skeleton-line']}
                      style={{width: '160px', height: '20px', marginTop: '8px'}}></View>
                <View className={styles['skeleton-line']}
                      style={{width: '140px', height: '20px', marginTop: '8px'}}></View>
              </View>
            </View>
            <View className={styles['article-content']}>
              <View className={styles['post-header']}>
                <View className={styles['skeleton-title']} style={{width: '80%', height: '48px'}}></View>
                <View className={styles['skeleton-subtitle']}
                      style={{width: '60%', height: '32px', marginTop: '16px'}}></View>
                <View className={styles['post-meta']}>
                  <View className={styles['skeleton-meta']}
                        style={{width: '100px', height: '20px'}}></View>
                  <View className={styles['skeleton-meta']}
                        style={{width: '100px', height: '20px'}}></View>
                  <View className={styles['skeleton-meta']}
                        style={{width: '100px', height: '20px'}}></View>
                  <View className={styles['skeleton-meta']}
                        style={{width: '100px', height: '20px'}}></View>
                </View>
              </View>
              <View className={styles['post-content']}>
                <View className={styles['skeleton-paragraph']}></View>
              </View>
            </View>
          </View>
        </View>
      </View>
    );
  }

  if (!post) {
    return (
      <View className={styles['post-detail-page']}>
        <View className={styles['empty-state']}>Post not found</View>
      </View>
    );
  }

  return (
    <View className={styles['post-detail-page']}>
      {/* Back Navigation */}
      <View className={styles['back-navigation']}>
        <View onClick={handleBack} className={styles['back-btn']} aria-label={t('post_detail.back')}>
          <Text>← </Text>
          <Text>{t('post_detail.back')}</Text>
        </View>
      </View>

      {/* Post Article */}
      <View className={styles['post-article']}>
        {/* Post Thumbnail Banner */}
        {displayThumbnail && (
          <View className={styles['post-banner']}>
            <Image src={displayThumbnail}/>
            <View className={styles['banner-overlay']}/>
          </View>
        )}

        <View className={`${styles['post-wrapper']} ${!isTocExpanded ? styles['toc-collapsed'] : ''}`}>
          {/* Table of Contents Sidebar */}
          {tableOfContents.length > 0 && isTocExpanded && (
            <View className={styles['toc-sidebar']}>
              <View className={styles['toc-container']}>
                <View className={styles['toc-header']}>
                  <Text className={styles['toc-title']}>
                    <Text>📑</Text>
                    <Text>{t('post_detail.table_of_contents')}</Text>
                  </Text>
                  <View
                    onClick={() => setIsTocExpanded(false)}
                    className={styles['toc-collapse-btn']}
                  >
                    <Text>◀</Text>
                  </View>
                </View>
                <View className={styles['toc-list']}>
                  {tableOfContents.map(item => (
                    <View
                      key={item.id}
                      className={`${styles['toc-item']} ${styles[`level-${item.level}`]} ${activeHeading === item.id ? styles['active'] : ''}`}
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

          {/* TOC Expand Button */}
          {tableOfContents.length > 0 && !isTocExpanded && (
            <View className={styles['toc-expand-trigger']}>
              <View onClick={() => setIsTocExpanded(true)}>
                <Text>📑 </Text>
                <Text>{t('post_detail.table_of_contents')}</Text>
                <Text> ▶</Text>
              </View>
            </View>
          )}

          <View className={styles['article-content']}>
            {/* Post Header */}
            <View className={styles['post-header']}>
              <Text className={styles['post-title']}>{displayTitle}</Text>
              <View className={styles['post-meta']}>
                <View className={styles['meta-item']}>
                  <Text>👤 </Text>
                  <Text>{post.authorName}</Text>
                </View>
                <View className={styles['meta-item']}>
                  <Text>📅 </Text>
                  <Text>{formatDate(post.createdAt)}</Text>
                </View>
                <View className={styles['meta-item']}>
                  <Text>👁️ </Text>
                  <Text>{post.visits || 0}</Text>
                </View>
                <View className={styles['meta-item']}>
                  <Text>👍 </Text>
                  <Text>{post.likes || 0}</Text>
                </View>
              </View>
            </View>

            {/* Post Content */}
            <View className={styles['post-content']} ref={contentRef}>
              <ContentViewer content={displayContent} type="markdown"/>
            </View>

            {/* Post Actions */}
            <View className={styles['post-actions']}>
              <View className={styles['action-buttons']}>
                <View
                  onClick={handleLike}
                  className={`${styles['action-btn']} ${isLiked ? styles['liked'] : ''}`}
                  aria-label={t('post_detail.likes')}
                >
                  <Text>{isLiked ? '👍' : '👎'}</Text>
                </View>
                <View
                  onClick={handleBookmark}
                  className={`${styles['action-btn']} ${isBookmarked ? styles['bookmarked'] : ''}`}
                  aria-label={t('post_detail.bookmark')}
                >
                  <Text>{isBookmarked ? '🔖' : '📑'}</Text>
                </View>
                <View
                  onClick={handleShare}
                  className={styles['action-btn']}
                  aria-label={t('post_detail.share')}
                >
                  <Text>📤</Text>
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
      <View className={styles['related-section']}>
        <View className={styles['section-header']}>
          <Text className={styles['section-title']}>
            <Text>📚 </Text>
            <Text>{t('post_detail.related_posts')}</Text>
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
