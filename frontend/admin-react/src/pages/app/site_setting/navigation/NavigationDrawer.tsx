import { useEffect, useRef, useState } from 'react';
import type { siteservicev1_Navigation } from '@/api/generated/admin/service/v1';
import type { ProFormInstance } from '@ant-design/pro-components';
import {
  DrawerForm,
  ProFormText,
  ProFormSelect,
  ProFormSwitch,
} from '@ant-design/pro-components';
import { App } from 'antd';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';

import {
  useCreateNavigation,
  useUpdateNavigation,
} from '@/api/hooks/navigation';
import { fetchListLanguages } from '@/api/hooks/language';
import { PaginationQuery } from '@/core';
import { locationOptions } from './constants';

interface NavigationDrawerProps {
  open: boolean;
  mode: 'create' | 'edit';
  data?: any;
  onClose: () => void;
  onSuccess: () => void;
}

/**
 * 导航编辑/创建抽屉组件
 */
const NavigationDrawer: React.FC<NavigationDrawerProps> = ({
  open,
  mode,
  data,
  onClose,
  onSuccess,
}) => {
  const { t } = useTranslation('navigation');
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
          location: data.location,
          locale: data.locale || 'zh-CN',
          isActive: data.isActive ?? true,
        });
      }, 0);
    }
  }, [open, mode, data]);

  // 创建 mutation
  const createMutation = useCreateNavigation({
    onSuccess: () => {
      message.success(t('createSuccess'));
      queryClient.invalidateQueries({ queryKey: ['listNavigations'] });
      onSuccess();
      onClose();
    },
    onError: (error: Error) => {
      message.error(error.message || t('createFailed'));
    },
  });

  // 更新 mutation
  const updateMutation = useUpdateNavigation({
    onSuccess: () => {
      message.success(t('updateSuccess'));
      queryClient.invalidateQueries({ queryKey: ['listNavigations'] });
      onSuccess();
      onClose();
    },
    onError: (error: Error) => {
      message.error(error.message || t('updateFailed'));
    },
  });

  // 提交表单
  // 导航项（items）由独立的导航项子表管理 CRUD，此抽屉仅处理导航自身字段；
  // payload 不含 items，避免整体替换语义误清空子项。
  const handleSubmit = async (values: Record<string, any>) => {
    try {
      setConfirmLoading(true);
      const payload: Record<string, any> = {
        name: values.name,
        location: values.location,
        locale: values.locale,
        isActive: values.isActive ?? true,
      };

      if (mode === 'edit' && data?.id) {
        await updateMutation.mutateAsync({ id: data.id, values: payload as unknown as siteservicev1_Navigation });
      } else {
        await createMutation.mutateAsync({ data: payload as unknown as siteservicev1_Navigation });
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
        locale: 'zh-CN',
        isActive: true,
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

      <ProFormSelect
        name="location"
        label={t('location')}
        placeholder={t('locationPlaceholder')}
        rules={[{ required: true, message: t('requiredLocation') }]}
        options={locationOptions(t)}
        fieldProps={{ allowClear: true }}
      />

      <ProFormSelect
        name="locale"
        label={t('locale')}
        placeholder={t('localePlaceholder')}
        rules={[{ required: true, message: t('requiredLocale') }]}
        options={languageOptions}
        fieldProps={{ allowClear: true }}
      />

      <ProFormSwitch name="isActive" label={t('isActive')} />
    </DrawerForm>
  );
};

export default NavigationDrawer;
