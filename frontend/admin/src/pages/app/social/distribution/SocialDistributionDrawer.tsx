import { useEffect, useMemo, useRef, useState } from 'react';
import type { ProFormInstance } from '@ant-design/pro-components';
import {
  DrawerForm,
  ProFormDateTimePicker,
  ProFormSelect,
  ProFormText,
  ProFormTextArea,
} from '@ant-design/pro-components';
import { App } from 'antd';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import dayjs from 'dayjs';

import AiGenerateButton from '@/components/common/AiGenerateButton';

import { fetchListPosts } from '@/api/hooks/post';
import { PaginationQuery } from '@/core';
import { fetchListMediaAssets } from '@/api/hooks/media-asset';
import { fetchListSocialAccounts, useCreateSocialDistribution, useUpdateSocialDistribution } from '@/api/hooks/social';
import { MAX_POST_CHARS, supportedPlatformOptions } from './constants';

/** 支持配图的平台 */
const MEDIA_PLATFORMS = new Set([
  'PLATFORM_TELEGRAM',
  'PLATFORM_BLUESKY',
  'PLATFORM_MASTODON',
  'PLATFORM_WORDPRESS',
  'PLATFORM_WEIBO',
]);

interface SocialDistributionDrawerProps {
  open: boolean;
  data?: any;
  /** 新建模式的预填（文章编辑页"分发到社交"/日历点日期入口带入） */
  prefill?: {
    postId?: number;
    postTitle?: string;
    overrideText?: string;
    publishTime?: string;
  };
  onClose: () => void;
  onSuccess: () => void;
}

/**
 * 内容分发抽屉：选源帖（CMS 文章）+ 平台账号 + 覆盖文案 + 排期。
 * 发布时间留空 = 保存后立即发布；状态机与重试由后端发布引擎负责。
 */
