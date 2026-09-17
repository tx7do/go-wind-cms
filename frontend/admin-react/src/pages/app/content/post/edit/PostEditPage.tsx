import { useCallback, useEffect, useMemo, useRef, useState } from 'react';
import { App, Button, Input, Select, Space, Tag } from 'antd';
import {
  ArrowLeftOutlined,
  CheckOutlined,
  CloseOutlined,
  MenuFoldOutlined,
  SettingOutlined,
  TranslationOutlined,
} from '@ant-design/icons';
import { useNavigate, useParams, useSearchParams, useBlocker } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

import { apiClient } from '@/api/client';
import { PaginationQuery } from '@/core';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import { Editor, EditorType } from '@/components/common/Editor';
import { fetchListLanguages } from '@/api/hooks/language';
import { fetchFlattenedCategoryOptions } from '@/api/hooks/category';
import { uploadMediaAsset } from '@/api/hooks/media-asset';
import { useTranslate } from '@/api/hooks/tag';
import { POST_SAVE_MASK_KEYS, useCreatePost, useUpdatePost } from '@/api/hooks/post';
import { compressImageFile } from '@/utils/image';
import {
  getPostStatusColor,
  getPostStatusLabel,
} from '../constants';
import PostSettingsPanel from './PostSettingsPanel';
import type { PostEditProps } from './types';

// 文章只提供富文本与 Markdown 两种编辑器；代码/JSON/纯文本编辑器不适合文章场景
const POST_EDITOR_TYPE_OPTIONS = [
  { value: EditorType.RICH_TEXT, key: 'EDITOR_TYPE_RICH_TEXT' },
  { value: EditorType.MARKDOWN, key: 'EDITOR_TYPE_MARKDOWN' },
];

const emptyForm = (lang: string): PostEditProps => ({
  title: '',
  content: '',
  lang,
  editorType: EditorType.RICH_TEXT,
  seo: {},
});

function getDraftKey(postId: null | number, lang: string, isCreateMode: boolean) {
  return isCreateMode ? `post-draft-create-${lang}` : `post-draft-edit-${postId}-${lang}`;
}

/** 表单中需要参与未保存检测的字段快照（对齐 vben snapshotPayload） */
function snapshotPayload(f: PostEditProps) {
  return {
    title: f.title,
    content: f.content,
    editorType: f.editorType,
    code: f.code ?? '',
    categoryIds: [...(f.categoryIds ?? [])].sort(),
    isFeatured: f.isFeatured ?? false,
    sortOrder: f.sortOrder ?? 0,
    disallowComment: f.disallowComment ?? false,
    publishTime: f.publishTime,
    summary: f.summary ?? '',
    slug: f.slug ?? '',
    thumbnail: f.thumbnail ?? '',
    seo: f.seo ?? {},
  };
}

/**
 * 文章编辑页（独立编辑页范式）：
 * 多语言翻译、本地草稿（自动暂存 + 恢复询问）、右栏设置、定时发布、未保存拦截
 */
