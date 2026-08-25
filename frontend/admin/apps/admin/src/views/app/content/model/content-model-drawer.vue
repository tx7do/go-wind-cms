<script lang="ts" setup>
import { ref, reactive } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';

import { notification, Collapse, CollapsePanel, Button, Select, Input, InputNumber, Switch, Textarea } from 'ant-design-vue';

import {
  fieldDefinitionTypeList,
  relationEntityTypeList,
  useCreateContentModel,
  useUpdateContentModel,
} from '#/api';
import { $t } from '#/locales';

const { mutateAsync: createContentModel } = useCreateContentModel();
const { mutateAsync: updateContentModel } = useUpdateContentModel();

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

// options 桥接：FieldDefinition.options 是后端 map<string,string>，
// 前端面板用 JSON 文本编辑（与 site-setting 的 optionsJson 同模式）。
// 非对象/解析失败一律回退空对象，避免脏数据落库。
function serializeOptions(options: Record<string, string> | undefined): string {
  try {
    return options ? JSON.stringify(options) : '{}';
  } catch {
    return '{}';
  }
}

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

const formData = reactive<ModelFormData>({
  id: 0,
  name: '',
  code: '',
  description: '',
  sortOrder: 0,
  fields: [],
});

const activeKeys = ref<string[]>([]);

const [Drawer, drawerApi] = useVbenDrawer({
  onCancel() {
    drawerApi.close();
  },
  async onConfirm() {
    setLoading(true);

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
        await updateContentModel({ id: formData.id, values: payload });
        notification.success({
          message: $t('ui.notification.update_success'),
        });
      } else {
        await createContentModel(payload);
        notification.success({
          message: $t('ui.notification.create_success'),
        });
      }

      drawerApi.close();
    } catch {
      notification.error({
        message: $t('ui.notification.create_failed'),
      });
    } finally {
      setLoading(false);
    }
  },
  onOpenChange(isOpen: boolean) {
    if (isOpen) {
      const data = drawerApi.getData<{ row?: any }>();
      if (data?.row) {
        formData.id = data.row.id ?? 0;
        formData.name = data.row.name ?? '';
        formData.code = data.row.code ?? '';
        formData.description = data.row.description ?? '';
        formData.sortOrder = data.row.sortOrder ?? 0;
        formData.fields = (data.row.fields ?? []).map((f: any) => {
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
            optionsJson: serializeOptions(f.options),
          };
        });
      } else {
        formData.id = 0;
        formData.name = '';
        formData.code = '';
        formData.description = '';
        formData.sortOrder = 0;
        formData.fields = [];
      }
      activeKeys.value = formData.fields.map((f) => f.clientKey);
      setLoading(false);
    }
  },
});

function setLoading(loading: boolean) {
  drawerApi.setState({ confirmLoading: loading });
}

function addField() {
  const f = makeEmptyField();
  formData.fields.push(f);
  activeKeys.value = [...activeKeys.value, f.clientKey];
}

function removeField(clientKey: string) {
  const idx = formData.fields.findIndex((f) => f.clientKey === clientKey);
  if (idx >= 0) {
    formData.fields.splice(idx, 1);
    activeKeys.value = activeKeys.value.filter((k) => k !== clientKey);
  }
}
</script>

