import { Button, Collapse, DatePicker, Form, Input, InputNumber, Select, Switch } from 'antd';
import dayjs, { type Dayjs } from 'dayjs';
import { useTranslation } from 'react-i18next';

import PostImageField from './PostImageField';
import type { PostEditProps, SeoMetaForm } from './types';

const SENTENCE_END_REGEX = /[。！？；.!?;]/g;

/**
 * 与后端 summary.GenerateSummaryByRule 对齐的摘要生成：
 * 剥离HTML/Markdown语法 → 合并空白 → 按句子边界截断（100字）→ 补省略号
 */
export function generateSummary(content: string, maxLength = 100): string {
  let text = content.replace(/<[^>]+>/g, ' ');
  text = text.replace(/```[\s\S]*?```/g, ' ');
  text = text.replace(/`[^`]*`/g, ' ');
  text = text.replace(/!\[[^\]]*\]\([^)]*\)/g, ' ');
  text = text.replace(/\[([^\]]*)\]\([^)]*\)/g, '$1');
  text = text.replace(/^[ \t]{0,3}#{1,6}[ \t]+/gm, '');
  text = text.replace(/^[ \t]*>[ \t]?/gm, '');
  text = text.replace(/^[ \t]*([-*_][ \t]*){3,}$/gm, ' ');
  text = text.replace(/^[ \t]*[-+*][ \t]+/gm, '');
  text = text.replace(/^[ \t]*\d+\.[ \t]+/gm, '');
  text = text.replace(/\|/g, ' ');
  text = text.replace(/[*_~]+/g, '');
  text = text.replace(/\s+/g, ' ').trim();

  if (!text) {
    return '暂无摘要';
  }

  const runes = Array.from(text);
  if (runes.length <= maxLength) {
    return text;
  }

  const truncated = runes.slice(0, maxLength).join('');
  const ends = [...truncated.matchAll(new RegExp(SENTENCE_END_REGEX.source, 'g'))];
  if (ends.length === 0) {
    return `${truncated}...`;
  }
  const lastEnd = ends[ends.length - 1]!;
  return `${truncated.slice(0, (lastEnd.index ?? -1) + 1)}...`;
}

interface PostSettingsPanelProps {
  formData: PostEditProps;
  patch: (p: Partial<PostEditProps>) => void;
  categoryOptions: { label: string; value: number }[];
  categoryOptionsLoading: boolean;
}

/**
 * 右栏设置面板：分类（必选，常驻顶部）+ 基础/媒体与高级/SEO 三组折叠
 */
const PostSettingsPanel = ({
  formData,
  patch,
  categoryOptions,
  categoryOptionsLoading,
}: PostSettingsPanelProps) => {
  const { t } = useTranslation('post');
  const seo: SeoMetaForm = formData.seo ?? {};

  const patchSeo = (p: Partial<SeoMetaForm>) => {
    patch({ seo: { ...seo, ...p } });
  };

  const publishTimeValue: Dayjs | null = formData.publishTime
    ? dayjs(formData.publishTime)
    : null;

  // 默认展开基础与媒体组，SEO 组默认收起
  const collapseItems = [
    {
      key: 'base',
      label: t('settingsGroupBase'),
      children: (
        <Form layout="vertical">
          <Form.Item label={t('code')}>
            <Input
              value={formData.code}
              placeholder={t('codePlaceholder')}
              allowClear
              onChange={(e) => patch({ code: e.target.value })}
            />
          </Form.Item>

          <Form.Item label={t('slug')}>
            <Input
              value={formData.slug}
              placeholder={t('slugPlaceholder')}
              allowClear
              onChange={(e) => patch({ slug: e.target.value })}
            />
          </Form.Item>

          <Form.Item
            label={
              <div style={{ display: 'flex', width: '100%', alignItems: 'center', justifyContent: 'space-between' }}>
                <span>{t('summary')}</span>
                <Button
                  size="small"
                  type="link"
                  style={{ padding: '0 4px' }}
                  disabled={!formData.content}
                  onClick={() => patch({ summary: generateSummary(formData.content || '') })}
                >
                  {t('autoSummary')}
                </Button>
              </div>
            }
          >
            <Input.TextArea
              rows={3}
              value={formData.summary}
              placeholder={t('summaryPlaceholder')}
              onChange={(e) => patch({ summary: e.target.value })}
            />
          </Form.Item>
        </Form>
      ),
    },
    {
      key: 'media',
      label: t('settingsGroupMedia'),
      children: (
        <Form layout="vertical">
          <Form.Item label={<span>{t('thumbnail')} <span style={{ fontSize: 12, color: 'var(--ant-color-text-quaternary)' }}>{t('thumbnailSharedHint')}</span></span>}>
            <PostImageField value={formData.thumbnail} onChange={(v) => patch({ thumbnail: v })} />
          </Form.Item>

          <Form.Item label={t('publishTime')}>
            <DatePicker
              showTime
              style={{ width: '100%' }}
              value={publishTimeValue}
              onChange={(value) =>
                patch({ publishTime: value ? value.toISOString() : undefined })
              }
            />
          </Form.Item>

          <Form.Item label={t('sortOrder')}>
            <InputNumber
              min={0}
              precision={0}
              style={{ width: '100%' }}
              value={formData.sortOrder}
              onChange={(v) => patch({ sortOrder: typeof v === 'number' ? v : 0 })}
            />
          </Form.Item>

          <Form.Item label={t('isFeatured')}>
            <Switch
              checked={formData.isFeatured ?? false}
              onChange={(checked) => patch({ isFeatured: checked })}
            />
          </Form.Item>

          <Form.Item label={t('disallowComment')}>
            <Switch
              checked={formData.disallowComment ?? false}
              onChange={(checked) => patch({ disallowComment: checked })}
            />
          </Form.Item>
        </Form>
      ),
    },
    {
      key: 'seo',
      label: t('seoGroup'),
      children: (
        // 护栏：旧本地草稿/异常路径下 seo 可能为空，patchSeo 会从空对象起步
        <Form layout="vertical">
          <Form.Item label={t('seoTitle')}>
            <Input
              value={seo.seoTitle}
              allowClear
              onChange={(e) => patchSeo({ seoTitle: e.target.value })}
            />
          </Form.Item>

          <Form.Item label={t('seoMetaKeywords')}>
            <Input
              value={seo.metaKeywords}
              allowClear
              onChange={(e) => patchSeo({ metaKeywords: e.target.value })}
            />
          </Form.Item>

          <Form.Item label={t('seoMetaDescription')}>
            <Input.TextArea
              rows={3}
              value={seo.metaDescription}
              onChange={(e) => patchSeo({ metaDescription: e.target.value })}
            />
          </Form.Item>

          <Form.Item label={t('seoOgTitle')}>
            <Input
              value={seo.ogTitle}
              allowClear
              onChange={(e) => patchSeo({ ogTitle: e.target.value })}
            />
          </Form.Item>

          <Form.Item label={t('seoOgDescription')}>
            <Input.TextArea
              rows={2}
              value={seo.ogDescription}
              onChange={(e) => patchSeo({ ogDescription: e.target.value })}
            />
          </Form.Item>

          <Form.Item label={t('seoOgImage')}>
            <PostImageField value={seo.ogImage} onChange={(v) => patchSeo({ ogImage: v })} />
          </Form.Item>

          <Form.Item label={t('seoCanonicalUrl')}>
            <Input
              value={seo.canonicalUrl}
              allowClear
              onChange={(e) => patchSeo({ canonicalUrl: e.target.value })}
            />
          </Form.Item>
        </Form>
      ),
    },
  ];

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
      {/* 分类（必选项，常驻侧栏顶部，不随手风琴收起） */}
      <Form layout="vertical" style={{ marginBottom: 0 }}>
        <Form.Item label={t('category')} required style={{ marginBottom: 8 }}>
          <Select
            mode="multiple"
            value={formData.categoryIds}
            options={categoryOptions}
            loading={categoryOptionsLoading}
            placeholder={t('categoryPlaceholder')}
            allowClear
            onChange={(v) => patch({ categoryIds: v })}
          />
        </Form.Item>
      </Form>

      <Collapse
        ghost
        defaultActiveKey={['base', 'media']}
        items={collapseItems}
      />
    </div>
  );
};

export default PostSettingsPanel;