const PostEditPage = () => {
  const { t } = useTranslation('post');
  const navigate = useNavigate();
  const { id } = useParams();
  const [searchParams, setSearchParams] = useSearchParams();
  const { message, notification, modal } = App.useApp();

  const initLang = searchParams.get('lang') || 'zh-CN';
  const isCreateMode = !id;
  const postId = isCreateMode ? null : Number(id);

  const [lang, setLang] = useState(initLang);
  const [formData, setFormData] = useState<PostEditProps>(() => emptyForm(initLang));
  const [status, setStatus] = useState<string | undefined>(undefined);
  const [languageOptions, setLanguageOptions] = useState<
    { label: string; value: string; hasTranslation?: boolean }[]
  >([]);
  const [saving, setSaving] = useState(false);
  const [settingsOpen, setSettingsOpen] = useState(true);
  const [categoryOptions, setCategoryOptions] = useState<{ label: string; value: number }[]>([]);
  const [categoryOptionsLoading, setCategoryOptionsLoading] = useState(false);

  const translateMutation = useTranslate();
  const createMutation = useCreatePost();
  const updateMutation = useUpdatePost();

  // 未保存检测基线；变更后防抖自动暂存本地草稿
  const savedSnapshotRef = useRef('');
  const draftSaveTimerRef = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);
  const draftPromptShownRef = useRef(false);
  /** 发布成功后的返回跳过拦截 */
  const leavingBySaveRef = useRef(false);
  /** 新建保存成功后拿到的新文章 id（state 异步，路由切换需同步值） */
  const createdIdRef = useRef<number | null>(null);

  const goBack = useCallback(() => navigate('/content/posts'), [navigate]);

  const patch = useCallback((p: Partial<PostEditProps>) => {
    setFormData((prev) => ({ ...prev, ...p }));
  }, []);

  const hasUnsavedChanges = useCallback(
    (f: PostEditProps) =>
      !!savedSnapshotRef.current &&
      JSON.stringify(snapshotPayload(f)) !== savedSnapshotRef.current,
    [],
  );

  const loadDraft = useCallback(
    (currentLang: string): PostEditProps | null => {
      const raw = localStorage.getItem(getDraftKey(postId, currentLang, isCreateMode));
      if (!raw) return null;
      try {
        const draft = JSON.parse(raw) as PostEditProps;
        // 早期草稿可能没有 seo 字段，补齐为对象避免绑定空引用
        if (!draft.seo) draft.seo = {};
        return draft;
      } catch {
        return null;
      }
    },
    [postId, isCreateMode],
  );

  const saveDraft = useCallback(
    (f: PostEditProps) => {
      localStorage.setItem(
        getDraftKey(postId, f.lang, isCreateMode),
        JSON.stringify(f),
      );
    },
    [postId, isCreateMode],
  );

  const clearDraft = useCallback(
    (currentLang: string) => {
      localStorage.removeItem(getDraftKey(postId, currentLang, isCreateMode));
    },
    [postId, isCreateMode],
  );

  const refreshSavedSnapshot = useCallback((f: PostEditProps) => {
    savedSnapshotRef.current = JSON.stringify(snapshotPayload(f));
  }, []);

  // 表单变化：有未保存修改时防抖自动暂存到本地草稿
  useEffect(() => {
    if (!hasUnsavedChanges(formData)) return;
    if (draftSaveTimerRef.current) {
      clearTimeout(draftSaveTimerRef.current);
    }
    draftSaveTimerRef.current = setTimeout(() => saveDraft(formData), 1500);
    return () => {
      if (draftSaveTimerRef.current) clearTimeout(draftSaveTimerRef.current);
    };
  }, [formData, hasUnsavedChanges, saveDraft]);

  // 关闭/刷新浏览器标签前拦截未保存修改
  useEffect(() => {
    const handler = (event: BeforeUnloadEvent) => {
      if (hasUnsavedChanges(formData)) {
        event.preventDefault();
        event.returnValue = '';
      }
    };
    window.addEventListener('beforeunload', handler);
    return () => window.removeEventListener('beforeunload', handler);
  }, [formData, hasUnsavedChanges]);

  // 站内路由离开拦截（发布成功后的返回跳过）
  const blocker = useBlocker(({ currentLocation, nextLocation }: { currentLocation: { pathname: string }; nextLocation: { pathname: string } }) => {
    if (leavingBySaveRef.current) return false;
    if (currentLocation.pathname === nextLocation.pathname) return false;
    return hasUnsavedChangesRef.current;
  });
  const hasUnsavedChangesRef = useRef(false);
  useEffect(() => {
    hasUnsavedChangesRef.current = hasUnsavedChanges(formData);
  }, [formData, hasUnsavedChanges]);

  useEffect(() => {
    if (blocker.state === 'blocked') {
      const modalRef = modal.confirm({
        title: t('unsavedTitle'),
        content: t('unsavedMessage'),
        okText: t('common:button.ok'),
        cancelText: t('common:button.cancel'),
        onOk: () => {
          blocker.proceed?.();
        },
        onCancel: () => {
          blocker.reset?.();
        },
      });
      return () => modalRef.destroy();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [blocker.state]);

  /** 加载文章（编辑模式）：翻译缺失时回退首条并提示需要翻译 */
  const fetchPost = useCallback(
    async (currentLang: string) => {
      if (isCreateMode || !postId) return;
      try {
        const item = await apiClient.postService.Get({ id: postId });
        if (!item.translations || item.translations.length === 0) {
          throw new Error('No translations found for post');
        }

        let langItem = item.translations.find((tr) => tr.languageCode === currentLang);
        let missing = false;
        if (!langItem) {
          langItem = item.translations[0];
          missing = true;
        }
        if (!langItem) {
          throw new Error('No translations found for post');
        }

        const availableLanguages = item.availableLanguages || [];
        setLanguageOptions((prev) =>
          prev.map((option) => ({
            ...option,
            hasTranslation: availableLanguages.includes(option.value),
          })),
        );

        const base: PostEditProps = {
          id: item.id,
          title: langItem.title || '',
          content: langItem.content || '',
          lang: currentLang,
          editorType: item.editorType || EditorType.RICH_TEXT,
          status: item.status ?? undefined,
          code: item.code || '',
          categoryIds: item.categoryIds ? [...item.categoryIds] : [],
          isFeatured: item.isFeatured ?? false,
          sortOrder: item.sortOrder ?? 0,
          disallowComment: item.disallowComment ?? false,
          publishTime: item.publishTime ?? undefined,
          // 缩略图帖子级，全语言共用
          thumbnail: item.thumbnail || '',
          summary: langItem.summary || '',
          slug: langItem.slug || '',
          seo: { ...(langItem.seo ?? {}) },
        };

        setFormData(base);
        setStatus(item.status ?? undefined);
        refreshSavedSnapshot(base);

        if (missing) {
          notification.info({ message: t('translationNotExists'), placement: 'bottomRight' });
        }
      } catch (error) {
        console.error('Failed to load post:', error);
        notification.error({ message: t('loadFailed'), placement: 'bottomRight' });
        throw error;
      }
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [isCreateMode, postId],
  );

  /** 存在未发布的本地草稿时询问用户是否恢复（不再静默覆盖服务端数据） */
  const maybePromptPendingDraft = useCallback(() => {
    if (draftPromptShownRef.current) return;
    const pending = !!loadDraft(initLang);
    if (!pending) return;
    draftPromptShownRef.current = true;
    modal.confirm({
      title: t('draftFound'),
      content: t('draftFoundMessage'),
      okText: t('draftRestore'),
      cancelText: t('common:button.cancel'),
      onOk: () => {
        const draft = loadDraft(initLang);
        if (draft) {
          setFormData(draft);
        }
      },
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [initLang, loadDraft]);

  useEffect(() => {
    (async () => {
      setCategoryOptionsLoading(true);
      try {
        const options = await fetchFlattenedCategoryOptions(initLang);
        setCategoryOptions(options);
      } catch {
        setCategoryOptions([]);
      } finally {
        setCategoryOptionsLoading(false);
      }
      try {
        const resp = await fetchListLanguages(
          new PaginationQuery({ orderBy: ['sortOrder'] }),
        );
        setLanguageOptions(
          resp.items
            ?.filter((l: any) => l.isEnabled !== false)
            .map((l: any) => ({
              label: l.nativeName || '',
              value: l.languageCode || '',
            })) || [],
        );
      } catch (error) {
        console.error('Failed to load language list:', error);
        notification.error({ message: t('loadLanguageFailed'), placement: 'bottomRight' });
      }

      if (isCreateMode) {
        const fresh = emptyForm(initLang);
        setFormData(fresh);
        refreshSavedSnapshot(fresh);
      } else {
        await fetchPost(initLang);
      }
      maybePromptPendingDraft();
    })();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [initLang, isCreateMode, postId]);

  const handleLanguageChange = async (newLang: string) => {
    setLang(newLang);
    setSearchParams((prev) => {
      const next = new URLSearchParams(prev);
      next.set('lang', newLang);
      return next;
    });
    // 分类名称依赖当前语言，按新语言重建选项
    setCategoryOptionsLoading(true);
    fetchFlattenedCategoryOptions(newLang)
      .then(setCategoryOptions)
      .catch(() => setCategoryOptions([]))
      .finally(() => setCategoryOptionsLoading(false));
    if (!isCreateMode) {
      await fetchPost(newLang);
    } else {
      patch({ lang: newLang });
    }
  };

  /** 一键翻译（当前已有内容时先确认覆盖） */
  const handleTranslate = async () => {
    if (formData.title || formData.content) {
      const confirmed = await new Promise<boolean>((resolve) => {
        modal.confirm({
          title: t('translationExists'),
          okText: t('common:button.ok'),
          cancelText: t('common:button.cancel'),
          onOk: () => resolve(true),
          onCancel: () => resolve(false),
        });
      });
      if (!confirmed) return;
    }

    try {
      const title = await translateMutation.mutateAsync({
        sourceLanguage: 'auto',
        targetLanguage: lang,
        content: formData.title,
      });
      patch({ title });
    } catch {
      message.error(t('translateTitleFailed'));
      return;
    }

    try {
      const content = await translateMutation.mutateAsync({
        sourceLanguage: 'auto',
        targetLanguage: lang,
        content: formData.content,
      });
      patch({ content });
    } catch {
      message.error(t('translateContentFailed'));
    }

    // 摘要非空时一并翻译（失败不影响整体流程）
    if (formData.summary) {
      try {
        const summary = await translateMutation.mutateAsync({
          sourceLanguage: 'auto',
          targetLanguage: lang,
          content: formData.summary,
        });
        patch({ summary });
      } catch {
        console.error('Summary translation failed');
      }
    }
  };

  /** 编辑器正文插图：先压缩再上传（onProgress 上报全局进度 toast） */
  const handleUploadImage = async (
    file: File,
    onProgress?: (percent: number) => void,
  ): Promise<string> => {
    try {
      const compressed = await compressImageFile(file);
      const resp = await uploadMediaAsset(
        {},
        compressed,
        (progressEvent: { loaded?: number; total?: number }) => {
          if (!onProgress) return;
          const total = progressEvent?.total ?? 0;
          onProgress(total > 0 ? Math.floor(((progressEvent.loaded ?? 0) / total) * 100) : -1);
        },
      );
      return (resp as { objectName?: string })?.objectName || '';
    } catch (error) {
      console.error('Image upload failed:', error);
      message.error((error as Error)?.message || t('publishFailed'));
      return '';
    }
  };

  /**
   * 保存文章（status 决定保存为草稿还是发布）。
   * 新建模式下保存成功后会自动切换为编辑模式（后续保存走 Update）。
   * 返回空字符串表示成功，否则返回错误提示文案。
   */
  const savePost = async (nextStatus: string): Promise<string> => {
    if (!formData.title) return t('requiredTitle');
    if (!formData.content) return t('requiredContent');
    if (!formData.categoryIds || formData.categoryIds.length === 0) {
      return t('requiredCategory');
    }

    const wasCreateMode = isCreateMode;
    // 全部为空的 seo 不提交，避免用空对象覆盖服务端已有值
    const seo = formData.seo ?? {};
    const hasSeo = Object.values(seo).some(
      (v) => typeof v === 'string' && v.length > 0,
    );
    const data: Record<string, any> = {
      editorType: formData.editorType,
      status: nextStatus,
      code: formData.code || '',
      isFeatured: formData.isFeatured ?? false,
      sortOrder: formData.sortOrder ?? 0,
      disallowComment: formData.disallowComment ?? false,
      categoryIds: formData.categoryIds ?? [],
      // 缩略图为帖子级字段，全语言共用一张
      thumbnail: formData.thumbnail || '',
      translations: [
        {
          title: formData.title,
          content: formData.content,
          languageCode: lang,
          summary: formData.summary || '',
          slug: formData.slug || '',
          seo: hasSeo ? seo : undefined,
        },
      ],
    };
    // 始终携带 publishTime 进掩码：未设置时提交 null，由后端置空该列（支持清空）
    (data as { publishTime?: string | null }).publishTime =
      formData.publishTime ?? null;

    try {
      if (wasCreateMode) {
        // useCreatePost 接收裸 values，内部再包 { data }
        const created = await createMutation.mutateAsync(data);
        // 先按创建模式的 key 清理本地草稿，再切换为编辑模式
        clearDraft(lang);
        if (created?.id) {
          createdIdRef.current = Number(created.id);
          setFormData((prev) => ({ ...prev, id: Number(created.id) }));
        }
      } else {
        await updateMutation.mutateAsync({
          id: formData.id || 0,
          data,
          maskKeys: POST_SAVE_MASK_KEYS,
        });
        clearDraft(lang);
      }

      setStatus(nextStatus);
      refreshSavedSnapshot({ ...formData, status: nextStatus });
      return '';
    } catch (error) {
      console.error('Failed to save post:', error);
      return wasCreateMode ? t('publishFailed') : t('saveDraftFailed');
    }
  };

  /** 保存为服务端草稿 */
  const handleSaveDraft = async () => {
    setSaving(true);
    try {
      const resp = await savePost('POST_STATUS_DRAFT');
      if (resp) {
        message.error(resp);
        return;
      }
      message.success(t('saveDraftSuccess'));

      // 新建模式下保存草稿后已生成文章，切换到编辑路由，后续保存走更新
      if (isCreateMode && createdIdRef.current) {
        leavingBySaveRef.current = true;
        await navigate(
          {
            pathname: `/content/posts/edit/${createdIdRef.current}`,
            search: searchParams.toString(),
          },
          { replace: true },
        );
        leavingBySaveRef.current = false;
      }
    } finally {
      setSaving(false);
    }
  };

  /**
   * 发布文章。设置了未来发布时间的未发布文章会转为定时发布（SCHEDULED），
   * 到达时间后由后端调度器自动置为已发布。
   */
  const handlePublish = async () => {
    setSaving(true);
    try {
      const publishTime = formData.publishTime;
      let target = 'POST_STATUS_PUBLISHED';
      if (publishTime) {
        const ts = new Date(publishTime).getTime();
        const isFuture = !Number.isNaN(ts) && ts > Date.now() + 60_000;
        const isScheduled = status === 'POST_STATUS_SCHEDULED';
        const isPublished = status === 'POST_STATUS_PUBLISHED';
        if (isFuture && !isScheduled && !isPublished) {
          target = 'POST_STATUS_SCHEDULED';
        }
      }

      const resp = await savePost(target);
      if (resp) {
        message.error(resp);
        return;
      }

      message.success(
        target === 'POST_STATUS_SCHEDULED' ? t('scheduleSuccess') : t('publishSuccess'),
      );
      leavingBySaveRef.current = true;
      goBack();
    } finally {
      setSaving(false);
    }
  };

  const editorTypeOptions = useMemo(
    () =>
      POST_EDITOR_TYPE_OPTIONS.map((o) => ({
        value: o.value,
        label: t(`editorTypeMap.${o.key}`),
      })),
    [t],
  );

  return (
    <ContentContainer heightMode="fixed" padding="0" bottomMargin={0}>
      <div style={{ height: '100%', display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
        {/* 头部：返回 + 标题 + 状态 + 语言 + 翻译 + 编辑器类型 */}
        <div style={{ display: 'flex', alignItems: 'center', gap: 8, padding: '12px 16px 8px', flexWrap: 'wrap' }}>
          <Button type="text" icon={<ArrowLeftOutlined />} onClick={goBack} />
          <Input
            placeholder={t('titlePlaceholder')}
            size="large"
            style={{ flex: 1, minWidth: 200 }}
            value={formData.title}
            onChange={(e) => patch({ title: e.target.value })}
          />
          {!isCreateMode && status && (
            <Tag color={getPostStatusColor(status)}>{getPostStatusLabel(t, status)}</Tag>
          )}
          <Select
            value={lang}
            style={{ width: 200 }}
            onChange={handleLanguageChange}
            options={languageOptions.map((option) => ({
              label: (
                <span>
                  {option.label}
                  {option.hasTranslation ? (
                    <CheckOutlined style={{ marginLeft: 8, color: 'var(--ant-color-success)' }} title={t('hasTranslation')} />
                  ) : (
                    <CloseOutlined style={{ marginLeft: 8, color: 'var(--ant-color-warning)' }} title={t('noTranslation')} />
                  )}
                </span>
              ),
              value: option.value,
            }))}
          />
          <Button
            type="primary"
            icon={<TranslationOutlined />}
            loading={translateMutation.isPending}
            onClick={handleTranslate}
          >
            {t('oneClickTranslate')}
          </Button>
          <Select
            style={{ width: 200 }}
            value={formData.editorType}
            onChange={(v: string) => patch({ editorType: v })}
            options={editorTypeOptions}
          />
        </div>

        {/* 主体：左侧编辑器 + 右侧内嵌设置栏 */}
        <div style={{ flex: 1, minHeight: 0, display: 'flex', overflow: 'hidden' }}>
          <div style={{ flex: 1, minWidth: 0, minHeight: 0 }}>
            <Editor
              height="100%"
              value={formData.content}
              editorType={formData.editorType}
              placeholder={t('contentPlaceholder')}
              uploadImage={handleUploadImage}
              onChange={(v) => patch({ content: v })}
            />
          </div>

          {/* 右侧：内嵌设置栏 + 常驻竖向开关轨 */}
          <div style={{ display: 'flex', height: '100%', flexShrink: 0 }}>
            <aside
              style={{
                height: '100%',
                overflow: 'hidden',
                transition: 'width 0.3s',
                width: settingsOpen ? 340 : 0,
                borderLeft: settingsOpen ? '1px solid var(--ant-color-split)' : 'none',
              }}
            >
              <div style={{ height: '100%', width: 340, display: 'flex', flexDirection: 'column' }}>
                <div
                  style={{
                    flexShrink: 0,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                    borderBottom: '1px solid var(--ant-color-split)',
                    padding: '8px 16px',
                  }}
                >
                  <span style={{ fontSize: 14, fontWeight: 500 }}>{t('settings')}</span>
                  <Button size="small" type="text" icon={<MenuFoldOutlined />} onClick={() => setSettingsOpen(false)} />
                </div>
                <div style={{ minHeight: 0, flex: 1, overflowY: 'auto', padding: '12px 16px' }}>
                  <PostSettingsPanel
                    formData={formData}
                    patch={patch}
                    categoryOptions={categoryOptions}
                    categoryOptionsLoading={categoryOptionsLoading}
                  />
                </div>
              </div>
            </aside>

            <button
              style={{
                display: 'flex',
                width: 32,
                flexShrink: 0,
                flexDirection: 'column',
                alignItems: 'center',
                gap: 8,
                border: 'none',
                borderLeft: settingsOpen ? 'none' : '1px solid var(--ant-color-split)',
                padding: '12px 0',
                cursor: 'pointer',
                background: 'transparent',
                color: 'var(--ant-color-text-tertiary)',
              }}
              title={t('settings')}
              onClick={() => setSettingsOpen((v) => !v)}
            >
              <SettingOutlined style={{ fontSize: 16, color: settingsOpen ? 'var(--ant-color-primary)' : undefined }} />
              <span style={{ fontSize: 12, writingMode: 'vertical-rl', letterSpacing: 2 }}>
                {t('settings')}
              </span>
            </button>
          </div>
        </div>

        {/* 底部操作 */}
        <div style={{ display: 'flex', justifyContent: 'flex-end', padding: '8px 16px', borderTop: '1px solid var(--ant-color-split)' }}>
          <Space>
            <Button onClick={handleSaveDraft} loading={saving}>
              {t('saveDraft')}
            </Button>
            <Button type="primary" danger onClick={handlePublish} loading={saving}>
              {t('publish')}
            </Button>
          </Space>
        </div>
      </div>
    </ContentContainer>
  );
};

export default PostEditPage;