<template>
  <Drawer
    :title="
      formData.id
        ? $t('page.contentModel.title.edit')
        : $t('page.contentModel.title.create')
    "
    class="w-[60%]"
  >
    <div class="flex flex-col gap-4 p-4">
      <!-- 模型基础信息 -->
      <div class="flex flex-col gap-2">
        <label class="text-sm">{{ $t('page.contentModel.name') }}</label>
        <Input v-model:value="formData.name" />
      </div>
      <div class="flex flex-col gap-2">
        <label class="text-sm">{{ $t('page.contentModel.code') }}</label>
        <Input v-model:value="formData.code" />
      </div>
      <div class="flex flex-col gap-2">
        <label class="text-sm">{{ $t('page.contentModel.description') }}</label>
        <Textarea v-model:value="formData.description" :rows="2" />
      </div>

      <a-divider />

      <!-- 字段定义编辑器 -->
      <div class="flex items-center justify-between">
        <span class="text-base font-medium">
          {{ $t('page.contentModel.fields.title') }}
        </span>
        <Button size="small" @click="addField">
          {{ $t('page.contentModel.fields.add') }}
        </Button>
      </div>

      <Collapse v-model:active-key="activeKeys" v-if="formData.fields.length > 0">
        <CollapsePanel
          v-for="field in formData.fields"
          :key="field.clientKey"
          :header="field.name || $t('page.contentModel.fields.unnamed')"
        >
          <div class="flex flex-col gap-3 p-2">
            <div class="flex justify-end">
              <Button danger size="small" @click="removeField(field.clientKey)">
                {{ $t('page.contentModel.fields.remove') }}
              </Button>
            </div>
            <div class="flex flex-col gap-2">
              <label class="text-sm">{{ $t('page.contentModel.fields.name') }}</label>
              <Input v-model:value="field.name" />
            </div>
            <div class="flex flex-col gap-2">
              <label class="text-sm">{{ $t('page.contentModel.fields.type') }}</label>
              <Select v-model:value="field.type" :options="fieldDefinitionTypeList" />
            </div>
            <div class="flex flex-col gap-2">
              <label class="text-sm">{{ $t('page.contentModel.fields.label') }}</label>
              <Input v-model:value="field.label" />
            </div>
            <div class="flex flex-col gap-2">
              <label class="text-sm">{{ $t('page.contentModel.fields.description') }}</label>
              <Textarea v-model:value="field.description" :rows="2" />
            </div>
            <div class="flex flex-col gap-2">
              <label class="text-sm">{{ $t('page.contentModel.fields.placeholder') }}</label>
              <Input v-model:value="field.placeholder" />
            </div>
            <div class="flex flex-col gap-2">
              <label class="text-sm">{{ $t('page.table.sortOrder') }}</label>
              <InputNumber v-model:value="field.sortOrder" class="w-full" />
            </div>
            <div class="flex items-center gap-2">
              <Switch v-model:checked="field.isRequired" />
              <span class="text-sm">{{ $t('page.contentModel.fields.isRequired') }}</span>
            </div>
            <div class="flex flex-col gap-2">
              <label class="text-sm">{{ $t('page.contentModel.fields.validationRegex') }}</label>
              <Input v-model:value="field.validationRegex" />
            </div>

            <!-- select 类型：选项 JSON 编辑（与 site-setting optionsJson 同模式） -->
            <div v-if="field.type === 'FIELD_TYPE_SELECT'" class="flex flex-col gap-2">
              <label class="text-sm">{{ $t('page.contentModel.fields.options') }}</label>
              <Textarea v-model:value="field.optionsJson" :rows="4" />
            </div>

            <!-- relation 配置 -->
            <div v-if="field.type === 'FIELD_TYPE_RELATION'" class="flex flex-col gap-2 border-t pt-3 mt-2">
              <label class="text-sm">{{ $t('page.contentModel.fields.relation.targetEntityType') }}</label>
              <Select
                v-model:value="field.relationTargetEntityType"
                :options="relationEntityTypeList"
              />
              <div class="flex items-center gap-2 mt-2">
                <Switch v-model:checked="field.relationAllowCrossTenant" />
                <span class="text-sm">{{ $t('page.contentModel.fields.relation.allowCrossTenant') }}</span>
              </div>
              <div class="flex flex-col gap-2 mt-2">
                <label class="text-sm">{{ $t('page.contentModel.fields.relation.filterCategoryId') }}</label>
                <Input v-model:value="field.relationFilterCategoryId" />
              </div>
            </div>
          </div>
        </CollapsePanel>
      </Collapse>

      <div v-else class="text-center text-gray-400 py-8">
        {{ $t('page.contentModel.fields.empty') }}
      </div>
    </div>
  </Drawer>
</template>
