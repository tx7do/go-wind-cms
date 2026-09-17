import { useEffect, useRef, useState } from 'react';
import { Button, Collapse, Drawer, Input, InputNumber, Select, Switch } from 'antd';
import { useTranslation } from 'react-i18next';

import {
  useCreateContentModel,
  useUpdateContentModel,
  fieldTypeOptions,
  RELATION_ENTITY_TYPES,
} from '@/api/hooks/content-model';

interface FieldPanel {
  clientKey: string;
  name: string;
  type: string;
  label: string;
  description: string;
  placeholder: string;
  isRequired: boolean;
  validationRegex: string;
  sortOrder: number;
  hasRelationConfig: boolean;
  relationTargetEntityType: string;
  relationAllowCrossTenant: boolean;
  relationFilterCategoryId: string;
  optionsJson: string;
}

interface ModelFormData {
  id: number;
  name: string;
  code: string;
  description: string;
  sortOrder: number;
  fields: FieldPanel[];
}

interface ContentModelDrawerProps {
  open: boolean;
  data?: any;
  onClose: () => void;
  onSuccess: () => void;
}

// options 桥接：FieldDefinition.options 是后端 map<string,string>，
// 前端面板用 JSON 文本编辑（与 site-setting 的 optionsJson 同模式）。
// 非对象/解析失败一律回退空对象，避免脏数据落库。
function parseOptionsJson(json: string | undefined): Record<string, string> {
  try {
    const parsed = json ? (JSON.parse(json) as unknown) : null;
    if (parsed && typeof parsed === 'object' && !Array.isArray(parsed)) {
      const result: Record<string, string> = {};
      for (const [k, v] of Object.entries(parsed as Record<string, unknown>)) {
        if (typeof v === 'string') result[k] = v;
      }
      return result;
    }
  } catch {
    // 忽略非法 JSON
  }
  return {};
}

let fieldCounter = 0;

function makeEmptyField(): FieldPanel {
  fieldCounter += 1;
  return {
    clientKey: `field-${fieldCounter}`,
    name: '',
    type: 'FIELD_TYPE_TEXT',
    label: '',
    description: '',
    placeholder: '',
    isRequired: false,
    validationRegex: '',
    sortOrder: 0,
    hasRelationConfig: false,
    relationTargetEntityType: '',
    relationAllowCrossTenant: false,
    relationFilterCategoryId: '',
    optionsJson: '{}',
  };
}

/**
 * 内容模型编辑抽屉：基础信息 + 字段定义编辑器（Collapse 面板，字段可增删）
 */
