import { useCallback, useEffect, useState } from 'react';
import { Button, Col, Form, Input, InputNumber, Row, Select, Switch, App } from 'antd';
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
import ContentContainer from '@/layouts/components/PageContainer/ContentContainer';
import {
  useCreateCategory,
  useUpdateCategory,
} from '@/api/hooks/category';
import {
  fetchListContentModels,
  fetchListFieldDefinitions,
} from '@/api/hooks/content-model';
import { fetchListLanguages } from '@/api/hooks/language';
import { useTranslate } from '@/api/hooks/tag';
import { categoryStatusOptions } from '../constants';

interface CategoryEditProps {
  id?: number;
  name: string;
  slug: string;
  description: string;
  lang: string;
  parentId?: number;
  icon?: string;
  isNav: boolean;
  sortOrder: number;
  status?: string;
  contentModelId?: number;
  customFields?: Record<string, any>;
}

interface LanguageOption {
  label: string;
  value: string;
  hasTranslation?: boolean;
}

const EMPTY_FORM: CategoryEditProps = {
  name: '',
  slug: '',
  description: '',
  lang: 'zh-CN',
  sortOrder: 0,
  isNav: false,
};

function getDraftKey(categoryId: null | number, lang: string, isCreateMode: boolean) {
  return isCreateMode
    ? `category-draft-create-${lang}`
    : `category-draft-edit-${categoryId}-${lang}`;
}

/**
 * 分类编辑页（独立编辑页范式）：
 * 多语言翻译（✓/○ 标记 + 一键翻译）、本地草稿、内容模型绑定 + 动态字段、发布
 */
