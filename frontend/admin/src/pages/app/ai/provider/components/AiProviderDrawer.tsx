import { useEffect, useRef } from 'react';
import type { ProFormInstance } from '@ant-design/pro-components';
import {
  DrawerForm,
  ProFormDependency,
  ProFormDigit,
  ProFormRadio,
  ProFormSelect,
  ProFormText,
  ProFormTextArea,
} from '@ant-design/pro-components';
import { App } from 'antd';
import { useTranslation } from 'react-i18next';
import { useQueryClient } from '@tanstack/react-query';
import type { aiservicev1_AiProvider as AiProvider } from '@/api/generated/admin/service/v1';
import { useCreateAiProvider, useUpdateAiProvider } from '@/api/hooks/ai-provider';
import { useDrawerCloseGuard } from '@/components/common/use-drawer-close-guard';

interface AiProviderDrawerProps {
  open: boolean;
  mode: 'create' | 'edit';
  data?: AiProvider;
  onClose: () => void;
  onSuccess: () => void;
}

const MODEL_TYPE_OPTIONS = (t: (k: string) => string) => [
  { label: t('modelTypeCloud'), value: 'CLOUD' },
  { label: t('modelTypeLocal'), value: 'LOCAL' },
];

export default function AiProviderDrawer({ open, mode, data, onClose, onSuccess }: AiProviderDrawerProps) {
  const { t } = useTranslation('aiProvider');
  const { message } = App.useApp();
  const queryClient = useQueryClient();
  const formRef = useRef<ProFormInstance>(null);

  // 关闭守卫：遮罩/ESC/X/取消统一，脏数据先确认（docs/form_drawer_close_guard.md）
  const { guardClose, handleValuesChange } = useDrawerCloseGuard(open, onClose);

  const createMutation = useCreateAiProvider();
  const updateMutation = useUpdateAiProvider();

  // 编辑回填：等抽屉渲染完成后再注入（initialValues 不会随编辑记录更新）
  useEffect(() => {
    if (!open) return;
    const timer = setTimeout(() => {
      if (mode === 'edit' && data) {
        formRef.current?.setFieldsValue({
          name: data.name,
          modelType: data.modelType ?? 'CLOUD',
          modelName: data.modelName,
          baseUrl: data.baseUrl,
          organization: data.organization,
          localHost: data.localHost,
          localPort: data.localPort,
          timeoutSeconds: data.timeoutSeconds,
          systemPrompt: data.systemPrompt,
          isDefault: data.isDefault,
          remark: data.remark,
          isEnabled: data.isEnabled ?? true,
          apiKey: undefined, // 密钥永不回显；留空表示不改
        });
      } else {
        formRef.current?.setFieldsValue({
          modelType: 'CLOUD',
          isDefault: false,
          isEnabled: true,
        });
      }
    }, 0);
    return () => clearTimeout(timer);
  }, [open, mode, data]);

  const handleSubmit = async (values: Record<string, any>) => {
    try {
      if (mode === 'create') {
        await createMutation.mutateAsync({ data: values as any });
        message.success(t('createSuccess'));
      } else {
        // api_key 未填则不进 values（useUpdateAiProvider 会按传入键构建 updateMask），
        // 避免空值经 FieldMask 被当成"显式清空"把已存的 Key 抹掉
        const payload: Record<string, any> = { ...values };
        if (!payload.apiKey) {
          delete payload.apiKey;
        }
        await updateMutation.mutateAsync({ id: data!.id!, values: payload });
        message.success(t('updateSuccess'));
      }
      queryClient.invalidateQueries({ queryKey: ['listAiProviders'] });
      onSuccess();
      onClose();
      return true;
    } catch (error) {
      console.error('save ai provider failed:', error);
      message.error((error as Error).message);
      return false;
    }
  };

  return (
    <DrawerForm
      formRef={formRef}
      title={mode === 'create' ? t('common:modal.create', { moduleName: t('moduleName') }) : t('common:modal.update', { moduleName: t('moduleName') })}
      open={open}
      onValuesChange={handleValuesChange}
      onFinish={handleSubmit}
      submitter={{ submitButtonProps: { loading: createMutation.isPending || updateMutation.isPending }, resetButtonProps: { onClick: guardClose } }}
      drawerProps={{ destroyOnHidden: true, onClose: guardClose, size: 600 }}
      layout="horizontal"
      labelCol={{ flex: '0 0 120px' }}
      wrapperCol={{ flex: 'auto' }}
    >
      <ProFormText
        name="name"
        label={t('name')}
        rules={[{ required: true, message: t('requiredName') }]}
        placeholder={t('namePlaceholder')}
      />
      <ProFormRadio.Group
        name="modelType"
        label={t('modelType')}
        options={MODEL_TYPE_OPTIONS(t)}
        rules={[{ required: true }]}
      />
      <ProFormText
        name="modelName"
        label={t('modelName')}
        rules={[{ required: true, message: t('requiredModelName') }]}
        placeholder={t('modelNamePlaceholder')}
      />
      <ProFormDependency name={['modelType']}>
        {({ modelType }) =>
          modelType === 'LOCAL' ? (
            <>
              <ProFormText name="localHost" label={t('localHost')} placeholder={t('localHostPlaceholder')} />
              <ProFormDigit name="localPort" label={t('localPort')} min={1} max={65535} fieldProps={{ precision: 0 }} placeholder={t('localPortPlaceholder')} />
            </>
          ) : (
            <>
              <ProFormText
                name="baseUrl"
                label={t('baseUrl')}
                rules={[{ required: true, message: t('requiredBaseUrl') }]}
                placeholder={t('baseUrlPlaceholder')}
              />
              <ProFormText name="organization" label={t('organization')} />
              <ProFormText.Password
                name="apiKey"
                label={`${t('apiKey')}${mode === 'edit' ? `（${t('apiKeyHint')}）` : ''}`}
                placeholder={mode === 'edit' ? t('apiKeyKeepExisting') : t('apiKeyPlaceholder')}
                extra={mode === 'edit' ? data?.apiKeyHint : undefined}
              />
            </>
          )
        }
      </ProFormDependency>
      <ProFormDigit
        name="timeoutSeconds"
        label={t('timeoutSeconds')}
        min={1}
        max={600}
        fieldProps={{ precision: 0 }}
      />
      <ProFormTextArea
        name="systemPrompt"
        label={t('systemPrompt')}
        placeholder={t('systemPromptPlaceholder')}
        fieldProps={{ rows: 3 }}
      />
      <ProFormSelect
        name="isDefault"
        label={t('isDefault')}
        options={[
          { label: 'Yes', value: true },
          { label: 'No', value: false },
        ]}
        tooltip={t('isDefaultTooltip')}
      />
      <ProFormSelect
        name="isEnabled"
        label={t('isEnabled')}
        options={[
          { label: t('statusMap.ON'), value: true },
          { label: t('statusMap.OFF'), value: false },
        ]}
        rules={[{ required: true }]}
      />
      <ProFormTextArea
        name="remark"
        label={t('remark')}
        fieldProps={{ rows: 2 }}
      />
    </DrawerForm>
  );
}
