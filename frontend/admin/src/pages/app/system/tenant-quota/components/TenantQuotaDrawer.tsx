import { useEffect, useRef, useState } from 'react';
import type { ProFormInstance } from '@ant-design/pro-components';
import {
  ProFormDigit,
  ProFormSelect,
  DrawerForm,
} from '@ant-design/pro-components';
import { App } from 'antd';
import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import type { identityservicev1_TenantQuota as TenantQuota } from '@/api/generated/admin/service/v1';
import {
  useCreateTenantQuota,
  useUpdateTenantQuota,
} from '@/api/hooks/tenant-quota';
import { useFetchTenantOptions } from '@/api/hooks/tenant';

// 从未知错误对象取本地化 message（request-client 拦截器已把 reason 译文写进 message）
const errMessage = (error: unknown, fallback: string): string => {
  if (error && typeof error === 'object' && 'message' in error) {
    const m = (error as { message?: unknown }).message;
    if (typeof m === 'string' && m) return m;
  }
  return fallback;
};

const quotaTypeOptions = (t: (k: string) => string) => [
  { value: 'USER_LIMIT', label: t('quotaTypeMap.USER_LIMIT') },
  { value: 'STORAGE', label: t('quotaTypeMap.STORAGE') },
  { value: 'API_CALL', label: t('quotaTypeMap.API_CALL') },
  { value: 'AI_TOKENS', label: t('quotaTypeMap.AI_TOKENS') },
];

interface TenantQuotaDrawerProps {
  open: boolean;
  mode: 'create' | 'edit';
  data?: TenantQuota;
  onClose: () => void;
  onSuccess: () => void;
}

/**
 * 配额创建/编辑抽屉。tenant_id/quota_type 创建后不可改（改口径=换行重建）。
 */
export default function TenantQuotaDrawer({
  open,
  mode,
  data,
  onClose,
  onSuccess,
}: TenantQuotaDrawerProps) {
  const { t } = useTranslation('tenant-quota');
  const { message } = App.useApp();
  const queryClient = useQueryClient();
  const formRef = useRef<ProFormInstance>(null);

  const [confirmLoading, setConfirmLoading] = useState(false);

  const createMutation = useCreateTenantQuota();
  const updateMutation = useUpdateTenantQuota();

  // 租户下拉（全量拉取，量有限）
  const { data: tenantOptions } = useFetchTenantOptions();

  useEffect(() => {
    if (!open) return;
    setTimeout(() => {
      if (mode === 'edit' && data) {
        formRef.current?.setFieldsValue({
          tenantId: data.tenantId,
          quotaType: data.quotaType,
          quotaValue: Number(data.quotaValue ?? 0),
        });
      }
    }, 0);
  }, [open, mode, data]);

  const resetAndClose = () => {
    formRef.current?.resetFields();
    onClose();
  };

  const handleSubmit = async (values: Record<string, any>) => {
    setConfirmLoading(true);
    try {
      if (mode === 'create') {
        await createMutation.mutateAsync({ values });
      } else {
        // 只提交配额值；tenantId/quotaType 不可变
        await updateMutation.mutateAsync({
          id: data!.id!,
          values: { quotaValue: values.quotaValue },
        });
      }
      message.success(t(mode === 'create' ? 'createSuccess' : 'updateSuccess'));
      queryClient.invalidateQueries({ queryKey: ['listTenantQuotas'] });
      queryClient.invalidateQueries({ queryKey: ['myTenantUsage'] });
      resetAndClose();
      onSuccess();
      return true;
    } catch (error: any) {
      // 原始错误必须留在控制台：用户可见的那句翻译不包含服务端的原因
      console.error('save tenant quota failed', error);
      message.error(errMessage(error, t('fetchFailed')));
      return false;
    } finally {
      setConfirmLoading(false);
    }
  };

  return (
    <DrawerForm
      title={mode === 'create' ? t('create') : t('edit')}
      open={open}
      formRef={formRef}
      drawerProps={{
        onClose: () => {
          formRef.current?.resetFields();
          onClose();
        },
        destroyOnClose: true,
        // antd6 Drawer width 已废弃：size 直接收数字
        size: 480,
      }}
      onFinish={handleSubmit}
      submitter={{ submitButtonProps: { loading: confirmLoading } }}
    >
      <ProFormSelect
        name="tenantId"
        label={t('tenantId')}
        placeholder={t('tenantIdPlaceholder')}
        rules={[{ required: true, message: t('requiredTenantId') }]}
        options={tenantOptions}
        disabled={mode === 'edit'}
        fieldProps={{ showSearch: true, optionFilterProp: 'label' }}
      />
      <ProFormSelect
        name="quotaType"
        label={t('quotaType')}
        placeholder={t('quotaTypePlaceholder')}
        rules={[{ required: true, message: t('requiredQuotaType') }]}
        options={quotaTypeOptions((k) => t(k))}
        disabled={mode === 'edit'}
      />
      <ProFormDigit
        name="quotaValue"
        label={t('quotaValue')}
        placeholder={t('quotaValuePlaceholder')}
        rules={[{ required: true, message: t('requiredQuotaValue') }]}
        min={0}
        fieldProps={{ precision: 0, style: { width: '100%' } }}
        extra={t('quotaValueHint')}
      />
    </DrawerForm>
  );
}
