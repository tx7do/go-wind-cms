<script lang="ts" setup>
import { computed, ref } from 'vue';

import { useVbenDrawer } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { notification } from 'ant-design-vue';

import { useVbenForm } from '#/adapter/form';
import { apiClient, commentStatusList, makeUpdateMask } from '#/api';

const data = ref<Record<string, any>>();

const getTitle = computed(() =>
  $t('ui.modal.update', { moduleName: $t('page.comment.moduleName') }),
);

const [BaseForm, baseFormApi] = useVbenForm({
  showDefaultActions: false,
  commonConfig: {
    componentProps: {
      class: 'w-full',
    },
  },
  schema: [
    {
      component: 'Textarea',
      fieldName: 'content',
      label: $t('page.comment.content'),
      componentProps: {
        rows: 4,
        placeholder: $t('ui.placeholder.input'),
        allowClear: true,
      },
      rules: 'required',
    },
    {
      component: 'Select',
      fieldName: 'status',
      label: $t('page.comment.status'),
      rules: 'selectRequired',
      componentProps: {
        options: commentStatusList,
        placeholder: $t('ui.placeholder.select'),
        filterOption: (input: string, option: any) =>
          option.label.toLowerCase().includes(input.toLowerCase()),
        allowClear: true,
        showSearch: true,
      },
    },
    {
      component: 'Switch',
      fieldName: 'isSpam',
      label: $t('page.comment.isSpam'),
      help: $t('page.comment.help.isSpam'),
    },
    {
      component: 'Switch',
      fieldName: 'isSticky',
      label: $t('page.comment.isSticky'),
      help: $t('page.comment.help.isSticky'),
    },
  ],
});

const [Drawer, drawerApi] = useVbenDrawer({
  onCancel() {
    drawerApi.close();
  },

  async onConfirm() {
    const validate = await baseFormApi.validate();
    if (!validate.valid) {
      return;
    }

    setLoading(true);
    const values = await baseFormApi.getValues();

    try {
      await apiClient.commentService.Update({
        id: data.value?.row?.id,
        data: values as any,
        updateMask: makeUpdateMask(Object.keys(values)),
      });

      notification.success({
        message: $t('ui.notification.update_success'),
      });
      drawerApi.close();
    } catch {
      notification.error({
        message: $t('ui.notification.update_failed'),
      });
    } finally {
      setLoading(false);
    }
  },

  onOpenChange(isOpen) {
    if (isOpen) {
      data.value = drawerApi.getData<Record<string, any>>();
      baseFormApi.setValues(data.value?.row);
      setLoading(false);
    }
  },
});

function setLoading(loading: boolean) {
  drawerApi.setState({ loading });
}
</script>

<template>
  <Drawer :title="getTitle">
    <BaseForm />
  </Drawer>
</template>
