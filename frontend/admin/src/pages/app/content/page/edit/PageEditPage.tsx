import { useCallback, useEffect, useMemo, useState } from 'react';
import {
  App,
  Button,
  Col,
  Input,
  Row,
  Select,
  Space,
} from 'antd';
import {
  ArrowLeftOutlined,
  CheckOutlined,
  CloseOutlined,
  DownOutlined,
  TranslationOutlined,
  UpOutlined,
} from '@ant-design/icons';
import { useNavigate, useParams, useSearchParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

import { apiClient } from '@/api/client';
import { PaginationQuery } from '@/core';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import { Editor, EditorType } from '@/components/common/Editor';
import {
  fetchListContentModels,
  fetchListFieldDefinitions,
} from '@/api/hooks/content-model';
import { fetchListLanguages } from '@/api/hooks/language';
import { uploadFile } from '@/api/hooks/file-transfer';
import { useTranslate } from '@/api/hooks/tag';
import {
  editorTypeOptions,
  pageStatusOptions,
  pageTypeOptions,
  sectionTypeOptions,
  getSectionTypeLabel,
} from '../constants';

interface SectionFormItemTranslation {
  id?: number;
  languageCode?: string;
  content?: Record<string, string>;
}

interface SectionFormItem {
  id?: number;
  type?: string;
  name?: string;
  sortOrder?: number;
  config?: Record<string, string>;
  translations?: SectionFormItemTranslation[];
}

interface PageEditProps {
  id?: number;
  title: string;
  slug: string;
  content: string;
  lang: string;
  editorType: string;
  parentId?: number;
  type?: string;
  status?: string;
  showInNavigation?: boolean;
  disallowComment?: boolean;
  template?: string;
  isCustomTemplate?: boolean;
  sortOrder?: number;
  contentModelId?: number;
  customFields?: Record<string, string>;
  /** 页面嵌套区块列表。作为页面子部件随页面整体读写（对齐后端 Page.Sections 整体替换语义） */
  sections?: SectionFormItem[];
}

interface LanguageOption {
  label: string;
  value: string;
  hasTranslation?: boolean;
}

const emptyForm = (lang: string): PageEditProps => ({
  title: '',
  slug: '',
  content: '',
  lang,
  editorType: EditorType.MARKDOWN,
  type: 'PAGE_TYPE_DEFAULT',
  status: 'PAGE_STATUS_DRAFT',
  sections: [],
});

function getDraftKey(pageId: null | number, lang: string, isCreateMode: boolean) {
  return isCreateMode
    ? `page-draft-create-${lang}`
    : `page-draft-edit-${pageId}-${lang}`;
}

/**
 * 页面编辑页（独立编辑页范式）：
 * 多语言翻译（✓/○ 标记 + 一键翻译）、本地草稿、编辑器切换、内容模型动态字段、
 * 嵌套区块管理、发布（主编辑器内容合并进首个正文区块）
 */
const PageEditPage = () => {
  const { t } = useTranslation('page');
  const navigate = useNavigate();
  const { id } = useParams();
  const [searchParams, setSearchParams] = useSearchParams();
  const { message, notification } = App.useApp();

  const initLang = searchParams.get('lang') || 'zh-CN';
  const isCreateMode = !id;
  const pageId = isCreateMode ? null : Number(id);

  const [lang, setLang] = useState(initLang);
  const [formData, setFormData] = useState<PageEditProps>(() => emptyForm(initLang));
  const [languageOptions, setLanguageOptions] = useState<LanguageOption[]>([]);
  const [needTranslate, setNeedTranslate] = useState(false);
  const [saving, setSaving] = useState(false);
  const [contentModelOptions, setContentModelOptions] = useState<
    { label: string; value: number }[]
  >([]);
  const [fieldDefinitions, setFieldDefinitions] = useState<any[]>([]);

  const translateMutation = useTranslate();
  const sections = useMemo(() => formData.sections ?? [], [formData.sections]);

  const goBack = () => navigate('/content/pages');

  const updateForm = useCallback((patch: Partial<PageEditProps>) => {
    setFormData((prev) => ({ ...prev, ...patch }));
  }, []);

  const loadDraft = useCallback(
    (currentLang: string): PageEditProps | null => {
      const raw = localStorage.getItem(getDraftKey(pageId, currentLang, isCreateMode));
      if (!raw) return null;
      try {
        return JSON.parse(raw) as PageEditProps;
      } catch {
        return null;
      }
    },
    [pageId, isCreateMode],
  );

  const saveDraft = () => {
    try {
      localStorage.setItem(
        getDraftKey(pageId, formData.lang, isCreateMode),
        JSON.stringify(formData),
      );
      message.success(t('saveDraftSuccess'));
    } catch (error) {
      console.error('Save draft failed:', error);
      message.error(t('saveDraftFailed'));
    }
  };

  const clearDraft = (currentLang: string) => {
    localStorage.removeItem(getDraftKey(pageId, currentLang, isCreateMode));
  };

  // 内容模型变化时拉取字段定义，驱动动态字段表单
  useEffect(() => {
    if (!formData.contentModelId) {
      setFieldDefinitions([]);
      return;
    }
    (async () => {
      try {
        const resp = await fetchListFieldDefinitions({
          contentModelId: formData.contentModelId,
        });
        setFieldDefinitions(resp.items ?? []);
      } catch (error) {
        console.error('Failed to load field definitions:', error);
        setFieldDefinitions([]);
      }
    })();
  }, [formData.contentModelId]);

  /** 加载页面（编辑模式）：元数据必载，翻译缺失时回退首条并提示需要翻译 */
  const fetchPage = useCallback(
    async (currentLang: string) => {
      if (isCreateMode || !pageId) return;
      try {
        const item = await apiClient.pageService.Get({ id: pageId });

        const hasTranslations = !!item.translations && item.translations.length > 0;
        let langItem = hasTranslations
          ? item.translations?.find((tr) => tr.languageCode === currentLang)
          : undefined;

        let missing = false;
        if (!langItem) {
          langItem = hasTranslations ? item.translations?.[0] : undefined;
          missing = true;
        }

        const availableLanguages = item.availableLanguages || [];
        setLanguageOptions((prev) =>
          prev.map((option) => ({
            ...option,
            hasTranslation: availableLanguages.includes(option.value),
          })),
        );

        const hydratedSections: SectionFormItem[] = (item.sections ?? []).map((s) => ({
          id: s.id,
          type: s.type,
          name: s.name,
          sortOrder: s.sortOrder,
          config: s.config,
          translations: (s.translations ?? []).map((tr) => ({
            id: tr.id,
            languageCode: tr.languageCode,
            content: tr.content,
          })),
        }));

        // 区块重构后正文持久化在 sections 中：从首个正文区块回填编辑器内容。
        // content 契约是 map<string,string>（proto 定义），正文取 value 键，
        // 兼容缺键/历史裸字符串两种形态。
        const bodySection = hydratedSections.find(
          (s) =>
            s.type === 'SECTION_TYPE_MARKDOWN' ||
            s.type === 'SECTION_TYPE_RICH_TEXT',
        );
        const readBodyContent = (
          tr?: SectionFormItemTranslation,
        ): string => {
          const c = tr?.content as Record<string, string> | string | undefined;
          if (!c) return '';
          if (typeof c === 'string') return c;
          if ('value' in c) return c.value ?? '';
          return Object.values(c)[0] ?? '';
        };
        const bodyContent = readBodyContent(
          bodySection?.translations?.find((tr) => tr.languageCode === currentLang) ??
            bodySection?.translations?.[0],
        );

        const base: PageEditProps = {
          id: item.id,
          title: langItem?.title || '',
          slug: langItem?.slug || '',
          content: bodyContent,
          lang: currentLang,
          editorType: item.editorType || EditorType.MARKDOWN,
          parentId: item.parentId ?? undefined,
          type: item.type,
          status: item.status,
          showInNavigation: item.showInNavigation,
          disallowComment: item.disallowComment,
          template: item.template,
          isCustomTemplate: item.isCustomTemplate,
          sortOrder: item.sortOrder,
          contentModelId: item.contentModelId ?? undefined,
          customFields: item.customFields ?? {},
          sections: hydratedSections,
        };

        setNeedTranslate(missing);
        // 草稿覆盖后端数据
        const draft = loadDraft(currentLang);
        setFormData(draft ?? base);
        if (missing) {
          notification.info({ message: t('translationNotExists'), placement: 'bottomRight' });
        }
      } catch (error) {
        console.error('Failed to load page:', error);
        notification.error({ message: t('loadFailed'), placement: 'bottomRight' });
        throw error;
      }
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [isCreateMode, pageId],
  );

  useEffect(() => {
    (async () => {
      try {
        const resp = await fetchListLanguages(
          new PaginationQuery({ orderBy: ['sortOrder'] }),
        );
        setLanguageOptions(
          resp.items?.map((l: any) => ({
            label: l.nativeName || '',
            value: l.languageCode || '',
          })) || [],
        );
      } catch (error) {
        console.error('Failed to load language list:', error);
        notification.error({ message: t('loadLanguageFailed'), placement: 'bottomRight' });
      }
      try {
        const models = await fetchListContentModels(
          new PaginationQuery({ paging: { page: 1, pageSize: 100 } }),
        );
        setContentModelOptions(
          models.items?.map((m: any) => ({
            label: m.translations?.[0]?.name || m.name || `Model #${m.id}`,
            value: m.id,
          })) || [],
        );
      } catch (error) {
        console.error('Failed to load content models:', error);
      }
    })();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    setLang(initLang);
    if (isCreateMode) {
      setNeedTranslate(false);
      setFormData(loadDraft(initLang) ?? emptyForm(initLang));
    } else {
      fetchPage(initLang);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [initLang, isCreateMode, pageId]);

  const handleLanguageChange = async (newLang: string) => {
    setLang(newLang);
    updateForm({ lang: newLang });
    setSearchParams((prev) => {
      const next = new URLSearchParams(prev);
      next.set('lang', newLang);
      return next;
    });
    if (isCreateMode) {
      setNeedTranslate(false);
      setFormData(loadDraft(newLang) ?? emptyForm(newLang));
    } else {
      await fetchPage(newLang);
    }
  };

  const handleTranslate = async () => {
    try {
      const title = await translateMutation.mutateAsync({
        sourceLanguage: 'auto',
        targetLanguage: lang,
        content: formData.title,
      });
      updateForm({ title });
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
      updateForm({ content });
    } catch {
      message.error(t('translateContentFailed'));
    }
  };

  const handleUploadImage = async (
    file: File,
    onProgress?: (percent: number) => void,
  ): Promise<string> => {
    try {
      // 与编辑器内置上传一致：取签名公开访问 URL 供富文本内嵌预览
      const resp = (await uploadFile(
        '',
        '',
        file,
        'post',
        (progressEvent: { loaded?: number; total?: number }) => {
          if (!onProgress) return;
          const total = progressEvent?.total ?? 0;
          onProgress(total > 0 ? Math.floor(((progressEvent.loaded ?? 0) / total) * 100) : -1);
        },
      )) as any;
      return resp?.publicUrl || '';
    } catch (error) {
      console.error('Image upload failed:', error);
      return '';
    }
  };

  // ==============================
  // 嵌套区块管理（页面子部件，随页面整体读写）
  // ==============================
  const setSections = (next: SectionFormItem[]) => updateForm({ sections: next });

  const addSection = () => {
    setSections([
      ...sections,
      {
        type: 'SECTION_TYPE_RICH_TEXT',
        name: '',
        sortOrder: sections.length,
        config: {},
        translations: [],
      },
    ]);
  };

  const removeSection = (index: number) => {
    setSections(sections.filter((_, i) => i !== index));
  };

  const moveSection = (index: number, delta: number) => {
    const target = index + delta;
    if (target < 0 || target >= sections.length) return;
    const next = [...sections];
    [next[index], next[target]] = [next[target], next[index]];
    setSections(next);
  };

  const patchSection = (index: number, patch: Partial<SectionFormItem>) => {
    setSections(sections.map((s, i) => (i === index ? { ...s, ...patch } : s)));
  };

  const addSectionTranslation = (index: number, code: string) => {
    const section = sections[index];
    if (!section || (section.translations ?? []).some((tr) => tr.languageCode === code)) {
      return;
    }
    patchSection(index, {
      translations: [...(section.translations ?? []), { languageCode: code, content: {} }],
    });
  };

  const removeSectionTranslation = (sectionIndex: number, trIndex: number) => {
    const section = sections[sectionIndex];
    if (!section) return;
    patchSection(sectionIndex, {
      translations: (section.translations ?? []).filter((_, i) => i !== trIndex),
    });
  };

  /** JSON textarea 双向桥：解析失败时静默保留旧值（对齐 vben 侧 ignore malformed） */
  const parseJsonRecord = (raw: string): Record<string, string> | undefined => {
    if (!raw) return {};
    try {
      const parsed = JSON.parse(raw);
      return parsed && typeof parsed === 'object' ? parsed : undefined;
    } catch {
      return undefined;
    }
  };

  // ==============================
  // 发布
  // ==============================
  const handlePublish = async () => {
    if (!formData.title) {
      message.error(t('requiredTitle'));
      return;
    }
    if (!formData.slug) {
      message.error(t('requiredSlug'));
      return;
    }

    setSaving(true);
    try {
      // 区块重构后正文持久化在 sections 中：发布前把主编辑器内容合并进
      // 首个正文区块（无则创建），否则编辑器内容会被后端静默丢弃
      const bodyType = String(formData.editorType || '').includes('RICH')
        ? 'SECTION_TYPE_RICH_TEXT'
        : 'SECTION_TYPE_MARKDOWN';
      const nextSections: SectionFormItem[] = sections.map((s) => ({
        ...s,
        translations: [...(s.translations ?? [])],
      }));
      if ((formData.content ?? '').trim() !== '') {
        let body = nextSections.find(
          (s) =>
            s.type === 'SECTION_TYPE_MARKDOWN' ||
            s.type === 'SECTION_TYPE_RICH_TEXT',
        );
        if (!body) {
          body = {
            type: bodyType,
            name: '',
            sortOrder: 0,
            config: {},
            translations: [],
          };
          nextSections.unshift(body);
        }
        const tr = body.translations!.find((tr2) => tr2.languageCode === lang);
        if (tr) {
          tr.content = { value: formData.content };
        } else {
          body.translations!.push({ languageCode: lang, content: { value: formData.content } });
        }
      }

      const hasSectionContent = nextSections.some((s) =>
        (s.translations ?? []).some(
          (tr) => Object.values(tr.content ?? {}).some((v) => String(v ?? '').trim() !== ''),
        ),
      );
      if ((formData.content ?? '').trim() === '' && !hasSectionContent) {
        message.error(t('requiredContent'));
        return;
      }

      const data: Record<string, any> = {
        editorType: formData.editorType,
        parentId: formData.parentId,
        type: formData.type,
        status: formData.status,
        showInNavigation: formData.showInNavigation,
        disallowComment: formData.disallowComment,
        template: formData.template,
        isCustomTemplate: formData.isCustomTemplate,
        sortOrder: formData.sortOrder,
        contentModelId: formData.contentModelId,
        customFields: formData.customFields,
        sections: nextSections,
        // translations 与 sections 均为整体替换字段，不纳入 updateMask
        translations: [
          {
            title: formData.title,
            slug: formData.slug,
            languageCode: lang,
          },
        ],
      };

      if (isCreateMode) {
        await apiClient.pageService.Create({ data: data as any });
      } else {
        await apiClient.pageService.Update({
          id: formData.id || 0,
          data: data as any,
          updateMask: Object.keys(data)
            .filter((k) => k !== 'translations' && k !== 'sections')
            .concat('id')
            .join(','),
        });
      }

      clearDraft(lang);
      message.success(t('publishSuccess'));
      goBack();
    } catch (error) {
      console.error('Failed to publish page:', error);
      message.error(t('publishFailed'));
    } finally {
      setSaving(false);
    }
  };

  const renderCustomField = (field: any) => {
    const name = field.name ?? '';
    const value = formData.customFields?.[name] ?? '';
    switch (field.type) {
      case 'FIELD_TYPE_RICHTEXT':
        return (
          <Input.TextArea
            rows={2}
            value={value}
            size="small"
            onChange={(e) =>
              updateForm({
                customFields: { ...(formData.customFields ?? {}), [name]: e.target.value },
              })
            }
          />
        );
      default:
        return (
          <Input
            size="small"
            value={value}
            onChange={(e) =>
              updateForm({
                customFields: { ...(formData.customFields ?? {}), [name]: e.target.value },
              })
            }
          />
        );
    }
  };

  return (
    <ContentContainer heightMode="fixed" padding="0" bottomMargin={0}>
      <div style={{ height: '100%', display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
        {/* 头部：返回 + 标题 + 别名 + 语言 + 翻译 + 编辑器/状态/类型 */}
        <div style={{ display: 'flex', alignItems: 'center', gap: 8, padding: '12px 16px 8px', flexWrap: 'wrap' }}>
          <Button type="text" icon={<ArrowLeftOutlined />} onClick={goBack} />
          <Input
            placeholder={t('titlePlaceholder')}
            size="large"
            style={{ flex: 1, minWidth: 200 }}
            value={formData.title}
            onChange={(e) => updateForm({ title: e.target.value })}
          />
          <Input
            placeholder={t('slugPlaceholder')}
            size="large"
            style={{ width: 200 }}
            value={formData.slug}
            onChange={(e) => updateForm({ slug: e.target.value })}
          />
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
          {needTranslate && (
            <Button
              type="primary"
              icon={<TranslationOutlined />}
              loading={translateMutation.isPending}
              onClick={handleTranslate}
            >
              {t('oneClickTranslate')}
            </Button>
          )}
          <Select
            style={{ width: 200 }}
            value={formData.editorType}
                      onChange={(v: string) => updateForm({ editorType: v })}
            options={editorTypeOptions(t)}
          />
          <Select
            style={{ width: 200 }}
            value={formData.status}
            onChange={(v) => updateForm({ status: v })}
            options={pageStatusOptions(t)}
          />
          <Select
            style={{ width: 200 }}
            value={formData.type}
            onChange={(v) => updateForm({ type: v })}
            options={pageTypeOptions(t)}
          />
        </div>

        <div style={{ flex: 1, minHeight: 0, display: 'flex', flexDirection: 'column', overflow: 'auto', borderBottom: '1px solid var(--ant-color-split)' }}>
          {/* 内容模型绑定 + 动态字段 */}
          {(contentModelOptions.length > 0 || fieldDefinitions.length > 0) && (
            <div style={{ padding: '8px 16px', borderBottom: '1px solid var(--ant-color-split)' }}>
              <Row gutter={16}>
                <Col span={8}>
                  <div style={{ fontSize: 12, color: 'var(--ant-color-text-secondary)', marginBottom: 4 }}>
                    {t('contentModel')}
                  </div>
                  <Select
                    value={formData.contentModelId}
                    options={contentModelOptions}
                    allowClear
                    size="small"
                    style={{ width: '100%' }}
                    onChange={(v) => updateForm({ contentModelId: v })}
                  />
                </Col>
                {fieldDefinitions.map((field) => (
                  <Col key={field.id} span={8}>
                    <div style={{ fontSize: 12, color: 'var(--ant-color-text-secondary)', marginBottom: 4 }}>
                      {field.label || field.name}
                    </div>
                    {renderCustomField(field)}
                  </Col>
                ))}
              </Row>
            </div>
          )}

          {/* 嵌套区块管理 */}
          <div style={{ padding: '8px 16px', borderBottom: '1px solid var(--ant-color-split)' }}>
            <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 8 }}>
              <span style={{ fontSize: 12, color: 'var(--ant-color-text-secondary)' }}>
                {t('section.moduleName')}
              </span>
              <Button size="small" type="primary" onClick={addSection}>
                {t('section.create')}
              </Button>
            </div>

            {sections.length === 0 && (
              <div style={{ fontSize: 12, color: 'var(--ant-color-text-quaternary)' }}>
                {t('section.contentPlaceholder')}
              </div>
            )}

            {sections.map((section, sIdx) => (
              <div
                key={sIdx}
                style={{ border: '1px solid var(--ant-color-border-secondary)', borderRadius: 6, marginBottom: 8, padding: 8 }}
              >
                <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 8 }}>
                  <Button size="small" icon={<UpOutlined />} onClick={() => moveSection(sIdx, -1)} />
                  <Button size="small" icon={<DownOutlined />} onClick={() => moveSection(sIdx, 1)} />
                  <span style={{ fontSize: 12, color: 'var(--ant-color-text-secondary)' }}>#{sIdx}</span>
                  {section.type && (
                    <span style={{ fontSize: 12, color: 'var(--ant-color-text-tertiary)' }}>
                      {getSectionTypeLabel(t, section.type)}
                    </span>
                  )}
                  <div style={{ flex: 1 }} />
                  <Button size="small" danger onClick={() => removeSection(sIdx)}>
                    ×
                  </Button>
                </div>

                <Row gutter={8} style={{ marginBottom: 8 }}>
                  <Col span={8}>
                    <div style={{ fontSize: 12, color: 'var(--ant-color-text-secondary)', marginBottom: 4 }}>
                      {t('section.type')}
                    </div>
                    <Select
                      value={section.type}
                      options={sectionTypeOptions(t)}
                      size="small"
                      style={{ width: '100%' }}
                      onChange={(v) => patchSection(sIdx, { type: v })}
                    />
                  </Col>
                  <Col span={10}>
                    <div style={{ fontSize: 12, color: 'var(--ant-color-text-secondary)', marginBottom: 4 }}>
                      {t('section.name')}
                    </div>
                    <Input
                      size="small"
                      value={section.name}
                      placeholder={t('section.namePlaceholder')}
                      onChange={(e) => patchSection(sIdx, { name: e.target.value })}
                    />
                  </Col>
                  <Col span={6}>
                    <div style={{ fontSize: 12, color: 'var(--ant-color-text-secondary)', marginBottom: 4 }}>
                      {t('section.sortOrder')}
                    </div>
                    <Input
                      size="small"
                      value={String(section.sortOrder ?? '')}
                      onChange={(e) => patchSection(sIdx, { sortOrder: Number(e.target.value) || 0 })}
                    />
                  </Col>
                </Row>

                <div style={{ marginBottom: 8 }}>
                  <div style={{ fontSize: 12, color: 'var(--ant-color-text-secondary)', marginBottom: 4 }}>
                    {t('section.config')}
                  </div>
                  <Input.TextArea
                    rows={2}
                    size="small"
                    value={section.config ? JSON.stringify(section.config) : ''}
                    onChange={(e) => {
                      const parsed = parseJsonRecord(e.target.value);
                      if (parsed) patchSection(sIdx, { config: parsed });
                    }}
                  />
                </div>

                <div style={{ borderTop: '1px solid var(--ant-color-split)', paddingTop: 8 }}>
                  <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 8 }}>
                    <span style={{ fontSize: 12, color: 'var(--ant-color-text-secondary)' }}>
                      {t('section.content')}
                    </span>
                    <Select
                      value={undefined}
                      size="small"
                      style={{ width: 160 }}
                      placeholder={t('section.addTranslation')}
                      options={languageOptions.map((o) => ({ label: o.label, value: o.value }))}
                      onChange={(code) => code && addSectionTranslation(sIdx, code)}
                    />
                  </div>

                  {(section.translations ?? []).map((tr, trIdx) => (
                    <div
                      key={trIdx}
                      style={{ border: '1px solid var(--ant-color-border-secondary)', borderRadius: 6, marginBottom: 8, padding: 8 }}
                    >
                      <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', marginBottom: 4 }}>
                        <span style={{ fontSize: 12, color: 'var(--ant-color-text-secondary)' }}>
                          {tr.languageCode}
                        </span>
                        <Button
                          size="small"
                          danger
                          onClick={() => removeSectionTranslation(sIdx, trIdx)}
                        >
                          ×
                        </Button>
                      </div>
                      <Input.TextArea
                        rows={3}
                        size="small"
                        value={tr.content ? JSON.stringify(tr.content) : ''}
                        onChange={(e) => {
                          const parsed = parseJsonRecord(e.target.value);
                          if (parsed) {
                            patchSection(sIdx, {
                              translations: (sections[sIdx].translations ?? []).map((item, i) =>
                                i === trIdx ? { ...item, content: parsed } : item,
                              ),
                            });
                          }
                        }}
                      />
                    </div>
                  ))}
                </div>
              </div>
            ))}
          </div>

          {/* 主编辑器 */}
          <div style={{ flex: 1, minHeight: 320 }}>
            <Editor
              height="100%"
              value={formData.content}
              editorType={formData.editorType}
              placeholder={t('contentPlaceholder')}
              uploadImage={handleUploadImage}
              onChange={(v) => updateForm({ content: v })}
            />
          </div>
        </div>

        {/* 底部操作 */}
        <div style={{ display: 'flex', justifyContent: 'flex-end', padding: '8px 16px' }}>
          <Space>
            <Button onClick={saveDraft}>{t('saveDraft')}</Button>
            <Button type="primary" danger loading={saving} onClick={handlePublish}>
              {t('publish')}
            </Button>
          </Space>
        </div>
      </div>
    </ContentContainer>
  );
};

export default PageEditPage;
