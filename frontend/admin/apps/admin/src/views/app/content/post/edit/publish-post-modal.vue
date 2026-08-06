<script setup lang="ts">
import { ref } from 'vue';

import { useVbenModal } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { notification } from 'ant-design-vue';

import { useVbenForm } from '#/adapter/form';
import { apiClient, makeUpdateMask } from '#/api';

const data = ref();

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
      component: 'Textarea',
      fieldName: 'summary',
      label: $t('page.post.summary'),
      componentProps: {
        placeholder: $t('ui.placeholder.input'),
      },
      rules: 'required',
    },
  ],
});

const [Modal, modalApi] = useVbenModal({
  onCancel() {
    modalApi.close();
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

    const postData = modalApi.getData<any>();
    if (postData?.id) {
      const data = {
        summary: values.summary,
      };
      try {
        await apiClient.postService.Update({
          id: postData.id,
          data,
          updateMask: makeUpdateMask(Object.keys(data)),
        });

        notification.success({
          message: $t('ui.notification.update_success'),
        });
      } catch {
        notification.error({
          message: $t('ui.notification.update_failed'),
        });
      }
    }

    setLoading(false);
    modalApi.close();
  },

  onOpenChange(isOpen: boolean) {
    if (isOpen) {
      // 获取传入的数据
      data.value = modalApi.getData<any>();

      setLoading(false);

    }
  },
});

function setLoading(loading: boolean) {
  modalApi.setState({ confirmLoading: loading });
}
</script>

<template>
  <Modal :title="$t('page.post.button.publish')">
    <BaseForm />
  </Modal>
</template>

<style scoped></style>
