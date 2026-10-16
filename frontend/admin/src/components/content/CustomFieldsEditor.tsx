import { useCallback, useEffect, useState } from 'react';
import { Button, Input, InputNumber, Select, Space, Tag } from 'antd';
import {
  DeleteOutlined,
  DownOutlined,
  PlusOutlined,
  UpOutlined,
} from '@ant-design/icons';
import { useTranslation } from 'react-i18next';
import type {
  contentservicev1_ComponentDefinition,
  contentservicev1_FieldDefinition,
} from '@/api/generated/admin/service/v1';
import {
  fetchComponentFieldDefinitions,
  fetchListComponentDefinitions,
} from '@/api/hooks/component-definition';

/**
 * 自定义字段编辑器（内容建模读侧写侧共用）：
 * 叶子字段简单输入；COMPONENT 单实例/可重复与 DYNAMIC_ZONE 渲染实例卡片
 * （增删/上下移），值按约定编码为 JSON 字符串（单实例=对象、数组=可重复/动态区，
 * 动态区元素带 _component 判别键）。既有值为非法 JSON 时原样保留在文本框，
 * 不阻断保存（后端校验兜底）。
 */

type CustomFields = Record<string, string>;
type FieldDef = Partial<contentservicev1_FieldDefinition>;
type ComponentSchema = {
  definition: Partial<contentservicev1_ComponentDefinition>;
  fields: FieldDef[];
};

interface InstanceData {
  _component?: string;
  [key: string]: unknown;
}

function parseInstances(value: string | undefined, expectArray: boolean): InstanceData[] | null {
  if (!value || !value.trim()) return [];
  try {
    const parsed = JSON.parse(value) as unknown;
    if (expectArray) {
      if (!Array.isArray(parsed)) return null;
      return parsed.filter(
        (x): x is InstanceData => Boolean(x) && typeof x === 'object' && !Array.isArray(x),
      );
    }
    if (!parsed || typeof parsed !== 'object' || Array.isArray(parsed)) return null;
    return [parsed as InstanceData];
  } catch {
    return null;
  }
}

/** 按模型字段集收集引用到的组件编码 */
function collectCodes(fields: FieldDef[]): string[] {
  const codes: string[] = [];
  for (const f of fields) {
    const cc = f.componentConfig;
    if (!cc) continue;
    if (cc.componentCode && !codes.includes(cc.componentCode)) codes.push(cc.componentCode);
    for (const c of cc.allowedComponentCodes ?? []) {
      if (!codes.includes(c)) codes.push(c);
    }
  }
  return codes;
}

/** 组件编码 → 字段 schema（一次拉取全部引用到的组件） */
function useComponentSchemas(fields: FieldDef[]) {
  const [schemas, setSchemas] = useState<Record<string, ComponentSchema>>({});
  const signature = collectCodes(fields).join(',');

  useEffect(() => {
    const codes = signature ? signature.split(',') : [];
    let cancelled = false;
    (async () => {
      if (codes.length === 0) {
        if (!cancelled) setSchemas({});
        return;
      }
      try {
        const list = await fetchListComponentDefinitions();
        const idByCode = new Map<string, number>();
        for (const c of list.items ?? []) {
          if (c.id && c.code) idByCode.set(c.code, c.id);
        }
        const next: Record<string, ComponentSchema> = {};
        for (const code of codes) {
          const id = idByCode.get(code);
          if (!id) continue;
          const resp = await fetchComponentFieldDefinitions(id);
          const def = (list.items ?? []).find((c) => c.code === code);
          next[code] = {
            definition: def ?? { code, name: code },
            fields: resp.items ?? [],
          };
        }
        if (!cancelled) setSchemas(next);
      } catch (error) {
        console.error('Load component schemas failed:', error);
      }
    })();
    return () => {
      cancelled = true;
    };
  }, [signature]);

  return schemas;
}

/** 实例卡片头（动态区显示组件名 + 移动/删除） */
function InstanceActions({
  onUp,
  onDown,
  onRemove,
  upDisabled,
  downDisabled,
  removeDisabled,
}: {
  onUp?: () => void;
  onDown?: () => void;
  onRemove?: () => void;
  upDisabled?: boolean;
  downDisabled?: boolean;
  removeDisabled?: boolean;
}) {
  return (
    <Space size={0}>
      {onUp && (
        <Button type="text" size="small" icon={<UpOutlined />} disabled={upDisabled} onClick={onUp} />
      )}
      {onDown && (
        <Button type="text" size="small" icon={<DownOutlined />} disabled={downDisabled} onClick={onDown} />
      )}
      {onRemove && (
        <Button type="text" size="small" danger icon={<DeleteOutlined />} disabled={removeDisabled} onClick={onRemove} />
      )}
    </Space>
  );
}