const CategoryEditPage = () => {
  const { t } = useTranslation('category');
  const navigate = useNavigate();
  const { id } = useParams();
  const [searchParams, setSearchParams] = useSearchParams();
  const { message } = App.useApp();
  const [form] = Form.useForm();

  const initLang = searchParams.get('lang') || 'zh-CN';
  const isCreateMode = !id;
  const categoryId = isCreateMode ? null : Number(id);

  const [lang, setLang] = useState(initLang);
  const [languageOptions, setLanguageOptions] = useState<LanguageOption[]>([]);
  const [needTranslate, setNeedTranslate] = useState(false);
  const [saving, setSaving] = useState(false);
  const [contentModelOptions, setContentModelOptions] = useState<
    { label: string; value: number }[]
  >([]);
  const [fieldDefinitions, setFieldDefinitions] = useState<any[]>([]);

  const translateMutation = useTranslate();
  const createMutation = useCreateCategory();
  const updateMutation = useUpdateCategory();

  const goBack = () => navigate('/content/categories');

  const loadDraft = useCallback(
    (currentLang: string): CategoryEditProps | null => {
      const raw = localStorage.getItem(
        getDraftKey(categoryId, currentLang, isCreateMode),
      );
      return raw ? (JSON.parse(raw) as CategoryEditProps) : null;
    },
    [categoryId, isCreateMode],
  );

  const saveDraft = (formData: CategoryEditProps) => {
    localStorage.setItem(
      getDraftKey(categoryId, formData.lang, isCreateMode),
      JSON.stringify(formData),
    );
  };

  const clearDraft = (currentLang: string) => {
    localStorage.removeItem(getDraftKey(categoryId, currentLang, isCreateMode));
  };

  const applyFormValues = (values: Partial<CategoryEditProps>) => {
    form.setFieldsValue({
      name: values.name ?? '',
      slug: values.slug ?? '',
      icon: values.icon,
      status: values.status,
      sortOrder: values.sortOrder ?? 0,
      isNav: values.isNav ?? false,
      description: values.description ?? '',
      contentModelId: values.contentModelId,
      customFields: values.customFields ?? {},
    });
  };

  const getFormValues = (): CategoryEditProps => {
    const v = form.getFieldsValue();
    return {
      id: categoryId ?? undefined,
      name: v.name ?? '',
      slug: v.slug ?? '',
      description: v.description ?? '',
      lang,
      parentId: v.parentId,
      icon: v.icon,
      isNav: v.isNav ?? false,
      sortOrder: v.sortOrder ?? 0,
      status: v.status,
      contentModelId: v.contentModelId,
      customFields: v.customFields ?? {},
    };
  };

  // 内容模型变化时拉取字段定义，驱动动态字段表单
  const contentModelId = Form.useWatch('contentModelId', form);
  useEffect(() => {
    if (!contentModelId) {
      setFieldDefinitions([]);
      return;
    }
    (async () => {
      try {
        const resp = await fetchListFieldDefinitions({ contentModelId });
        setFieldDefinitions(resp.items ?? []);
      } catch (error) {
        console.error('Failed to load field definitions:', error);
        setFieldDefinitions([]);
      }
    })();
  }, [contentModelId]);

  const fetchCategory = useCallback(
    async (currentLang: string) => {
      if (isCreateMode || !categoryId) return;
      try {
        const item = await apiClient.categoryService.Get({ id: categoryId });
        if (!item?.translations?.length) {
          throw new Error('No translations found for category');
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

        const base: Partial<CategoryEditProps> = {
          id: item.id,
          parentId: item.parentId ?? undefined,
          icon: item.icon ?? undefined,
          isNav: item.isNav ?? false,
          sortOrder: item.sortOrder ?? 0,
          status: item.status,
          contentModelId: item.contentModelId ?? undefined,
          customFields: item.customFields ?? {},
          name: langItem.name || '',
          slug: langItem.slug || '',
          description: langItem.description || '',
        };

        const draft = loadDraft(currentLang);
        applyFormValues(draft ?? base);
        if (missing) message.info(t('translationNotExists'));
      } catch (error) {
        console.error('Failed to load category:', error);
        message.error(t('loadFailed'));
      }
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [isCreateMode, categoryId],
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
        message.error(t('loadLanguageFailed'));
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
      const draft = loadDraft(initLang);
      applyFormValues(draft ?? { ...EMPTY_FORM, lang: initLang });
    } else {
      fetchCategory(initLang);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [initLang, isCreateMode, categoryId]);

  const handleLanguageChange = async (newLang: string) => {
    setLang(newLang);
    setSearchParams((prev) => {
      const next = new URLSearchParams(prev);
      next.set('lang', newLang);
      return next;
    });
    if (!isCreateMode) {
      await fetchCategory(newLang);
    } else {
      const draft = loadDraft(newLang);
      applyFormValues(draft ?? { ...EMPTY_FORM, lang: newLang });
    }
  };

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

  const handleSaveDraft = () => {
    try {
      saveDraft({ ...getFormValues(), lang });
      message.success(t('saveDraftSuccess'));
    } catch (error) {
      console.error('Save draft failed:', error);
      message.error(t('saveDraftFailed'));
    }
  };

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
      const data: Record<string, any> = {
        parentId: values.parentId,
        icon: values.icon,
        isNav: values.isNav,
        sortOrder: values.sortOrder,
        // "发布"语义：未显式选择状态时默认启用，避免落库为空状态
        status: values.status || 'CATEGORY_STATUS_ACTIVE',
        contentModelId: values.contentModelId,
        customFields: values.customFields,
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
        await createMutation.mutateAsync({ data });
      } else {
        await updateMutation.mutateAsync({
          id: values.id || 0,
          data,
          excludeMaskKeys: ['translations'],
        });
      }

      clearDraft(lang);
      message.success(t('publishSuccess'));
      goBack();
    } catch (error) {
      console.error('Failed to publish category:', error);
      message.error(t('publishFailed'));
    } finally {
      setSaving(false);
    }
  };

  const renderCustomField = (field: any) => {
    switch (field.type) {
      case 'FIELD_TYPE_RICHTEXT':
        return (
          <Input.TextArea
            rows={4}
            placeholder={field.label || field.name}
          />
        );
      // RELATION 关联选择器待 content/model 子模块移植后补齐，暂以文本输入兜底
      default:
        return <Input placeholder={field.label || field.name} />;
    }
  };

  return (
    <ContentContainer heightMode="fixed" padding="16px" bottomMargin={0}>
      <Form
        form={form}
        layout="vertical"
        style={{ height: '100%', display: 'flex', flexDirection: 'column', overflow: 'hidden' }}
        initialValues={{ sortOrder: 0, isNav: false }}
      >
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

        <div style={{ flex: 1, minHeight: 0, overflow: 'auto', maxWidth: 960, margin: '0 auto', width: '100%' }}>
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="slug"
                label={t('slugPlaceholder')}
                rules={[{ required: true, message: t('requiredSlug') }]}
              >
                <Input placeholder={t('slugPlaceholder')} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="icon" label={t('iconPlaceholder')}>
                <Input placeholder={t('iconPlaceholder')} />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item name="status" label={t('status')}>
                <Select options={categoryStatusOptions(t)} placeholder={t('statusPlaceholder')} />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item name="contentModelId" label={t('contentModel')}>
                <Select
                  options={contentModelOptions}
                  placeholder={t('contentModelPlaceholder')}
                  allowClear
                />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={8}>
              <Form.Item name="sortOrder" label={t('sortOrder')}>
                <InputNumber min={0} style={{ width: '100%' }} />
              </Form.Item>
            </Col>
            <Col span={8}>
              <Form.Item name="isNav" label={t('isNav')} valuePropName="checked">
                <Switch />
              </Form.Item>
            </Col>
          </Row>

          <Form.Item name="description" label={t('descriptionPlaceholder')}>
            <Input.TextArea rows={4} placeholder={t('descriptionPlaceholder')} />
          </Form.Item>

          {/* 动态字段表单：由绑定的内容模型字段定义驱动 */}
          <Form.Item
            noStyle
            shouldUpdate={(prev, cur) => prev.contentModelId !== cur.contentModelId}
          >
            {() => {
              return fieldDefinitions.map((field) => (
                <Form.Item
                  key={field.id}
                  name={['customFields', field.name ?? '']}
                  label={field.label || field.name}
                >
                  {renderCustomField(field)}
                </Form.Item>
              ));
            }}
          </Form.Item>
        </div>
      </Form>
    </ContentContainer>
  );
};

export default CategoryEditPage;
