import { useCallback, useEffect, useState } from 'react';
import { Button, Form, Input, InputNumber, Select, Switch, App } from 'antd';
import {
  ArrowLeftOutlined,
  CheckOutlined,
  CloseOutlined,
  TranslationOutlined,
} from '@ant-design/icons';
import { useNavigate, useParams, useSearchParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';

import { apiClient } from '@/api/client';
import { PaginationQuery } from '@/core';
import {
  useCreateTag,
  useTranslate,
  useUpdateTag,
} from '@/api/hooks/tag';
import { fetchListLanguages } from '@/api/hooks/language';
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import { tagStatusOptions } from '../constants';

interface TagEditProps {
  id?: number;
  name: string;
  slug: string;
  description: string;
  lang: string;
  color?: string;
  icon?: string;
  group?: string;
  coverImage?: string;
  template?: string;
  sortOrder: number;
  isFeatured: boolean;
  status?: string;
}

interface LanguageOption {
  label: string;
  value: string;
  hasTranslation?: boolean;
}

const EMPTY_FORM: TagEditProps = {
  name: '',
  slug: '',
  description: '',
  lang: 'zh-CN',
  sortOrder: 0,
  isFeatured: false,
};

// 草稿存储：key 按 模式+id+语言 区分（对齐 vben StorageManager 语义）
function getDraftKey(tagId: null | number, lang: string, isCreateMode: boolean) {
  return isCreateMode ? `tag-draft-create-${lang}` : `tag-draft-edit-${tagId}-${lang}`;
}

/**
 * 标签编辑页（独立编辑页范式）：
 * 多语言翻译（✓/○ 标记 + 一键翻译）、本地草稿、发布
 */
const TagEditPage = () => {
  const { t } = useTranslation('tag');
  const navigate = useNavigate();
  const { id } = useParams();
  const [searchParams, setSearchParams] = useSearchParams();
  const { message } = App.useApp();
  const [form] = Form.useForm();

  const initLang = searchParams.get('lang') || 'zh-CN';
  const isCreateMode = !id;
  const tagId = isCreateMode ? null : Number(id);

  const [lang, setLang] = useState(initLang);
  const [languageOptions, setLanguageOptions] = useState<LanguageOption[]>([]);
  const [needTranslate, setNeedTranslate] = useState(false);
  const [saving, setSaving] = useState(false);

  const translateMutation = useTranslate();
  const createMutation = useCreateTag();
  const updateMutation = useUpdateTag();

  const goBack = () => navigate('/content/tags');

  // ---- 草稿 ----
  const loadDraft = useCallback(
    (currentLang: string): TagEditProps | null => {
      const raw = localStorage.getItem(getDraftKey(tagId, currentLang, isCreateMode));
      return raw ? (JSON.parse(raw) as TagEditProps) : null;
    },
    [tagId, isCreateMode],
  );

  const saveDraft = (formData: TagEditProps) => {
    localStorage.setItem(
      getDraftKey(tagId, formData.lang, isCreateMode),
      JSON.stringify(formData),
    );
  };

  const clearDraft = (currentLang: string) => {
    localStorage.removeItem(getDraftKey(tagId, currentLang, isCreateMode));
  };

  // ---- 表单值 ----
  const applyFormValues = (values: Partial<TagEditProps>) => {
    form.setFieldsValue({
      name: values.name ?? '',
      slug: values.slug ?? '',
      color: values.color,
      icon: values.icon,
      group: values.group,
      status: values.status,
      sortOrder: values.sortOrder ?? 0,
      isFeatured: values.isFeatured ?? false,
      description: values.description ?? '',
    });
  };

  const getFormValues = (): TagEditProps => {
    const v = form.getFieldsValue();
    return {
      id: tagId ?? undefined,
      name: v.name ?? '',
      slug: v.slug ?? '',
      description: v.description ?? '',
      lang,
      color: v.color,
      icon: v.icon,
      group: v.group,
      sortOrder: v.sortOrder ?? 0,
      isFeatured: v.isFeatured ?? false,
      status: v.status,
    };
  };

  // ---- 拉取标签（编辑模式）----
  const fetchTag = useCallback(
    async (currentLang: string) => {
      if (isCreateMode || !tagId) return;
      try {
        const item = await apiClient.tagService.Get({ id: tagId });
        if (!item?.translations?.length) {
          throw new Error('No translations found for tag');
        }

        let langItem = item.translations.find((tr) => tr.languageCode === currentLang);
        let missing = false;
        if (!langItem) {
          langItem = item.translations[0];
          missing = true;
          setNeedTranslate(true);
        } else {
          setNeedTranslate(false);
        }

        const availableLanguages = item.availableLanguages || [];
        setLanguageOptions((prev) =>
          prev.map((option) => ({
            ...option,
            hasTranslation: availableLanguages.includes(option.value),
          })),
        );

        const base: Partial<TagEditProps> = {
          id: item.id,
          color: item.color ?? undefined,
          icon: item.icon ?? undefined,
          group: item.group ?? undefined,
          sortOrder: item.sortOrder ?? 0,
          isFeatured: item.isFeatured ?? false,
          status: item.status,
          name: langItem.name || '',
          slug: langItem.slug || '',
          description: langItem.description || '',
        };

        // 草稿覆盖后端数据（若存在）
        const draft = loadDraft(currentLang);
        applyFormValues(draft ?? base);
        if (missing) message.info(t('translationNotExists'));
      } catch (error) {
        console.error('Failed to load tag:', error);
        message.error(t('loadFailed'));
      }
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [isCreateMode, tagId],
  );

  // ---- 初始化 ----
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
        message.error(t('loadLanguageFailed'));
      }
    })();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    setLang(initLang);
    if (isCreateMode) {
      const draft = loadDraft(initLang);
      applyFormValues(draft ?? { ...EMPTY_FORM, lang: initLang });
    } else {
      fetchTag(initLang);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [initLang, isCreateMode, tagId]);

  // ---- 语言切换 ----
  const handleLanguageChange = async (newLang: string) => {
    setLang(newLang);
    setSearchParams((prev) => {
      const next = new URLSearchParams(prev);
      next.set('lang', newLang);
      return next;
    });
    if (!isCreateMode) {
      await fetchTag(newLang);
    } else {
      const draft = loadDraft(newLang);
      applyFormValues(draft ?? { ...EMPTY_FORM, lang: newLang });
    }
  };

  // ---- 一键翻译 ----
  const handleTranslate = async () => {
    const values = getFormValues();
    try {
      const name = await translateMutation.mutateAsync({
        sourceLanguage: 'auto',
        targetLanguage: lang,
        content: values.name,
      });
      form.setFieldValue('name', name);
    } catch {
      message.error(t('translateNameFailed'));
      return;
    }
    try {
      const description = await translateMutation.mutateAsync({
        sourceLanguage: 'auto',
        targetLanguage: lang,
        content: values.description,
      });
      form.setFieldValue('description', description);
    } catch {
      message.error(t('translateDescriptionFailed'));
    }
  };

  // ---- 保存草稿 ----
  const handleSaveDraft = () => {
    try {
      saveDraft({ ...getFormValues(), lang });
      message.success(t('saveDraftSuccess'));
    } catch (error) {
      console.error('Save draft failed:', error);
      message.error(t('saveDraftFailed'));
    }
  };

  // ---- 发布 ----
  const handlePublish = async () => {
    const values = getFormValues();
    if (!values.name) {
      message.error(t('requiredName'));
      return;
    }
    if (!values.slug) {
      message.error(t('requiredSlug'));
      return;
    }

    setSaving(true);
    try {
      const data = {
        color: values.color,
        icon: values.icon,
        group: values.group,
        sortOrder: values.sortOrder,
        isFeatured: values.isFeatured,
        status: values.status,
        translations: [
          {
            name: values.name,
            slug: values.slug,
            description: values.description,
            languageCode: values.lang,
          },
        ],
      };

      if (isCreateMode) {
        await createMutation.mutateAsync({ data: data as any });
      } else {
        await updateMutation.mutateAsync({
          id: values.id || 0,
          data: data as any,
          // translations 走整体 upsert，其余标量字段进 mask
          excludeMaskKeys: ['translations'],
        });
      }

      clearDraft(lang);
      message.success(t('publishSuccess'));
      goBack();
    } catch (error) {
      console.error('Failed to publish tag:', error);
      message.error(t('publishFailed'));
    } finally {
      setSaving(false);
    }
  };

  return (
    <ContentContainer heightMode="fixed" padding="16px" bottomMargin={0}>
      <Form
        form={form}
        layout="vertical"
        style={{ height: '100%', display: 'flex', flexDirection: 'column', overflow: 'hidden' }}
        initialValues={{ sortOrder: 0, isFeatured: false }}
      >
        {/* 头部：返回 + 名称 + 语言切换 + 翻译/草稿/发布 */}
        <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 16 }}>
          <Button type="text" icon={<ArrowLeftOutlined />} onClick={goBack} />
          <Form.Item name="name" noStyle>
            <Input placeholder={t('namePlaceholder')} size="large" style={{ flex: 1 }} />
          </Form.Item>
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
          <Button onClick={handleSaveDraft}>{t('saveDraft')}</Button>
          <Button type="primary" danger loading={saving} onClick={handlePublish}>
            {t('publish')}
          </Button>
        </div>

        {/* 表单体 */}
        <div style={{ flex: 1, minHeight: 0, overflow: 'auto', maxWidth: 960, margin: '0 auto', width: '100%' }}>
          <div style={{ display: 'flex', gap: 16 }}>
            <Form.Item
              name="slug"
              label={t('slugPlaceholder')}
              rules={[{ required: true, message: t('requiredSlug') }]}
              style={{ flex: 1 }}
            >
              <Input placeholder={t('slugPlaceholder')} />
            </Form.Item>
            <Form.Item name="color" label={t('color')} style={{ flex: 1 }}>
              <Input placeholder={t('colorPlaceholder')} />
            </Form.Item>
          </div>

          <div style={{ display: 'flex', gap: 16 }}>
            <Form.Item name="icon" label={t('icon')} style={{ flex: 1 }}>
              <Input placeholder={t('iconPlaceholder')} />
            </Form.Item>
            <Form.Item name="group" label={t('group')} style={{ flex: 1 }}>
              <Input placeholder={t('groupPlaceholder')} />
            </Form.Item>
          </div>

          <div style={{ display: 'flex', gap: 16 }}>
            <Form.Item name="status" label={t('status')} style={{ flex: 1 }}>
              <Select options={tagStatusOptions(t)} placeholder={t('statusPlaceholder')} />
            </Form.Item>
            <Form.Item name="sortOrder" label={t('sortOrder')} style={{ flex: 1 }}>
              <InputNumber min={0} style={{ width: '100%' }} />
            </Form.Item>
            <Form.Item name="isFeatured" label={t('isFeatured')} valuePropName="checked">
              <Switch />
            </Form.Item>
          </div>

          <Form.Item name="description" label={t('descriptionPlaceholder')}>
            <Input.TextArea rows={4} placeholder={t('descriptionPlaceholder')} />
          </Form.Item>
        </div>
      </Form>
    </ContentContainer>
  );
};

export default TagEditPage;