/** 组件实例内的叶子字段输入 */
function LeafInput({
  field,
  value,
  onChange,
}: {
  field: FieldDef;
  value: unknown;
  onChange: (next: string) => void;
}) {
  const { t } = useTranslation('custom-fields');
  const text = typeof value === 'string' ? value : '';
  if (field.type === 'FIELD_TYPE_RICHTEXT') {
    return (
      <Input.TextArea
        rows={2}
        size="small"
        value={text}
        placeholder={field.label || field.name}
        onChange={(e) => onChange(e.target.value)}
      />
    );
  }
  if (field.type === 'FIELD_TYPE_NUMBER') {
    return (
      <InputNumber
        size="small"
        style={{ width: '100%' }}
        value={text === '' ? undefined : Number(text)}
        placeholder={field.label || field.name}
        onChange={(v) => onChange(v == null ? '' : String(v))}
      />
    );
  }
  return (
    <Input
      size="small"
      value={text}
      placeholder={
        field.type === 'FIELD_TYPE_IMAGE' || field.type === 'FIELD_TYPE_FILE'
          ? t('mediaRefPlaceholder')
          : field.label || field.name
      }
      onChange={(e) => onChange(e.target.value)}
    />
  );
}

/** 组件实例卡片的字段区 */
function InstanceFields({
  schemaFields,
  instance,
  qualName,
  onPatch,
}: {
  schemaFields: FieldDef[];
  instance: InstanceData;
  qualName: string;
  onPatch: (key: string, next: string) => void;
}) {
  return (
    <>
      {schemaFields.map((f) => {
        const key = f.name ?? '';
        return (
          <div key={String(f.id ?? key)} style={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
            <span style={{ fontSize: 12, color: 'var(--ant-color-text-secondary)' }}>
              {f.label || key}
              {f.isRequired ? <Tag color="orange" style={{ marginLeft: 4 }}>*</Tag> : null}
            </span>
            <LeafInput
              field={f}
              value={instance[key]}
              onChange={(next) => onPatch(key, next)}
            />
            <span style={{ fontSize: 11, color: 'var(--ant-color-text-quaternary)' }}>{qualName}</span>
          </div>
        );
      })}
    </>
  );
}

export default function CustomFieldsEditor({
  fields,
  value,
  onChange,
}: {
  fields: FieldDef[];
  value?: CustomFields;
  onChange: (next: CustomFields) => void;
}) {
  const { t } = useTranslation('custom-fields');
  const schemas = useComponentSchemas(fields);

  const setValue = useCallback(
    (name: string, next: string) => {
      onChange({ ...(value ?? {}), [name]: next });
    },
    [onChange, value],
  );

  const renderComponentField = (field: FieldDef) => {
    const name = field.name ?? '';
    const cc = field.componentConfig;
    if (!cc) return null;
    const expectArray = !!cc.repeatable || field.type === 'FIELD_TYPE_DYNAMIC_ZONE';
    const instances = parseInstances(value?.[name], expectArray);

    // 非法 JSON：原样文本框，不阻断（保存由后端校验拦截）
    if (instances === null) {
      return (
        <Input
          status="warning"
          size="small"
          value={value?.[name] ?? ''}
          onChange={(e) => setValue(name, e.target.value)}
        />
      );
    }

    const patchInstance = (index: number, key: string, next: string) => {
      const updated = instances.map((inst, i) =>
        i === index ? { ...inst, [key]: next } : inst,
      );
      setValue(name, JSON.stringify(expectArray ? updated : updated[0] ?? {}));
    };

    const schemaOf = (code: string): FieldDef[] => schemas[code]?.fields ?? [];
    const nameOf = (code: string): string =>
      schemas[code]?.definition?.name ?? code;

    if (field.type === 'FIELD_TYPE_COMPONENT') {
      const code = cc.componentCode ?? '';
      const schemaFields = schemaOf(code);
      return (
        <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
          {instances.map((inst, i) => (
            <div
              key={i}
              style={{
                border: '1px solid var(--ant-color-border-secondary)',
                borderRadius: 8,
                padding: 8,
                display: 'flex',
                flexDirection: 'column',
                gap: 8,
              }}
            >
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <Tag color="purple">{nameOf(code)}</Tag>
                {cc.repeatable ? (
                  <InstanceActions
                    onUp={i > 0 ? () => {
                      const arr = [...instances];
                      [arr[i - 1], arr[i]] = [arr[i], arr[i - 1]];
                      setValue(name, JSON.stringify(arr));
                    } : undefined}
                    upDisabled={i === 0}
                    onDown={i < instances.length - 1 ? () => {
                      const arr = [...instances];
                      [arr[i + 1], arr[i]] = [arr[i], arr[i + 1]];
                      setValue(name, JSON.stringify(arr));
                    } : undefined}
                    downDisabled={i === instances.length - 1}
                    onRemove={() =>
                      setValue(name, JSON.stringify(instances.filter((_, j) => j !== i)))
                    }
                  />
                ) : null}
              </div>
              <InstanceFields
                schemaFields={schemaFields}
                instance={inst}
                qualName={`${name}[${i}]`}
                onPatch={(key, next) => patchInstance(i, key, next)}
              />
            </div>
          ))}
          {cc.repeatable ? (
            <Button
              block
              type="dashed"
              size="small"
              icon={<PlusOutlined />}
              onClick={() =>
                setValue(name, JSON.stringify([...instances, {}]))
              }
            >
              {t('addInstance', { name: nameOf(code) })}
            </Button>
          ) : null}
        </div>
      );
    }

    // DYNAMIC_ZONE
    const allowed = cc.allowedComponentCodes ?? [];
    return (
      <div style={{ display: 'flex', flexDirection: 'column', gap: 8 }}>
        {instances.map((inst, i) => {
          const code = String(inst._component ?? '');
          return (
            <div
              key={i}
              style={{
                border: '1px solid var(--ant-color-border-secondary)',
                borderRadius: 8,
                padding: 8,
                display: 'flex',
                flexDirection: 'column',
                gap: 8,
              }}
            >
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                <Tag color="geekblue">{nameOf(code)}</Tag>
                <InstanceActions
                  onUp={i > 0 ? () => {
                    const arr = [...instances];
                    [arr[i - 1], arr[i]] = [arr[i], arr[i - 1]];
                    setValue(name, JSON.stringify(arr));
                  } : undefined}
                  upDisabled={i === 0}
                  onDown={i < instances.length - 1 ? () => {
                    const arr = [...instances];
                    [arr[i + 1], arr[i]] = [arr[i], arr[i + 1]];
                    setValue(name, JSON.stringify(arr));
                  } : undefined}
                  downDisabled={i === instances.length - 1}
                  onRemove={() =>
                    setValue(name, JSON.stringify(instances.filter((_, j) => j !== i)))
                  }
                />
              </div>
              <InstanceFields
                schemaFields={schemaOf(code)}
                instance={inst}
                qualName={`${name}[${i}]`}
                onPatch={(key, next) => patchInstance(i, key, next)}
              />
            </div>
          );
        })}
        <Select
          size="small"
          value={undefined}
          placeholder={t('addBlock')}
          options={allowed.map((code) => ({ value: code, label: nameOf(code) }))}
          onChange={(code) => {
            setValue(name, JSON.stringify([...instances, { _component: code }]));
          }}
        />
      </div>
    );
  };

  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
      {fields.map((field) => {
        const name = field.name ?? '';
        const isComponent =
          field.type === 'FIELD_TYPE_COMPONENT' || field.type === 'FIELD_TYPE_DYNAMIC_ZONE';
        return (
          <div key={String(field.id ?? name)} style={{ display: 'flex', flexDirection: 'column', gap: 4 }}>
            <span style={{ fontSize: 13, fontWeight: 500 }}>{field.label || name}</span>
            {isComponent ? (
              renderComponentField(field)
            ) : field.type === 'FIELD_TYPE_RICHTEXT' ? (
              <Input.TextArea
                rows={2}
                size="small"
                value={value?.[name] ?? ''}
                onChange={(e) => setValue(name, e.target.value)}
              />
            ) : (
              <LeafInput
                field={field}
                value={value?.[name] ?? ''}
                onChange={(next) => setValue(name, next)}
              />
            )}
          </div>
        );
      })}
    </div>
  );
}
