import { useRef, useState, useEffect } from 'react';
import type { siteservicev1_Site } from '@/api/generated/admin/service/v1';
import type { ProFormInstance } from '@ant-design/pro-components';
import {
  DrawerForm,
  ProFormText,
  ProFormSelect,
  ProFormSwitch,
  ProFormTextArea,
} from '@ant-design/pro-components';
import { App } from 'antd';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';

import { useCreateSite, useUpdateSite } from '@/api/hooks/site';
import { fetchListLanguages } from '@/api/hooks/language';
import { PaginationQuery } from '@/core';
import { statusOptions } from './constants';

interface SiteDrawerProps {
  open: boolean;
  mode: 'create' | 'edit';
  data?: any;
  onClose: () => void;
  onSuccess: () => void;
}

/**
 * 站点编辑/创建抽屉组件
 */
const SiteDrawer: React.FC<SiteDrawerProps> = ({
  open,
  mode,
  data,
  onClose,
  onSuccess,
}) => {
  const { t } = useTranslation('site');
  const formRef = useRef<ProFormInstance>(null);
  const queryClient = useQueryClient();
  const { message } = App.useApp();

  const [confirmLoading, setConfirmLoading] = useState(false);
  const [languageOptions, setLanguageOptions] = useState<
    { label: string; value: string }[]
  >([]);

  useEffect(() => {
    if (!open) return;
    (async () => {
      try {
        const resp = await fetchListLanguages(
          new PaginationQuery({ orderBy: ['sortOrder'] }),
        );
        setLanguageOptions(
          resp.items?.map((lang: any) => ({
            label: lang.nativeName || lang.languageCode || '',
            value: lang.languageCode || '',
          })) || [],
        );
      } catch (error) {
        console.error('Failed to load language list:', error);
      }
    })();
  }, [open]);

  // 编辑模式下设置表单值（destroyOnHidden 时需延迟赋值）
  useEffect(() => {
    if (open && mode === 'edit' && data) {
      setTimeout(() => {
        formRef.current?.setFieldsValue({
          name: data.name || '',
          slug: data.slug || '',
          domain: data.domain || '',
          status: data.status || 'SITE_STATUS_ACTIVE',
          isDefault: data.isDefault ?? false,
          defaultLocale: data.defaultLocale || 'zh-CN',
          template: data.template || '',
          theme: data.theme || '',
          // 桥接回填：alternateDomains string[] → 换行分隔文本
          alternateDomainsText: Array.isArray(data.alternateDomains)
            ? data.alternateDomains.join('\n')
            : '',
        });
      }, 0);
    }
  }, [open, mode, data]);

  // 创建 mutation
  const createMutation = useCreateSite({
    onSuccess: () => {
      message.success(t('createSuccess'));
      queryClient.invalidateQueries({ queryKey: ['listSites'] });
      onSuccess();
      onClose();
    },
    onError: (error: Error) => {
      message.error(error.message || t('createFailed'));
    },
  });

  // 更新 mutation
  const updateMutation = useUpdateSite({
    onSuccess: () => {
      message.success(t('updateSuccess'));
      queryClient.invalidateQueries({ queryKey: ['listSites'] });
      onSuccess();
      onClose();
    },
    onError: (error: Error) => {
      message.error(error.message || t('updateFailed'));
    },
  });

  // 提交表单
  const handleSubmit = async (values: Record<string, any>) => {
    try {
      setConfirmLoading(true);
      // 桥接：alternateDomainsText（换行分隔文本）→ alternateDomains（string[]）
      const alternateDomains = String(values.alternateDomainsText ?? '')
        .split('\n')
        .map((s) => s.trim())
        .filter(Boolean);

      // 显式白名单构造 payload：row 回填残留（id/updatedBy/createdAt 等）或
      // 表单辅助字段（alternateDomainsText）一旦混入 updateMask，proto 校验会整体拒绝
      const payload: Record<string, any> = {
        name: values.name,
        slug: values.slug,
        domain: values.domain,
        status: values.status,
        isDefault: values.isDefault,
        defaultLocale: values.defaultLocale,
        template: values.template,
        theme: values.theme,
        alternateDomains,
      };

      if (mode === 'edit' && data?.id) {
        await updateMutation.mutateAsync({ id: data.id, values: payload as unknown as siteservicev1_Site });
      } else {
        await createMutation.mutateAsync({ data: payload as unknown as siteservicev1_Site });
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
      title={mode === 'create' ? t('create') : t('edit')}
      open={open}
      onOpenChange={(visible) => {
        if (!visible) {
          formRef.current?.resetFields();
          onClose();
        }
      }}
      initialValues={{
        status: 'SITE_STATUS_ACTIVE',
        isDefault: false,
        defaultLocale: 'zh-CN',
      }}
      onFinish={handleSubmit}
      submitter={{
        searchConfig: {
          submitText: t('common:button.submit'),
          resetText: t('common:button.cancel'),
        },
        submitButtonProps: {
          loading: confirmLoading || createMutation.isPending || updateMutation.isPending,
        },
        resetButtonProps: { onClick: onClose },
      }}
      drawerProps={{ destroyOnHidden: true, onClose, placement: 'left', size: 800 }}
    >
      <ProFormText
        name="name"
        label={t('name')}
        placeholder={t('namePlaceholder')}
        rules={[{ required: true, message: t('requiredName') }]}
        fieldProps={{ allowClear: true }}
      />

      <ProFormText
        name="slug"
        label={t('slug')}
        placeholder={t('slugPlaceholder')}
        rules={[{ required: true, message: t('requiredSlug') }]}
        fieldProps={{ allowClear: true }}
      />

      <ProFormText
        name="domain"
        label={t('domain')}
        placeholder={t('domainPlaceholder')}
        fieldProps={{ allowClear: true }}
      />

      <ProFormSelect
        name="status"
        label={t('status')}
        placeholder={t('statusPlaceholder')}
        rules={[{ required: true, message: t('requiredStatus') }]}
        options={statusOptions(t)}
        fieldProps={{ allowClear: true }}
      />

      <ProFormSwitch
        name="isDefault"
        label={t('isDefault')}
        help={t('helpIsDefault')}
      />

      <ProFormSelect
        name="defaultLocale"
        label={t('defaultLocale')}
        placeholder={t('defaultLocalePlaceholder')}
        options={languageOptions}
        fieldProps={{ allowClear: true }}
        help={t('helpDefaultLocale')}
      />

      <ProFormText
        name="template"
        label={t('template')}
        placeholder={t('templatePlaceholder')}
        rules={[{ required: true, message: t('requiredTemplate') }]}
        fieldProps={{ allowClear: true }}
      />

      <ProFormText
        name="theme"
        label={t('theme')}
        placeholder={t('themePlaceholder')}
        fieldProps={{ allowClear: true }}
      />

      <ProFormTextArea
        name="alternateDomainsText"
        label={t('alternateDomains')}
        placeholder={t('placeholderAlternateDomains')}
        help={t('helpAlternateDomains')}
        fieldProps={{ rows: 3, allowClear: true }}
      />
    </DrawerForm>
  );
};

export default SiteDrawer;
