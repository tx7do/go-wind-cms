<script lang="ts" setup>
import { computed, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { notification } from 'ant-design-vue';

import { useVbenForm } from '#/adapter/form';
import {
  apiClient,
  enableBoolList,
  type internal_messageservicev1_InternalMessageCategory as InternalMessageCategory,
  fetchListMessageCategories,
  makeUpdateMask,
  PaginationQuery,
} from '#/api';

const data = ref();

const getTitle = computed(() =>
  data.value?.create
    ? $t('ui.modal.create', {
        moduleName: $t('page.internalMessageCategory.moduleName'),
      })
    : $t('ui.modal.update', {
        moduleName: $t('page.internalMessageCategory.moduleName'),
      }),
);
// const isCreate = computed(() => data.value?.create);

const [BaseForm, baseFormApi] = useVbenForm({
  showDefaultActions: false,
  // 所有表单项共用，可单独在表单内覆盖
  commonConfig: {
    // 所有表单项
    componentProps: {
      class: 'w-full',
    },
  },
  schema: [
    {
      component: 'Input',
      fieldName: 'name',
      label: $t('page.internalMessageCategory.name'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
      rules: 'required',
    },
    {
      component: 'Input',
      fieldName: 'code',
      label: $t('page.internalMessageCategory.code'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
      rules: 'required',
    },
    {
      component: 'Select',
      fieldName: 'parentId',
      label: $t('page.internalMessageCategory.parentId'),
      componentProps: {
        placeholder: $t('ui.placeholder.select'),
        allowClear: true,
        options: [] as { label: string; value: number }[],
      },
    },
    {
      component: 'InputNumber',
      fieldName: 'sortOrder',
      defaultValue: 1,
      label: $t('ui.table.sortOrder'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
    },
    {
      component: 'RadioGroup',
      fieldName: 'isEnabled',
      label: $t('ui.table.status'),
      defaultValue: true,
      rules: 'selectRequired',
      componentProps: {
        optionType: 'button',
        buttonStyle: 'solid',
        class: 'flex flex-wrap', // 如果选项过多，可以添加class来自动折叠
        options: enableBoolList,
      },
    },
    {
      component: 'Input',
      fieldName: 'iconUrl',
      label: $t('page.internalMessageCategory.iconUrl'),
      componentProps: {
        placeholder: $t('page.internalMessageCategory.placeholder.iconUrl'),
        allowClear: true,
      },
    },
  ],
});

const [Drawer, drawerApi] = useVbenDrawer({
  onCancel() {
    drawerApi.close();
  },

  async onConfirm() {

    // 校验输入的数据
    const validate = await baseFormApi.validate();
    if (!validate.valid) {
      return;
    }

    setLoading(true);

    // 获取表单数据
    const values = await baseFormApi.getValues();

    // parentId 空值归一为 0（顶级），支持通过清空选择把分类移回顶级
    if (values.parentId === undefined || values.parentId === null) {
      values.parentId = 0;
    }


    try {
      await (data.value?.create
        ? apiClient.internalMessageCategoryService.Create({
            data: { ...values } as InternalMessageCategory,
          })
        : apiClient.internalMessageCategoryService.Update({
            id: data.value.row.id,
            data: { ...values } as InternalMessageCategory,
            updateMask: makeUpdateMask(
              Object.keys(values).filter((k) => !['code'].includes(k)),
            ),
          }));

      notification.success({
        message: data.value?.create
          ? $t('ui.notification.create_success')
          : $t('ui.notification.update_success'),
      });
    } catch {
      notification.error({
        message: data.value?.create
          ? $t('ui.notification.create_failed')
          : $t('ui.notification.update_failed'),
      });
    } finally {
      // 关闭窗口
      drawerApi.close();
      setLoading(false);
    }
  },

  onOpenChange(isOpen: boolean) {
    if (isOpen) {
      // 获取传入的数据
      data.value = drawerApi.getData<Record<string, any>>();

      // 为表单赋值
      baseFormApi.setValues(data.value?.row);

      // 加载父分类选项（排除自身，防自环；后代环由后端祖先链校验兜底）
      loadParentOptions();

      setLoading(false);

    }
  },
});

function setLoading(loading: boolean) {
  drawerApi.setState({ confirmLoading: loading });
}

// loadParentOptions 拉取全量分类（扁平）构建父分类下拉选项，编辑模式排除自身
async function loadParentOptions() {
  try {
    const resp = await fetchListMessageCategories(
      new PaginationQuery({
        paging: { page: 1, pageSize: 1000 },
        formValues: {},
      }),
    );
    const selfId = data.value?.row?.id;
    const items = resp?.items ?? [];
    const options = items
      .filter((item: Record<string, any>) => item.id !== selfId)
      .map((item: Record<string, any>) => ({
        label: `${item.name} (${item.code})`,
        value: item.id,
      }));

    await baseFormApi.updateSchema([
      {
        fieldName: 'parentId',
        componentProps: {
          placeholder: $t('ui.placeholder.select'),
          allowClear: true,
          options,
        },
      },
    ]);
  } catch {
    // 选项加载失败不阻塞表单，仅留下空下拉
  }
}
</script>

<template>
  <Drawer :title="getTitle" class="w-full max-w-[800px]">
    <BaseForm />
  </Drawer>
</template>