const ContentModelDrawer: React.FC<ContentModelDrawerProps> = ({
  open,
  data,
  onClose,
  onSuccess,
}) => {
  const { t } = useTranslation('content-model');

  const [saving, setSaving] = useState(false);
  const [formData, setFormData] = useState<ModelFormData>({
    id: 0,
    name: '',
    code: '',
    description: '',
    sortOrder: 0,
    fields: [],
  });
  const [activeKeys, setActiveKeys] = useState<string[]>([]);

  const setField = (clientKey: string, patch: Partial<FieldPanel>) => {
    setFormData((prev) => ({
      ...prev,
      fields: prev.fields.map((f) =>
        f.clientKey === clientKey ? { ...f, ...patch } : f,
      ),
    }));
  };

  const addField = () => {
    const f = makeEmptyField();
    setFormData((prev) => ({ ...prev, fields: [...prev.fields, f] }));
    setActiveKeys((prev) => [...prev, f.clientKey]);
  };

  const removeField = (clientKey: string) => {
    setFormData((prev) => ({
      ...prev,
      fields: prev.fields.filter((f) => f.clientKey !== clientKey),
    }));
    setActiveKeys((prev) => prev.filter((k) => k !== clientKey));
  };

  useEffect(() => {
    if (!open) return;
    if (data?.id) {
      setFormData({
        id: data.id ?? 0,
        name: data.name ?? '',
        code: data.code ?? '',
        description: data.description ?? '',
        sortOrder: data.sortOrder ?? 0,
        fields: (data.fields ?? []).map((f: any) => {
          fieldCounter += 1;
          return {
            clientKey: `field-${fieldCounter}`,
            name: f.name ?? '',
            type: f.type ?? 'FIELD_TYPE_TEXT',
            label: f.label ?? '',
            description: f.description ?? '',
            placeholder: f.placeholder ?? '',
            isRequired: f.isRequired ?? false,
            validationRegex: f.validationRegex ?? '',
            sortOrder: f.sortOrder ?? 0,
            hasRelationConfig: !!f.relationConfig,
            relationTargetEntityType: f.relationConfig?.targetEntityType ?? '',
            relationAllowCrossTenant: f.relationConfig?.allowCrossTenant ?? false,
            relationFilterCategoryId: f.relationConfig?.filterCategoryId ?? '',
            optionsJson: f.options ? JSON.stringify(f.options) : '{}',
          };
        }),
      });
    } else {
      setFormData({ id: 0, name: '', code: '', description: '', sortOrder: 0, fields: [] });
    }
  }, [open, data]);

  const handleSubmit = async () => {
    setSaving(true);
    try {
      const fieldsPayload = formData.fields.map((f) => {
        const field: Record<string, any> = {
          name: f.name,
          type: f.type,
          label: f.label,
          description: f.description,
          placeholder: f.placeholder,
          isRequired: f.isRequired,
          validationRegex: f.validationRegex,
          sortOrder: f.sortOrder,
        };
        if (f.type === 'FIELD_TYPE_SELECT') {
          field.options = parseOptionsJson(f.optionsJson);
        }
        if (f.type === 'FIELD_TYPE_RELATION' && f.hasRelationConfig) {
          field.relationConfig = {
            targetEntityType: f.relationTargetEntityType,
            allowCrossTenant: f.relationAllowCrossTenant,
            filterCategoryId: f.relationFilterCategoryId,
          };
        }
        return field;
      });

      const payload: Record<string, any> = {
        name: formData.name,
        code: formData.code,
        description: formData.description,
        sortOrder: formData.sortOrder,
        fields: fieldsPayload,
      };

      if (formData.id) {
        await updateMutationRef.current.mutateAsync({ id: formData.id, values: payload });
      } else {
        await createMutationRef.current.mutateAsync(payload);
      }
      onSuccess();
      onClose();
    } catch {
      // 错误消息由 hooks 的 onError 统一提示
    } finally {
      setSaving(false);
    }
  };

  const createMutation = useCreateContentModel();
  const updateMutation = useUpdateContentModel();
  const createMutationRef = useRef(createMutation);
  const updateMutationRef = useRef(updateMutation);
  createMutationRef.current = createMutation;
  updateMutationRef.current = updateMutation;

  const updateBase = (patch: Partial<ModelFormData>) => {
    setFormData((prev) => ({ ...prev, ...patch }));
  };

  return (
    <Drawer
      title={formData.id ? t('editTitle') : t('createTitle')}
      width="60%"
      open={open}
      onClose={onClose}
      destroyOnHidden
      footer={
        <div style={{ display: 'flex', width: '100%', justifyContent: 'flex-end' }}>
          <Button
            type="primary"
            loading={saving}
            onClick={handleSubmit}
          >
            {t('common:button.submit')}
          </Button>
        </div>
      }
    >
      <div style={{ display: 'flex', flexDirection: 'column', gap: 16, padding: 4 }}>
        <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
          <label style={{ fontSize: 14 }}>{t('name')}</label>
          <Input
            value={formData.name}
            onChange={(e) => updateBase({ name: e.target.value })}
          />
        </div>
        <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
          <label style={{ fontSize: 14 }}>{t('code')}</label>
          <Input
            value={formData.code}
            onChange={(e) => updateBase({ code: e.target.value })}
          />
        </div>
        <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
          <label style={{ fontSize: 14 }}>{t('description')}</label>
          <Input.TextArea
            rows={2}
            value={formData.description}
            onChange={(e) => updateBase({ description: e.target.value })}
          />
        </div>

        <hr style={{ border: 'none', borderTop: '1px solid var(--ant-color-border)' }} />

        <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <span style={{ fontSize: 15, fontWeight: 500 }}>{t('fieldsTitle')}</span>
          <Button size="small" onClick={addField}>
            {t('addField')}
          </Button>
        </div>

        {formData.fields.length > 0 ? (
          <Collapse
            activeKey={activeKeys}
            onChange={(keys) => setActiveKeys(Array.isArray(keys) ? keys : [keys])}
            items={formData.fields.map((field) => ({
              key: field.clientKey,
              label: field.name || t('unnamedField'),
              children: (
                <div style={{ display: 'flex', flexDirection: 'column', gap: 12, padding: 8 }}>
                  <div style={{ display: 'flex', justifyContent: 'flex-end' }}>
                    <Button
                      danger
                      size="small"
                      onClick={() => removeField(field.clientKey)}
                    >
                      {t('removeField')}
                    </Button>
                  </div>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                    <label style={{ fontSize: 14 }}>{t('fieldName')}</label>
                    <Input
                      value={field.name}
                      onChange={(e) => setField(field.clientKey, { name: e.target.value })}
                    />
                  </div>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                    <label style={{ fontSize: 14 }}>{t('fieldType')}</label>
                    <Select
                      value={field.type}
                      options={fieldTypeOptions(t)}
                      onChange={(v) => setField(field.clientKey, { type: v })}
                    />
                  </div>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                    <label style={{ fontSize: 14 }}>{t('fieldLabel')}</label>
                    <Input
                      value={field.label}
                      onChange={(e) => setField(field.clientKey, { label: e.target.value })}
                    />
                  </div>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                    <label style={{ fontSize: 14 }}>{t('fieldDescription')}</label>
                    <Input.TextArea
                      rows={2}
                      value={field.description}
                      onChange={(e) => setField(field.clientKey, { description: e.target.value })}
                    />
                  </div>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                    <label style={{ fontSize: 14 }}>{t('fieldPlaceholder')}</label>
                    <Input
                      value={field.placeholder}
                      onChange={(e) => setField(field.clientKey, { placeholder: e.target.value })}
                    />
                  </div>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                    <label style={{ fontSize: 14 }}>{t('sortOrder')}</label>
                    <InputNumber
                      min={0}
                      style={{ width: '100%' }}
                      value={field.sortOrder}
                      onChange={(v) => setField(field.clientKey, { sortOrder: v ?? 0 })}
                    />
                  </div>
                  <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                    <Switch
                      checked={field.isRequired}
                      onChange={(v) => setField(field.clientKey, { isRequired: v })}
                    />
                    <span style={{ fontSize: 14 }}>{t('fieldIsRequired')}</span>
                  </div>
                  <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                    <label style={{ fontSize: 14 }}>{t('fieldValidationRegex')}</label>
                    <Input
                      value={field.validationRegex}
                      onChange={(e) => setField(field.clientKey, { validationRegex: e.target.value })}
                    />
                  </div>

                  {field.type === 'FIELD_TYPE_SELECT' && (
                    <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
                      <label style={{ fontSize: 14 }}>{t('fieldOptions')}</label>
                      <Input.TextArea
                        rows={4}
                        value={field.optionsJson}
                        onChange={(e) => setField(field.clientKey, { optionsJson: e.target.value })}
                      />
                    </div>
                  )}

                  {field.type === 'FIELD_TYPE_RELATION' && (
                    <div
                      style={{
                        display: 'flex',
                        flexDirection: 'column',
                        gap: 8,
                        borderTop: '1px solid var(--ant-color-border)',
                        paddingTop: 12,
                        marginTop: 8,
                      }}
                    >
                      <label style={{ fontSize: 14 }}>{t('relationTargetEntityType')}</label>
                      <Select
                        value={field.relationTargetEntityType || undefined}
                        options={RELATION_ENTITY_TYPES}
                        onChange={(v) =>
                          setField(field.clientKey, {
                            hasRelationConfig: true,
                            relationTargetEntityType: v ?? '',
                          })
                        }
                      />
                      <div style={{ display: 'flex', alignItems: 'center', gap: 8, marginTop: 8 }}>
                        <Switch
                          checked={field.relationAllowCrossTenant}
                          onChange={(v) =>
                            setField(field.clientKey, {
                              hasRelationConfig: true,
                              relationAllowCrossTenant: v,
                            })
                          }
                        />
                        <span style={{ fontSize: 14 }}>{t('relationAllowCrossTenant')}</span>
                      </div>
                      <div style={{ display: 'flex', flexDirection: 'column', gap: 8, marginTop: 8 }}>
                        <label style={{ fontSize: 14 }}>{t('relationFilterCategoryId')}</label>
                        <Input
                          value={field.relationFilterCategoryId}
                          onChange={(e) =>
                            setField(field.clientKey, {
                              hasRelationConfig: true,
                              relationFilterCategoryId: e.target.value,
                            })
                          }
                        />
                      </div>
                    </div>
                  )}
                </div>
              ),
            }))}
          />
        ) : (
          <div style={{ textAlign: 'center', color: 'var(--ant-color-text-quaternary)', padding: '32px 0' }}>
            {t('emptyFields')}
          </div>
        )}
      </div>
    </Drawer>
  );
};

export default ContentModelDrawer;