const SocialDistributionDrawer: React.FC<SocialDistributionDrawerProps> = ({
  open,
  data,
  prefill,
  onClose,
  onSuccess,
}) => {
  const { t } = useTranslation('social');
  const formRef = useRef<ProFormInstance>(null);
  const queryClient = useQueryClient();
  const { message } = App.useApp();

  const [confirmLoading, setConfirmLoading] = useState(false);
  const [platform, setPlatform] = useState<string | undefined>();

  const isEdit = !!data?.id;
  const prefillPostId = prefill?.postId;
  const prefillPostTitle = prefill?.postTitle;
  // 已发布/发布中的记录不可编辑
  const immutable =
    isEdit && ['STATUS_PUBLISHED', 'STATUS_PUBLISHING'].includes(data?.status);

  useEffect(() => {
    if (open) {
      setTimeout(() => {
        setPlatform(data?.platform);
        formRef.current?.setFieldsValue({
          platform: data?.platform,
          accountId: data?.accountId,
          postId: data?.postId,
          title: data?.title || '',
          mediaIds: data?.mediaIds || [],
          overrideText: data?.overrideText || '',
          publishTime: data?.publishTime ? dayjs(data.publishTime) : undefined,
        });
      }, 0);
    }
  }, [open, data]);

  // 平台账号下拉（启用中的账号）
  const accountOptions = useMemo(() => {
    return async (keywords?: string) => {
      const query = new PaginationQuery({
        paging: { page: 1, pageSize: 100 },
        formValues: {
          ...(platform ? { platform } : {}),
          ...(keywords ? { name: keywords } : {}),
        },
        orderBy: ['-created_at'],
      });
      const resp = await fetchListSocialAccounts(query);
      return (resp.items || [])
        .filter((item: any) => item.isEnabled !== false)
        .map((item: any) => ({
          label: `${item.name}（${item.handle || item.platform}）`,
          value: item.id,
        }));
    };
  }, [platform]);

  // 配图下拉（媒体库中的图片资产）
  const mediaOptions = useMemo(() => {
    return async (keywords?: string) => {
      const query = new PaginationQuery({
        paging: { page: 1, pageSize: 100 },
        formValues: keywords ? { filename: keywords } : null,
        orderBy: ['-created_at'],
      });
      const resp = await fetchListMediaAssets(query);
      return (resp.items || [])
        .filter((item: any) => (item.mimeType || '').startsWith('image/'))
        .map((item: any) => ({
          label: `${item.filename || item.title || `#${item.id}`}`,
          value: item.id,
        }));
    };
  }, []);

  // 源帖下拉（CMS 文章，取标题做选项文案）
  const postOptions = useMemo(() => {
    return async (keywords?: string) => {
      const query = new PaginationQuery({
        paging: { page: 1, pageSize: 50 },
        formValues: keywords ? { title: keywords } : null,
        orderBy: ['-created_at'],
      });
      const resp = await fetchListPosts(query);
      const options = (resp.items || []).map((item: any) => {
        const title = item.translations?.find((tr: any) => tr.title)?.title;
        return {
          label: `${title || item.code || `#${item.id}`}`,
          value: item.id,
        };
      });
      // 预填的源帖可能不在最近列表里，前置补一条保证回显有 label
      if (prefillPostId && !options.some((o: any) => o.value === prefillPostId)) {
        options.unshift({
          label: prefillPostTitle || `#${prefillPostId}`,
          value: prefillPostId,
        });
      }
      return options;
    };
  }, [prefillPostId, prefillPostTitle]);

  const createMutation = useCreateSocialDistribution({
    onSuccess: () => {
      message.success(t('createSuccessDist'));
      queryClient.invalidateQueries({ queryKey: ['listSocialDistributions'] });
      onSuccess();
      onClose();
    },
    onError: (error: Error) => {
      message.error(error.message || t('createFailedDist'));
    },
  });

  const updateMutation = useUpdateSocialDistribution({
    onSuccess: () => {
      message.success(t('updateSuccessDist'));
      queryClient.invalidateQueries({ queryKey: ['listSocialDistributions'] });
      onSuccess();
      onClose();
    },
    onError: (error: Error) => {
      message.error(error.message || t('updateFailedDist'));
    },
  });

  const handleSubmit = async (values: Record<string, any>) => {
    try {
      setConfirmLoading(true);
      const publishTime = values.publishTime
        ? dayjs(values.publishTime).toISOString()
        : undefined;

      // 配图仅支持特定平台：不支持时明确报错（后端也有同样校验兜底）
      if (values.mediaIds?.length && !MEDIA_PLATFORMS.has(values.platform)) {
        message.error(t('mediaUnsupported'));
        return false;
      }
      const mediaIds: number[] = values.mediaIds || [];

      if (isEdit) {
        const maskKeys: string[] = ['overrideText', 'publishTime'];
        const payload: Record<string, any> = {
          overrideText: values.overrideText,
          publishTime: publishTime ?? null,
        };
        if (platform === 'PLATFORM_DEVTO') {
          payload.title = values.title || null;
          maskKeys.push('title');
        }
        if (values.accountId && values.accountId !== data.accountId) {
          payload.accountId = values.accountId;
          maskKeys.push('accountId');
        }
        if (values.postId && values.postId !== data.postId) {
          payload.postId = values.postId;
          maskKeys.push('postId');
        }
        if (MEDIA_PLATFORMS.has(platform || '')) {
          payload.mediaIds = mediaIds;
          maskKeys.push('mediaIds');
        }
        await updateMutation.mutateAsync({ id: data.id, values: payload, maskKeys });
      } else {
        await createMutation.mutateAsync({
          data: {
            platform: values.platform,
            accountId: values.accountId,
            postId: values.postId,
            title: values.title,
            mediaIds,
            overrideText: values.overrideText,
            // 缺省 SCHEDULED：publishTime 为空时后端立即入队
            status: 'STATUS_SCHEDULED',
            publishTime,
          } as any,
        });
      }
      return true;
    } catch {
      return false;
    } finally {
      setConfirmLoading(false);
    }
  };

  return (
    <DrawerForm
      formRef={formRef}
      title={isEdit ? t('editDistribution') : t('createDistribution')}
      open={open}
      onOpenChange={(visible) => {
        if (!visible) {
          formRef.current?.resetFields();
          onClose();
        }
      }}
      onFinish={handleSubmit}
      submitter={{
        searchConfig: {
          submitText: t('common:button.submit'),
          resetText: t('common:button.cancel'),
        },
        submitButtonProps: {
          loading: confirmLoading || createMutation.isPending || updateMutation.isPending,
          disabled: immutable,
        },
        resetButtonProps: { onClick: onClose },
      }}
      drawerProps={{ destroyOnHidden: true, onClose, size: 600 }}
    >
      <ProFormSelect
        name="platform"
        label={t('platform')}
        placeholder={t('platformPlaceholder')}
        rules={[{ required: true, message: t('requiredPlatform') }]}
        options={supportedPlatformOptions(t)}
        disabled={isEdit}
        fieldProps={{
          allowClear: false,
          onChange: (value: string) => {
            setPlatform(value);
            formRef.current?.setFieldValue('accountId', undefined);
          },
        }}
      />

      <ProFormSelect
        name="accountId"
        label={t('socialAccount')}
        placeholder={platform ? t('socialAccountPlaceholder') : t('socialAccountNeedPlatform')}
        rules={[{ required: true, message: t('requiredSocialAccount') }]}
        request={() => accountOptions()}
        params={[platform]}
        disabled={isEdit}
        fieldProps={{ showSearch: true, optionFilterProp: 'label' }}
      />

      <ProFormSelect
        name="postId"
        label={t('sourcePost')}
        placeholder={t('sourcePostPlaceholder')}
        extra={t('helpSourcePost')}
        request={() => postOptions()}
        fieldProps={{ showSearch: true, optionFilterProp: 'label', allowClear: true }}
      />

      {platform === 'PLATFORM_DEVTO' && (
        <ProFormText
          name="title"
          label={t('title')}
          placeholder={t('titlePlaceholder')}
          rules={[{ required: true, message: t('requiredTitle') }]}
          extra={t('helpTitle')}
          fieldProps={{ allowClear: true, maxLength: 128, showCount: true }}
        />
      )}

      {MEDIA_PLATFORMS.has(platform || '') && (
        <ProFormSelect
          name="mediaIds"
          label={t('mediaIds')}
          placeholder={t('mediaIdsPlaceholder')}
          extra={t('helpMediaIds')}
          request={() => mediaOptions()}
          fieldProps={{
            mode: 'multiple',
            showSearch: true,
            optionFilterProp: 'label',
            allowClear: true,
            maxCount: 4,
          }}
        />
      )}

      <ProFormTextArea
        name="overrideText"
        label={
          <span style={{ display: 'inline-flex', alignItems: 'center', gap: 8 }}>
            {t('overrideText')}
            <AiGenerateButton
              scene="SOCIAL_POST"
              context={platform ? `${t('platform')}：${t(`platformMap.${platform}`)}` : undefined}
              defaultTopic={data?.postTitle || undefined}
              onGenerate={(content: string) => formRef.current?.setFieldValue('overrideText', content)}
            />
          </span>
        }
        placeholder={t('overrideTextPlaceholder')}
        rules={[{ required: true, message: t('requiredOverrideText') }]}
        extra={t('helpOverrideText', { max: MAX_POST_CHARS })}
        fieldProps={{
          rows: 5,
          showCount: true,
          maxLength: MAX_POST_CHARS,
          allowClear: true,
        }}
      />

      <ProFormDateTimePicker
        name="publishTime"
        label={t('publishTime')}
        placeholder={t('publishTimePlaceholder')}
        extra={t('helpPublishTime')}
        fieldProps={{ style: { width: '100%' }, showTime: { format: 'HH:mm' }, format: 'YYYY-MM-DD HH:mm' }}
        disabled={immutable}
      />
    </DrawerForm>
  );
};

export default SocialDistributionDrawer;
