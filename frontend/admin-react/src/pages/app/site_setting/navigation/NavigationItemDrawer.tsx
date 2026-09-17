import { useEffect, useRef, useState } from 'react';
import type { siteservicev1_NavigationItem } from '@/api/generated/admin/service/v1';
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
  useCreateNavigationItem,
  useUpdateNavigationItem,
} from '@/api/hooks/navigation-item';
import { linkTypeOptions } from './constants';

interface NavigationItemDrawerProps {
  open: boolean;
  mode: 'create' | 'edit';
  data?: any;
  /** 父级导航归属：Create 时由列表上下文注入（当前选中导航），编辑时以行数据为准 */
  navigationId: number | null;
  onClose: () => void;
  onSuccess: () => void;
}

/**
 * 导航项编辑/创建抽屉组件
 */
const NavigationItemDrawer: React.FC<NavigationItemDrawerProps> = ({
  open,
  mode,
  data,
  navigationId,
  onClose,
  onSuccess,
}) => {
  const { t } = useTranslation('navigation-item');
  const formRef = useRef<ProFormInstance>(null);
  const queryClient = useQueryClient();
  const { message } = App.useApp();

  const [confirmLoading, setConfirmLoading] = useState(false);

  // 编辑模式下设置表单值（destroyOnHidden 时需延迟赋值）
  useEffect(() => {
    if (open && mode === 'edit' && data) {
      setTimeout(() => {
        formRef.current?.setFieldsValue({
          title: data.title || '',
          description: data.description || '',
          icon: data.icon || '',
          url: data.url || '',
          linkType: data.linkType || 'LINK_TYPE_CUSTOM',
          isOpenNewTab: data.isOpenNewTab ?? false,
          isInvalid: data.isInvalid ?? false,
        });
      }, 0);
    }
  }, [open, mode, data]);

  // 创建 mutation
  const createMutation = useCreateNavigationItem({
    onSuccess: () => {
      message.success(t('createSuccess'));
      queryClient.invalidateQueries({ queryKey: ['listNavigationItems'] });
      onSuccess();
      onClose();
    },
    onError: (error: Error) => {
      message.error(error.message || t('createFailed'));
    },
  });

  // 更新 mutation
  const updateMutation = useUpdateNavigationItem({
    onSuccess: () => {
      message.success(t('updateSuccess'));
      queryClient.invalidateQueries({ queryKey: ['listNavigationItems'] });
      onSuccess();
      onClose();
    },
    onError: (error: Error) => {
      message.error(error.message || t('updateFailed'));
    },
  });

  // 提交表单
  const handleSubmit = async (values: Record<string, any>) => {
    // 父级导航归属：Create 时由列表上下文注入（当前选中导航），
    // 不在表单渲染；缺则拒绝创建以避免产生无归属的孤儿导航项。
    const parentNavigationId =
      mode === 'create' ? navigationId : (data?.navigationId ?? navigationId);
    if (mode === 'create' && !parentNavigationId) {
      message.error(t('navigationIdRequired'));
      return false;
    }

    try {
      setConfirmLoading(true);
      const payload: Record<string, any> = {
        title: values.title,
        description: values.description,
        icon: values.icon,
        url: values.url,
        linkType: values.linkType,
        isOpenNewTab: values.isOpenNewTab ?? false,
        isInvalid: values.isInvalid ?? false,
        navigationId: parentNavigationId,
      };

      if (mode === 'edit' && data?.id) {
        await updateMutation.mutateAsync({ id: data.id, values: payload as unknown as siteservicev1_NavigationItem });
      } else {
        await createMutation.mutateAsync({ data: payload as unknown as siteservicev1_NavigationItem });
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
        linkType: 'LINK_TYPE_CUSTOM',
        isOpenNewTab: false,
        isInvalid: false,
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
        name="title"
        label={t('title')}
        placeholder={t('titlePlaceholder')}
        rules={[{ required: true, message: t('requiredTitle') }]}
        fieldProps={{ allowClear: true }}
      />

      <ProFormText
        name="description"
        label={t('description')}
        placeholder={t('descriptionPlaceholder')}
        fieldProps={{ allowClear: true }}
      />

      <ProFormText
        name="icon"
        label={t('icon')}
        placeholder={t('iconPlaceholder')}
        rules={[{ required: true, message: t('requiredIcon') }]}
        fieldProps={{ allowClear: true }}
      />

      <ProFormText
        name="url"
        label={t('url')}
        placeholder={t('urlPlaceholder')}
        rules={[{ required: true, message: t('requiredUrl') }]}
        fieldProps={{ allowClear: true }}
      />

      <ProFormSelect
        name="linkType"
        label={t('linkType')}
        placeholder={t('linkTypePlaceholder')}
        rules={[{ required: true, message: t('requiredLinkType') }]}
        options={linkTypeOptions(t)}
        fieldProps={{ allowClear: true }}
      />

      <ProFormSwitch name="isOpenNewTab" label={t('isOpenNewTab')} />

      <ProFormSwitch name="isInvalid" label={t('isInvalid')} />
    </DrawerForm>
  );
};

export default NavigationItemDrawer;
